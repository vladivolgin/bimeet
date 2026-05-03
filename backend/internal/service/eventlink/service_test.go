package eventlinksvc_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"bimeet/internal/model"
	eventlinksvc "bimeet/internal/service/eventlink"
	"bimeet/internal/service/eventlink/mock"
)

func newService(t *testing.T) (*eventlinksvc.Service, *mock.MockEventLinkRepo, *mock.MockEventRepo) {
	t.Helper()
	ctrl := gomock.NewController(t)
	lr := mock.NewMockEventLinkRepo(ctrl)
	er := mock.NewMockEventRepo(ctrl)
	return eventlinksvc.New(lr, er), lr, er
}

// ─── List ──────────────────────────────────────────────────────────────────

func TestList_Forbidden(t *testing.T) {
	svc, _, er := newService(t)
	eventID, userID := uuid.New(), uuid.New()
	er.EXPECT().IsParticipant(gomock.Any(), eventID, userID).Return(false, nil)

	_, err := svc.List(context.Background(), eventID, userID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestList_ReturnsEmptySlice(t *testing.T) {
	svc, lr, er := newService(t)
	eventID, userID := uuid.New(), uuid.New()
	er.EXPECT().IsParticipant(gomock.Any(), eventID, userID).Return(true, nil)
	lr.EXPECT().List(gomock.Any(), eventID).Return(nil, nil)

	result, err := svc.List(context.Background(), eventID, userID)
	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Empty(t, result)
}

// ─── Create ────────────────────────────────────────────────────────────────

func TestCreate_MissingTitle(t *testing.T) {
	svc, _, _ := newService(t)
	_, err := svc.Create(context.Background(), uuid.New(), uuid.New(), model.CreateEventLinkRequest{URL: "https://x.com"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "title is required")
}

func TestCreate_MissingURL(t *testing.T) {
	svc, _, _ := newService(t)
	_, err := svc.Create(context.Background(), uuid.New(), uuid.New(), model.CreateEventLinkRequest{Title: "Agenda"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "url is required")
}

func TestCreate_NotOrdinaryEvent(t *testing.T) {
	svc, _, er := newService(t)
	organizerID := uuid.New()
	eventID := uuid.New()
	er.EXPECT().GetByID(gomock.Any(), eventID).Return(model.Event{OrganizerID: organizerID, Category: "ordinary"}, nil)

	_, err := svc.Create(context.Background(), eventID, organizerID, model.CreateEventLinkRequest{Title: "Zoom", URL: "https://zoom.us"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "only available for business")
}

func TestCreate_NotOrganizer(t *testing.T) {
	svc, _, er := newService(t)
	eventID, callerID := uuid.New(), uuid.New()
	er.EXPECT().GetByID(gomock.Any(), eventID).Return(model.Event{OrganizerID: uuid.New(), Category: "business"}, nil)

	_, err := svc.Create(context.Background(), eventID, callerID, model.CreateEventLinkRequest{Title: "Zoom", URL: "https://zoom.us"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestCreate_HappyPath(t *testing.T) {
	svc, lr, er := newService(t)
	organizerID := uuid.New()
	eventID := uuid.New()
	req := model.CreateEventLinkRequest{Title: "Zoom", URL: "https://zoom.us"}
	expected := model.EventLink{ID: uuid.New(), EventID: eventID}

	er.EXPECT().GetByID(gomock.Any(), eventID).Return(model.Event{OrganizerID: organizerID, Category: "business"}, nil)
	lr.EXPECT().Create(gomock.Any(), eventID, organizerID, req).Return(expected, nil)

	result, err := svc.Create(context.Background(), eventID, organizerID, req)
	require.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
}

// ─── Delete ────────────────────────────────────────────────────────────────

func TestDelete_Forbidden(t *testing.T) {
	svc, _, er := newService(t)
	eventID, linkID, callerID := uuid.New(), uuid.New(), uuid.New()
	er.EXPECT().GetByID(gomock.Any(), eventID).Return(model.Event{OrganizerID: uuid.New()}, nil)

	err := svc.Delete(context.Background(), eventID, linkID, callerID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestDelete_HappyPath(t *testing.T) {
	svc, lr, er := newService(t)
	organizerID := uuid.New()
	eventID, linkID := uuid.New(), uuid.New()

	er.EXPECT().GetByID(gomock.Any(), eventID).Return(model.Event{OrganizerID: organizerID}, nil)
	lr.EXPECT().Delete(gomock.Any(), linkID, eventID).Return(nil)

	require.NoError(t, svc.Delete(context.Background(), eventID, linkID, organizerID))
}
