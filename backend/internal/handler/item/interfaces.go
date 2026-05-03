package itemhandler

//go:generate mockgen -destination=mock/mock.go -package=mock . ItemService

import (
	"context"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

type ItemService interface {
	List(ctx context.Context, eventID, userID uuid.UUID) ([]model.Item, error)
	Create(ctx context.Context, eventID, userID uuid.UUID, req model.CreateItemRequest) (model.Item, error)
	UpdateAssignment(ctx context.Context, eventID, itemID, userID uuid.UUID, req model.UpdateItemRequest) (model.Item, error)
}
