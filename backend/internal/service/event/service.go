package eventsvc

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"bimeet/internal/model"
)

var (
	ErrForbidden = errors.New("forbidden")
	ErrFull      = errors.New("event is full")
)

type Service struct {
	eventRepo        EventRepo
	userRepo         UserRepo
	notificationRepo NotificationRepo
	mailer           Mailer
}

func New(
	eventRepo EventRepo,
	userRepo UserRepo,
	notificationRepo NotificationRepo,
	mailer Mailer,
) *Service {
	return &Service{
		eventRepo:        eventRepo,
		userRepo:         userRepo,
		notificationRepo: notificationRepo,
		mailer:           mailer,
	}
}

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]model.EventListItem, error) {
	return s.eventRepo.ListForUserEnriched(ctx, userID)
}

func (s *Service) Create(ctx context.Context, req model.CreateEventRequest, organizerID uuid.UUID) (model.Event, error) {
	if req.Title == "" {
		return model.Event{}, fmt.Errorf("title is required")
	}
	if req.DateStart.IsZero() || req.DateEnd.IsZero() {
		return model.Event{}, fmt.Errorf("date_start and date_end are required")
	}
	if req.DateEnd.Before(req.DateStart) {
		return model.Event{}, fmt.Errorf("date_end must be after date_start")
	}
	if req.Category == "" {
		req.Category = "ordinary"
	}
	if req.Category != "ordinary" && req.Category != "business" {
		return model.Event{}, fmt.Errorf("category must be ordinary or business")
	}
	return s.eventRepo.Create(ctx, req, organizerID)
}

func (s *Service) GetByID(ctx context.Context, id, userID uuid.UUID) (model.EventDetail, error) {
	detail, err := s.eventRepo.GetDetail(ctx, id, userID)
	if err != nil {
		return model.EventDetail{}, fmt.Errorf("event not found")
	}

	if !detail.IsPublic {
		if detail.Organizer.ID != userID {
			isParticipant := false
			for _, p := range detail.Participants {
				if p.User.ID == userID {
					isParticipant = true
					break
				}
			}
			if !isParticipant {
				return model.EventDetail{}, ErrForbidden
			}
		}
	}

	return detail, nil
}

func (s *Service) ListPublic(ctx context.Context, userID uuid.UUID) ([]model.PublicEventListItem, error) {
	return s.eventRepo.ListPublic(ctx, userID)
}

func (s *Service) JoinPublic(ctx context.Context, eventID, userID uuid.UUID) error {
	detail, err := s.eventRepo.GetDetail(ctx, eventID, userID)
	if err != nil {
		return fmt.Errorf("event not found")
	}
	if !detail.IsPublic {
		return ErrForbidden
	}
	if detail.Status != "active" {
		return fmt.Errorf("event is not active")
	}
	if detail.Organizer.ID == userID {
		return fmt.Errorf("you are the organizer")
	}
	if detail.MaxGuests != nil && detail.ConfirmedCount >= *detail.MaxGuests {
		return ErrFull
	}

	if _, err := s.eventRepo.JoinPublic(ctx, eventID, userID); err != nil {
		return err
	}

	go func() {
		bgCtx := context.Background()
		user, err := s.userRepo.GetByID(bgCtx, userID)
		if err != nil {
			return
		}
		_, _ = s.notificationRepo.Create(bgCtx, detail.Organizer.ID, &eventID,
			"event_joined", fmt.Sprintf("%s joined your event '%s'", user.Name, detail.Title))
	}()

	return nil
}

func (s *Service) GetByInviteToken(ctx context.Context, token uuid.UUID) (model.InviteEventInfo, error) {
	return s.eventRepo.GetByInviteToken(ctx, token)
}

func (s *Service) JoinByInviteToken(ctx context.Context, token, userID uuid.UUID, req model.JoinByInviteTokenRequest) (model.EventParticipant, error) {
	event, err := s.eventRepo.GetByInviteToken(ctx, token)
	if err != nil {
		return model.EventParticipant{}, fmt.Errorf("event not found")
	}

	if event.Organizer.ID == userID {
		return model.EventParticipant{}, fmt.Errorf("you are the organizer of this event")
	}

	action := req.Action
	if action == "" {
		action = "join"
	}

	var status string
	switch action {
	case "join":
		status = "confirmed"
	case "decline":
		status = "declined"
	default:
		return model.EventParticipant{}, fmt.Errorf("action must be join or decline")
	}

	p, err := s.eventRepo.AddParticipant(ctx, event.ID, userID, status)
	if err != nil {
		return model.EventParticipant{}, err
	}
	return p, nil
}

