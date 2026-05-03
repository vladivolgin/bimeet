package notificationhandler

//go:generate mockgen -destination=mock/mock.go -package=mock . NotificationService

import (
	"context"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

type NotificationService interface {
	List(ctx context.Context, userID uuid.UUID) ([]model.Notification, error)
	MarkRead(ctx context.Context, notifID, userID uuid.UUID) (model.Notification, error)
	MarkAllRead(ctx context.Context, userID uuid.UUID) error
	Delete(ctx context.Context, notifID, userID uuid.UUID) error
	DeleteAll(ctx context.Context, userID uuid.UUID) error
}
