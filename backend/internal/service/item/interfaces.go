package itemsvc

//go:generate mockgen -destination=mock/mock.go -package=mock . ItemRepo,EventRepo

import (
	"context"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

type ItemRepo interface {
	Create(ctx context.Context, eventID uuid.UUID, req model.CreateItemRequest) (model.Item, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Item, error)
	ListByEvent(ctx context.Context, eventID uuid.UUID) ([]model.Item, error)
	UpdateAssignment(ctx context.Context, id uuid.UUID, assignedTo *uuid.UUID) (model.Item, error)
}

type EventRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.Event, error)
	IsParticipant(ctx context.Context, eventID, userID uuid.UUID) (bool, error)
	IsConfirmedParticipant(ctx context.Context, eventID, userID uuid.UUID) (bool, error)
}
