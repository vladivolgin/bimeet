package eventsvc

//go:generate mockgen -destination=mock/mock.go -package=mock . EventRepo,UserRepo,NotificationRepo,Mailer

import (
	"context"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

type EventRepo interface {
	Create(ctx context.Context, req model.CreateEventRequest, organizerID uuid.UUID) (model.Event, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.Event, error)
	GetByInviteToken(ctx context.Context, token uuid.UUID) (model.InviteEventInfo, error)
	GetDetail(ctx context.Context, eventID, userID uuid.UUID) (model.EventDetail, error)
	ListForUserEnriched(ctx context.Context, userID uuid.UUID) ([]model.EventListItem, error)
	ListPublic(ctx context.Context, userID uuid.UUID) ([]model.PublicEventListItem, error)
	JoinPublic(ctx context.Context, eventID, userID uuid.UUID) (model.EventParticipant, error)
	Update(ctx context.Context, id uuid.UUID, req model.UpdateEventRequest) (model.Event, error)
	Delete(ctx context.Context, id uuid.UUID) error
	Complete(ctx context.Context, id uuid.UUID) error
	AddParticipant(ctx context.Context, eventID, userID uuid.UUID, status string) (model.EventParticipant, error)
	UpdateParticipantStatus(ctx context.Context, eventID, userID uuid.UUID, status string) (model.EventParticipant, error)
	ListParticipants(ctx context.Context, eventID uuid.UUID) ([]model.EventParticipant, error)
	AddChangeLog(ctx context.Context, eventID, changedBy uuid.UUID, field, oldVal, newVal string) error
}

type UserRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
}

type NotificationRepo interface {
	Create(ctx context.Context, userID uuid.UUID, eventID *uuid.UUID, notifType, message string) (model.Notification, error)
}

type Mailer interface {
	SendInvite(toEmail, eventTitle, organizerName, inviteToken string) error
}
