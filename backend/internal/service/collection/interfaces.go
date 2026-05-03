package collectionsvc

//go:generate mockgen -destination=mock/mock.go -package=mock . CollectionRepo,EventRepo,NotificationRepo

import (
	"context"
	"io"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

type CollectionRepo interface {
	Create(ctx context.Context, eventID uuid.UUID, req model.CreateCollectionRequest, createdBy uuid.UUID) (model.Collection, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Collection, error)
	ListByEvent(ctx context.Context, eventID uuid.UUID) ([]model.Collection, error)
	Delete(ctx context.Context, id uuid.UUID) error
	GetContribution(ctx context.Context, collectionID, userID uuid.UUID) (model.CollectionContribution, error)
	GetContributionByID(ctx context.Context, id uuid.UUID) (model.CollectionContribution, error)
	SubmitContribution(ctx context.Context, collectionID, userID uuid.UUID, receiptURL string) (model.CollectionContribution, error)
	ConfirmContribution(ctx context.Context, id uuid.UUID) (model.CollectionContribution, error)
	RejectContribution(ctx context.Context, id uuid.UUID) (model.CollectionContribution, error)
	MarkPaid(ctx context.Context, collectionID, userID uuid.UUID) (model.CollectionContribution, error)
	CountPaidContributions(ctx context.Context, collectionID uuid.UUID) (int, error)
}

type EventRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.Event, error)
	IsParticipant(ctx context.Context, eventID, userID uuid.UUID) (bool, error)
	IsConfirmedParticipant(ctx context.Context, eventID, userID uuid.UUID) (bool, error)
	CountConfirmedParticipants(ctx context.Context, eventID uuid.UUID) (int, error)
}

type NotificationRepo interface {
	Create(ctx context.Context, userID uuid.UUID, eventID *uuid.UUID, notifType, message string) (model.Notification, error)
}

type Uploader interface {
	Upload(ctx context.Context, key, contentType string, body io.Reader) (string, error)
}
