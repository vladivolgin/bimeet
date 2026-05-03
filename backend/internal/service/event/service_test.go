package eventsvc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"bimeet/internal/model"
	eventsvc "bimeet/internal/service/event"
	"bimeet/internal/service/event/mock"
)

// ─── Helpers ───────────────────────────────────────────────────────────────

func newService(t *testing.T) (*eventsvc.Service, *mock.MockEventRepo, *mock.MockUserRepo, *mock.MockNotificationRepo, *mock.MockMailer) {
	t.Helper()
	ctrl := gomock.NewController(t)
	er := mock.NewMockEventRepo(ctrl)
	ur := mock.NewMockUserRepo(ctrl)
	nr := mock.NewMockNotificationRepo(ctrl)
	ml := mock.NewMockMailer(ctrl)
	return eventsvc.New(er, ur, nr, ml), er, ur, nr, ml
}

// ─── Create ────────────────────────────────────────────────────────────────

func TestCreate_MissingTitle(t *testing.T) {
	svc, _, _, _, _ := newService(t)
	_, err := svc.Create(context.Background(), model.CreateEventRequest{}, uuid.New())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "title is required")
}

func TestCreate_MissingDates(t *testing.T) {
	svc, _, _, _, _ := newService(t)
	_, err := svc.Create(context.Background(), model.CreateEventRequest{Title: "Party"}, uuid.New())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "required")
}

func TestCreate_DateEndBeforeDateStart(t *testing.T) {
	svc, _, _, _, _ := newService(t)
	now := time.Now()
	req := model.CreateEventRequest{
		Title:     "Party",
		DateStart: now.Add(2 * time.Hour),
		DateEnd:   now.Add(1 * time.Hour),
	}
	_, err := svc.Create(context.Background(), req, uuid.New())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "date_end must be after date_start")
}

func TestCreate_InvalidCategory(t *testing.T) {
	svc, _, _, _, _ := newService(t)
	now := time.Now()
	req := model.CreateEventRequest{
		Title:     "Party",
		DateStart: now,
		DateEnd:   now.Add(1 * time.Hour),
		Category:  "invalid",
	}
	_, err := svc.Create(context.Background(), req, uuid.New())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "category must be")
}

func TestCreate_DefaultsToOrdinary(t *testing.T) {
	svc, er, _, _, _ := newService(t)

	organizerID := uuid.New()
	now := time.Now()
	req := model.CreateEventRequest{Title: "Party", DateStart: now, DateEnd: now.Add(2 * time.Hour)}
	expected := model.Event{ID: uuid.New(), Title: "Party"}

	er.EXPECT().Create(gomock.Any(), gomock.AssignableToTypeOf(model.CreateEventRequest{}), organizerID).
		DoAndReturn(func(_ context.Context, r model.CreateEventRequest, _ uuid.UUID) (model.Event, error) {
			assert.Equal(t, "ordinary", r.Category)
			return expected, nil
		})

	result, err := svc.Create(context.Background(), req, organizerID)
	require.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
}

// ─── Delete ────────────────────────────────────────────────────────────────

func TestDelete_NotFound(t *testing.T) {
	svc, er, _, _, _ := newService(t)

	id := uuid.New()
	er.EXPECT().GetByID(gomock.Any(), id).Return(model.Event{}, errors.New("no rows"))

	err := svc.Delete(context.Background(), id, uuid.New())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestDelete_Forbidden(t *testing.T) {
	svc, er, _, _, _ := newService(t)

	id := uuid.New()
	callerID := uuid.New()
	er.EXPECT().GetByID(gomock.Any(), id).Return(model.Event{ID: id, OrganizerID: uuid.New()}, nil)

	err := svc.Delete(context.Background(), id, callerID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestDelete_HappyPath(t *testing.T) {
	svc, er, _, _, _ := newService(t)

	id := uuid.New()
	organizerID := uuid.New()
	er.EXPECT().GetByID(gomock.Any(), id).Return(model.Event{ID: id, OrganizerID: organizerID}, nil)
	er.EXPECT().Delete(gomock.Any(), id).Return(nil)

	require.NoError(t, svc.Delete(context.Background(), id, organizerID))
}

// ─── InviteParticipant ─────────────────────────────────────────────────────

func TestInviteParticipant_Forbidden(t *testing.T) {
	svc, er, _, _, _ := newService(t)

	eventID := uuid.New()
	callerID := uuid.New()
	er.EXPECT().GetByID(gomock.Any(), eventID).Return(model.Event{ID: eventID, OrganizerID: uuid.New()}, nil)

	_, err := svc.InviteParticipant(context.Background(), eventID, callerID, model.InviteParticipantRequest{Email: "x@x.com"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestInviteParticipant_UnknownEmailSendsMailAsync(t *testing.T) {
	svc, er, ur, _, ml := newService(t)

	organizerID := uuid.New()
	eventID := uuid.New()
	inviteToken := uuid.New()
	event := model.Event{ID: eventID, OrganizerID: organizerID, Title: "Gala", InviteToken: inviteToken}
	organizer := model.User{ID: organizerID, Name: "Alice"}

	er.EXPECT().GetByID(gomock.Any(), eventID).Return(event, nil)
	ur.EXPECT().GetByID(gomock.Any(), organizerID).Return(organizer, nil)
	ur.EXPECT().GetByEmail(gomock.Any(), "unknown@x.com").Return(model.User{}, errors.New("not found"))
	ml.EXPECT().SendInvite("unknown@x.com", "Gala", "Alice", inviteToken.String()).Return(nil)

	p, err := svc.InviteParticipant(context.Background(), eventID, organizerID, model.InviteParticipantRequest{Email: "unknown@x.com"})
	require.NoError(t, err)
	assert.Equal(t, uuid.Nil, p.ID)

	// Wait briefly for goroutine
	time.Sleep(10 * time.Millisecond)
}

// ─── UpdateParticipantStatus ───────────────────────────────────────────────

func TestUpdateParticipantStatus_InvalidStatus(t *testing.T) {
	svc, _, _, _, _ := newService(t)
	id := uuid.New()
	_, err := svc.UpdateParticipantStatus(context.Background(), id, id, id, model.UpdateParticipantStatusRequest{Status: "maybe"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "status must be")
}

func TestUpdateParticipantStatus_Forbidden(t *testing.T) {
	svc, _, _, _, _ := newService(t)
	eventID := uuid.New()
	targetID := uuid.New()
	callerID := uuid.New() // different from targetID
	_, err := svc.UpdateParticipantStatus(context.Background(), eventID, targetID, callerID, model.UpdateParticipantStatusRequest{Status: "confirmed"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestUpdateParticipantStatus_HappyPath(t *testing.T) {
	svc, er, _, _, _ := newService(t)

	eventID := uuid.New()
	userID := uuid.New()
	participant := model.EventParticipant{EventID: eventID, UserID: userID, Status: "confirmed"}

	er.EXPECT().GetByID(gomock.Any(), eventID).Return(model.Event{ID: eventID}, nil)
	er.EXPECT().UpdateParticipantStatus(gomock.Any(), eventID, userID, "confirmed").Return(participant, nil)

	result, err := svc.UpdateParticipantStatus(context.Background(), eventID, userID, userID, model.UpdateParticipantStatusRequest{Status: "confirmed"})
	require.NoError(t, err)
	assert.Equal(t, "confirmed", result.Status)
}