func (s *Service) Update(ctx context.Context, id, userID uuid.UUID, req model.UpdateEventRequest) (model.Event, error) {
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return model.Event{}, fmt.Errorf("event not found")
	}

	if event.OrganizerID != userID {
		return model.Event{}, fmt.Errorf("forbidden")
	}

	type fieldChange struct {
		name    string
		oldVal  string
		newVal  string
		changed bool
	}

	if req.Category != nil && *req.Category != "ordinary" && *req.Category != "business" {
		return model.Event{}, fmt.Errorf("category must be ordinary or business")
	}

	changes := []fieldChange{
		{name: "title", oldVal: event.Title, newVal: func() string {
			if req.Title != nil {
				return *req.Title
			}
			return event.Title
		}(), changed: req.Title != nil && *req.Title != event.Title},
		{name: "description", oldVal: event.Description, newVal: func() string {
			if req.Description != nil {
				return *req.Description
			}
			return event.Description
		}(), changed: req.Description != nil && *req.Description != event.Description},
		{name: "location", oldVal: event.Location, newVal: func() string {
			if req.Location != nil {
				return *req.Location
			}
			return event.Location
		}(), changed: req.Location != nil && *req.Location != event.Location},
		{name: "date_start", oldVal: event.DateStart.String(), newVal: func() string {
			if req.DateStart != nil {
				return req.DateStart.String()
			}
			return event.DateStart.String()
		}(), changed: req.DateStart != nil && !req.DateStart.Equal(event.DateStart)},
		{name: "date_end", oldVal: event.DateEnd.String(), newVal: func() string {
			if req.DateEnd != nil {
				return req.DateEnd.String()
			}
			return event.DateEnd.String()
		}(), changed: req.DateEnd != nil && !req.DateEnd.Equal(event.DateEnd)},
		{name: "category", oldVal: event.Category, newVal: func() string {
			if req.Category != nil {
				return *req.Category
			}
			return event.Category
		}(), changed: req.Category != nil && *req.Category != event.Category},
	}

	updated, err := s.eventRepo.Update(ctx, id, req)
	if err != nil {
		return model.Event{}, err
	}

	for _, ch := range changes {
		if ch.changed {
			_ = s.eventRepo.AddChangeLog(ctx, id, userID, ch.name, ch.oldVal, ch.newVal)
		}
	}

	go func() {
		bgCtx := context.Background()
		participants, err := s.eventRepo.ListParticipants(bgCtx, id)
		if err != nil {
			return
		}
		for _, p := range participants {
			if p.UserID == userID {
				continue
			}
			_, _ = s.notificationRepo.Create(bgCtx, p.UserID, &id, "event_updated",
				fmt.Sprintf("Event '%s' has been updated", updated.Title))
		}
	}()

	return updated, nil
}

func (s *Service) Delete(ctx context.Context, id, userID uuid.UUID) error {
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("event not found")
	}
	if event.OrganizerID != userID {
		return fmt.Errorf("forbidden")
	}
	return s.eventRepo.Delete(ctx, id)
}

func (s *Service) Complete(ctx context.Context, id, userID uuid.UUID) error {
	event, err := s.eventRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("event not found")
	}
	if event.OrganizerID != userID {
		return fmt.Errorf("forbidden")
	}
	if err := s.eventRepo.Complete(ctx, id); err != nil {
		return err
	}
	go func() {
		bgCtx := context.Background()
		participants, err := s.eventRepo.ListParticipants(bgCtx, id)
		if err != nil {
			return
		}
		for _, p := range participants {
			if p.UserID == userID {
				continue
			}
			_, _ = s.notificationRepo.Create(bgCtx, p.UserID, &id, "event_completed",
				fmt.Sprintf("Event '%s' has been completed", event.Title))
		}
	}()
	return nil
}

// ─── Participants ─────────────────────────────────────────────────────────

func (s *Service) InviteParticipant(ctx context.Context, eventID, organizerID uuid.UUID, req model.InviteParticipantRequest) (model.EventParticipant, error) {
	event, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return model.EventParticipant{}, fmt.Errorf("event not found")
	}
	if event.OrganizerID != organizerID {
		return model.EventParticipant{}, fmt.Errorf("forbidden")
	}

	organizer, err := s.userRepo.GetByID(ctx, organizerID)
	if err != nil {
		return model.EventParticipant{}, fmt.Errorf("organizer not found")
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		go func() {
			if err := s.mailer.SendInvite(req.Email, event.Title, organizer.Name, event.InviteToken.String()); err != nil {
				log.Printf("send invite to %s: %v", req.Email, err)
			}
		}()
		return model.EventParticipant{}, nil
	}

	if user.ID == organizerID {
		return model.EventParticipant{}, fmt.Errorf("cannot invite yourself")
	}

	p, err := s.eventRepo.AddParticipant(ctx, eventID, user.ID, "invited")
	if err != nil {
		return model.EventParticipant{}, err
	}

	go func() {
		_, _ = s.notificationRepo.Create(context.Background(), user.ID, &eventID, "event_invite",
			fmt.Sprintf("You have been invited to event '%s'", event.Title))
		if err := s.mailer.SendInvite(user.Email, event.Title, organizer.Name, event.InviteToken.String()); err != nil {
			log.Printf("send invite to %s: %v", user.Email, err)
		}
	}()

	return p, nil
}

func (s *Service) UpdateParticipantStatus(ctx context.Context, eventID, targetUserID, callerID uuid.UUID, req model.UpdateParticipantStatusRequest) (model.EventParticipant, error) {
	if req.Status != "confirmed" && req.Status != "declined" {
		return model.EventParticipant{}, fmt.Errorf("status must be confirmed or declined")
	}

	if callerID != targetUserID {
		return model.EventParticipant{}, fmt.Errorf("forbidden")
	}

	_, err := s.eventRepo.GetByID(ctx, eventID)
	if err != nil {
		return model.EventParticipant{}, fmt.Errorf("event not found")
	}

	p, err := s.eventRepo.UpdateParticipantStatus(ctx, eventID, targetUserID, req.Status)
	if err != nil {
		if isNotFound(err) {
			return model.EventParticipant{}, fmt.Errorf("participant not found")
		}
		return model.EventParticipant{}, err
	}
	return p, nil
}

func isNotFound(err error) bool {
	return err == pgx.ErrNoRows
}
