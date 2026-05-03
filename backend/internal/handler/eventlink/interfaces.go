package eventlinkhandler

//go:generate mockgen -destination=mock/mock.go -package=mock . EventLinkService

import (
	"context"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

type EventLinkService interface {
	List(ctx context.Context, eventID, userID uuid.UUID) ([]model.EventLink, error)
	Create(ctx context.Context, eventID, userID uuid.UUID, req model.CreateEventLinkRequest) (model.EventLink, error)
	Delete(ctx context.Context, eventID, linkID, userID uuid.UUID) error
}
