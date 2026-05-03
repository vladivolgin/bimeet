package collectionsvc

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"strings"
	"time"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

type Service struct {
	collectionRepo   CollectionRepo
	eventRepo        EventRepo
	notificationRepo NotificationRepo
	uploader         Uploader
}

func New(collectionRepo CollectionRepo, eventRepo EventRepo, notificationRepo NotificationRepo, uploader Uploader) *Service {
	return &Service{
		collectionRepo:   collectionRepo,
		eventRepo:        eventRepo,
		notificationRepo: notificationRepo,
		uploader:         uploader,
	}
}

func (s *Service) List(ctx context.Context, eventID, userID uuid.UUID) ([]model.Collection, error) {
	ok, err := s.eventRepo.IsParticipant(ctx, eventID, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, fmt.Errorf("forbidden")
	}
	return s.collectionRepo.ListByEvent(ctx, eventID)
}

func (s *Service) Create(ctx context.Context, eventID, userID uuid.UUID, req model.CreateCollectionRequest) (model.Collection, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return model.Collection{}, fmt.Errorf("event not found")
	}
	if event.OrganizerID != userID {
		return model.Collection{}, fmt.Errorf("forbidden")
	}
	if req.Title == "" {
		return model.Collection{}, fmt.Errorf("title is required")
	}
	if req.PerPersonAmount <= 0 {
		return model.Collection{}, fmt.Errorf("per_person_amount must be positive")
	}
	return s.collectionRepo.Create(ctx, eventID, req, userID)
}

func (s *Service) Delete(ctx context.Context, eventID, collectionID, userID uuid.UUID) error {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return fmt.Errorf("event not found")
	}
	if event.OrganizerID != userID {
		return fmt.Errorf("forbidden")
	}

	col, err := s.collectionRepo.GetByID(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("collection not found")
	}
	if col.EventID != eventID {
		return fmt.Errorf("collection does not belong to this event")
	}

	paidCount, err := s.collectionRepo.CountPaidContributions(ctx, collectionID)
	if err != nil {
		return fmt.Errorf("cannot check contributions: %w", err)
	}
	if paidCount > 0 {
		return fmt.Errorf("cannot delete collection with paid contributions")
	}

	return s.collectionRepo.Delete(ctx, collectionID)
}

// SubmitContribution uploads receipt to S3 and sets contribution status to pending.
func (s *Service) SubmitContribution(ctx context.Context, eventID, collectionID, userID uuid.UUID, file io.Reader, header *multipart.FileHeader) (model.CollectionContribution, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return model.CollectionContribution{}, fmt.Errorf("event not found")
	}

	// Organizer uses MarkPaid instead
	if event.OrganizerID == userID {
		return model.CollectionContribution{}, fmt.Errorf("forbidden: organizer should use mark-paid")
	}

	ok, err := s.eventRepo.IsConfirmedParticipant(ctx, eventID, userID)
	if err != nil {
		return model.CollectionContribution{}, err
	}
	if !ok {
		return model.CollectionContribution{}, fmt.Errorf("forbidden: must be a confirmed participant")
	}

	col, err := s.collectionRepo.GetByID(ctx, collectionID)
	if err != nil {
		return model.CollectionContribution{}, fmt.Errorf("collection not found")
	}
	if col.EventID != eventID {
		return model.CollectionContribution{}, fmt.Errorf("collection does not belong to this event")
	}

	// Check current status — cannot resubmit if already paid
	existing, err := s.collectionRepo.GetContribution(ctx, collectionID, userID)
	if err == nil && existing.Status == "paid" {
		return model.CollectionContribution{}, fmt.Errorf("already paid: cannot resubmit")
	}

	contentType := header.Header.Get("Content-Type")
	if !isAllowedReceiptType(contentType) {
		return model.CollectionContribution{}, fmt.Errorf("unsupported file type: must be image or PDF")
	}

	key := fmt.Sprintf("receipts/%s/%s/%d", collectionID, userID, time.Now().UnixMilli())
	receiptURL, err := s.uploader.Upload(ctx, key, contentType, file)
	if err != nil {
		return model.CollectionContribution{}, fmt.Errorf("upload receipt: %w", err)
	}

	contrib, err := s.collectionRepo.SubmitContribution(ctx, collectionID, userID, receiptURL)
	if err != nil {
		return model.CollectionContribution{}, err
	}

	go func() {
		bgCtx := context.Background()
		msg := fmt.Sprintf("Новый чек на подтверждение в сборе «%s»", col.Title)
		_, _ = s.notificationRepo.Create(bgCtx, event.OrganizerID, &eventID, "collection_contribution_pending", msg)
	}()

	return contrib, nil
}

