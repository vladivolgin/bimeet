package carpoolhandler

//go:generate mockgen -destination=mock/mock.go -package=mock . CarpoolService

import (
	"context"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

type CarpoolService interface {
	List(ctx context.Context, eventID, userID uuid.UUID) ([]model.Carpool, error)
	Create(ctx context.Context, eventID, driverID uuid.UUID, req model.CreateCarpoolRequest) (model.Carpool, error)
	Join(ctx context.Context, eventID, carpoolID, userID uuid.UUID) error
}
