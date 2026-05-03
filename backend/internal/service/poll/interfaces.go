package pollsvc

//go:generate mockgen -destination=mock/mock.go -package=mock . PollRepo,EventRepo

import (
	"context"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

type PollRepo interface {
	Create(ctx context.Context, eventID uuid.UUID, req model.CreatePollRequest, createdBy uuid.UUID) (model.Poll, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Poll, error)
	ListByEvent(ctx context.Context, eventID uuid.UUID) ([]model.Poll, error)
	GetOption(ctx context.Context, optionID uuid.UUID) (model.PollOption, error)
	Vote(ctx context.Context, optionID, userID uuid.UUID) error
}

type EventRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.Event, error)
	IsParticipant(ctx context.Context, eventID, userID uuid.UUID) (bool, error)
	IsConfirmedParticipant(ctx context.Context, eventID, userID uuid.UUID) (bool, error)
}
