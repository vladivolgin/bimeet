package reminder

import (
	"context"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

type EventRepo interface {
	ListForReminder(ctx context.Context) ([]model.ReminderEvent, error)
	MarkReminder3dSent(ctx context.Context, id uuid.UUID) error
	MarkReminder1dSent(ctx context.Context, id uuid.UUID) error
}

type NotificationRepo interface {
	Create(ctx context.Context, userID uuid.UUID, eventID *uuid.UUID, notifType, message string) (model.Notification, error)
}
