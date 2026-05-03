package pollhandler

//go:generate mockgen -destination=mock/mock.go -package=mock . PollService

import (
	"context"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

type PollService interface {
	List(ctx context.Context, eventID, userID uuid.UUID) ([]model.Poll, error)
	Create(ctx context.Context, eventID, userID uuid.UUID, req model.CreatePollRequest) (model.Poll, error)
	Vote(ctx context.Context, eventID, pollID, userID uuid.UUID, req model.VoteRequest) error
}