// ConfirmContribution approves a pending contribution (organizer only).
func (s *Service) ConfirmContribution(ctx context.Context, eventID, collectionID, contributionID, organizerID uuid.UUID) (model.CollectionContribution, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return model.CollectionContribution{}, fmt.Errorf("event not found")
	}
	if event.OrganizerID != organizerID {
		return model.CollectionContribution{}, fmt.Errorf("forbidden")
	}

	contrib, err := s.collectionRepo.GetContributionByID(ctx, contributionID)
	if err != nil {
		return model.CollectionContribution{}, fmt.Errorf("contribution not found")
	}
	if contrib.CollectionID != collectionID {
		return model.CollectionContribution{}, fmt.Errorf("contribution does not belong to this collection")
	}
	if contrib.Status != "pending" {
		return model.CollectionContribution{}, fmt.Errorf("contribution is not pending")
	}

	result, err := s.collectionRepo.ConfirmContribution(ctx, contributionID)
	if err != nil {
		return model.CollectionContribution{}, err
	}

	go func() {
		bgCtx := context.Background()
		col, _ := s.collectionRepo.GetByID(bgCtx, collectionID)
		msg := fmt.Sprintf("Ваш взнос в сборе «%s» подтверждён", col.Title)
		_, _ = s.notificationRepo.Create(bgCtx, contrib.UserID, &eventID, "collection_contribution_confirmed", msg)
	}()

	return result, nil
}

// RejectContribution rejects a pending contribution (organizer only).
func (s *Service) RejectContribution(ctx context.Context, eventID, collectionID, contributionID, organizerID uuid.UUID) (model.CollectionContribution, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return model.CollectionContribution{}, fmt.Errorf("event not found")
	}
	if event.OrganizerID != organizerID {
		return model.CollectionContribution{}, fmt.Errorf("forbidden")
	}

	contrib, err := s.collectionRepo.GetContributionByID(ctx, contributionID)
	if err != nil {
		return model.CollectionContribution{}, fmt.Errorf("contribution not found")
	}
	if contrib.CollectionID != collectionID {
		return model.CollectionContribution{}, fmt.Errorf("contribution does not belong to this collection")
	}
	if contrib.Status != "pending" {
		return model.CollectionContribution{}, fmt.Errorf("contribution is not pending")
	}

	result, err := s.collectionRepo.RejectContribution(ctx, contributionID)
	if err != nil {
		return model.CollectionContribution{}, err
	}

	go func() {
		bgCtx := context.Background()
		col, _ := s.collectionRepo.GetByID(bgCtx, collectionID)
		msg := fmt.Sprintf("Ваш чек в сборе «%s» отклонён организатором", col.Title)
		_, _ = s.notificationRepo.Create(bgCtx, contrib.UserID, &eventID, "collection_contribution_rejected", msg)
	}()

	return result, nil
}

// MarkPaid forcibly marks a participant as paid (organizer action, no receipt required).
func (s *Service) MarkPaid(ctx context.Context, eventID, collectionID, targetUserID, organizerID uuid.UUID) (model.CollectionContribution, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return model.CollectionContribution{}, fmt.Errorf("event not found")
	}
	if event.OrganizerID != organizerID {
		return model.CollectionContribution{}, fmt.Errorf("forbidden")
	}

	col, err := s.collectionRepo.GetByID(ctx, collectionID)
	if err != nil {
		return model.CollectionContribution{}, fmt.Errorf("collection not found")
	}
	if col.EventID != eventID {
		return model.CollectionContribution{}, fmt.Errorf("collection does not belong to this event")
	}

	ok, err := s.eventRepo.IsParticipant(ctx, eventID, targetUserID)
	if err != nil {
		return model.CollectionContribution{}, err
	}
	if !ok {
		return model.CollectionContribution{}, fmt.Errorf("target user is not a participant")
	}

	return s.collectionRepo.MarkPaid(ctx, collectionID, targetUserID)
}

func (s *Service) Summary(ctx context.Context, eventID, userID uuid.UUID) (model.EventCollectionsSummaryResponse, error) {
	ok, err := s.eventRepo.IsParticipant(ctx, eventID, userID)
	if err != nil {
		return model.EventCollectionsSummaryResponse{}, err
	}
	if !ok {
		return model.EventCollectionsSummaryResponse{}, fmt.Errorf("forbidden")
	}

	confirmedCount, err := s.eventRepo.CountConfirmedParticipants(ctx, eventID)
	if err != nil {
		return model.EventCollectionsSummaryResponse{}, err
	}

	collections, err := s.collectionRepo.ListByEvent(ctx, eventID)
	if err != nil {
		return model.EventCollectionsSummaryResponse{}, err
	}

	var summaries []model.CollectionSummary
	var grandTotal, totalPaid float64

	for _, col := range collections {
		paidCount, err := s.collectionRepo.CountPaidContributions(ctx, col.ID)
		if err != nil {
			return model.EventCollectionsSummaryResponse{}, err
		}

		perPerson := col.PerPersonAmount
		expectedTotal := perPerson * float64(confirmedCount)
		collectionPaid := float64(paidCount) * perPerson
		remaining := float64(confirmedCount-paidCount) * perPerson

		summaries = append(summaries, model.CollectionSummary{
			CollectionID:    col.ID,
			Title:           col.Title,
			PerPersonAmount: perPerson,
			ExpectedTotal:   expectedTotal,
			PaidCount:       paidCount,
			TotalCount:      confirmedCount,
			TotalPaid:       collectionPaid,
			Remaining:       remaining,
		})

		grandTotal += expectedTotal
		totalPaid += collectionPaid
	}

	return model.EventCollectionsSummaryResponse{
		Collections: summaries,
		GrandTotal:  grandTotal,
		TotalPaid:   totalPaid,
		Remaining:   grandTotal - totalPaid,
	}, nil
}

func isAllowedReceiptType(contentType string) bool {
	ct := strings.ToLower(contentType)
	return strings.HasPrefix(ct, "image/") || ct == "application/pdf"
}
