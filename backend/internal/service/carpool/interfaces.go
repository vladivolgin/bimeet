package carpoolsvc

//go:generate mockgen -destination=mock/mock.go -package=mock . CarpoolRepo,EventRepo

import (
	"context"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

type CarpoolRepo interface {
	Create(ctx context.Context, eventID, driverID uuid.UUID, req model.CreateCarpoolRequest) (model.Carpool, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Carpool, error)
	ListByEvent(ctx context.Context, eventID uuid.UUID) ([]model.Carpool, error)
	AddPassenger(ctx context.Context, carpoolID, userID uuid.UUID) error
	CountPassengers(ctx context.Context, carpoolID uuid.UUID) (int, error)
	IsPassenger(ctx context.Context, carpoolID, userID uuid.UUID) (bool, error)
}

type EventRepo interface {
	IsParticipant(ctx context.Context, eventID, userID uuid.UUID) (bool, error)
	IsConfirmedParticipant(ctx context.Context, eventID, userID uuid.UUID) (bool, error)
}
