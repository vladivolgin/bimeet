package eventlinksvc

//go:generate mockgen -destination=mock/mock.go -package=mock . EventLinkRepo,EventRepo

import (
	"context"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

type EventLinkRepo interface {
	Create(ctx context.Context, eventID, userID uuid.UUID, req model.CreateEventLinkRequest) (model.EventLink, error)
	List(ctx context.Context, eventID uuid.UUID) ([]model.EventLink, error)
	Delete(ctx context.Context, linkID, eventID uuid.UUID) error
}

type EventRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.Event, error)
	IsParticipant(ctx context.Context, eventID, userID uuid.UUID) (bool, error)
}
