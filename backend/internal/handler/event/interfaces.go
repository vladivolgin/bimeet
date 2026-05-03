package eventhandler

//go:generate mockgen -destination=mock/mock.go -package=mock . EventService

import (
	"context"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

type EventService interface {
	List(ctx context.Context, userID uuid.UUID) ([]model.EventListItem, error)
	ListPublic(ctx context.Context, userID uuid.UUID) ([]model.PublicEventListItem, error)
	Create(ctx context.Context, req model.CreateEventRequest, organizerID uuid.UUID) (model.Event, error)
	GetByID(ctx context.Context, id, userID uuid.UUID) (model.EventDetail, error)
	GetByInviteToken(ctx context.Context, token uuid.UUID) (model.InviteEventInfo, error)
	JoinByInviteToken(ctx context.Context, token, userID uuid.UUID, req model.JoinByInviteTokenRequest) (model.EventParticipant, error)
	JoinPublic(ctx context.Context, eventID, userID uuid.UUID) error
	Update(ctx context.Context, id, userID uuid.UUID, req model.UpdateEventRequest) (model.Event, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
	Complete(ctx context.Context, id, userID uuid.UUID) error
	InviteParticipant(ctx context.Context, eventID, organizerID uuid.UUID, req model.InviteParticipantRequest) (model.EventParticipant, error)
	UpdateParticipantStatus(ctx context.Context, eventID, targetUserID, callerID uuid.UUID, req model.UpdateParticipantStatusRequest) (model.EventParticipant, error)
}
