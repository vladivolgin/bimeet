package collectionsvc_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"bimeet/internal/model"
	collectionsvc "bimeet/internal/service/collection"
	"bimeet/internal/service/collection/mock"
)

func newService(t *testing.T) (*collectionsvc.Service, *mock.MockCollectionRepo, *mock.MockEventRepo) {
	t.Helper()
	ctrl := gomock.NewController(t)
	cr := mock.NewMockCollectionRepo(ctrl)
	er := mock.NewMockEventRepo(ctrl)
	return collectionsvc.New(cr, er), cr, er
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

func TestList_HappyPath(t *testing.T) {
	svc, cr, er := newService(t)
	eventID, userID := uuid.New(), uuid.New()
	er.EXPECT().IsParticipant(gomock.Any(), eventID, userID).Return(true, nil)
	cr.EXPECT().ListByEvent(gomock.Any(), eventID).Return([]model.Collection{}, nil)

	result, err := svc.List(context.Background(), eventID, userID)
	require.NoError(t, err)
	assert.NotNil(t, result)
}

// ─── Create ────────────────────────────────────────────────────────────────

func TestCreate_NotOrganizer(t *testing.T) {
	svc, _, er := newService(t)
	eventID, userID := uuid.New(), uuid.New()
	er.EXPECT().GetByID(gomock.Any(), eventID).Return(model.Event{OrganizerID: uuid.New()}, nil)

	_, err := svc.Create(context.Background(), eventID, userID, model.CreateCollectionRequest{Title: "T"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestCreate_MissingTitle(t *testing.T) {
	svc, _, er := newService(t)
	organizerID := uuid.New()
	eventID := uuid.New()
	er.EXPECT().GetByID(gomock.Any(), eventID).Return(model.Event{OrganizerID: organizerID}, nil)

	_, err := svc.Create(context.Background(), eventID, organizerID, model.CreateCollectionRequest{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "title is required")
}

func TestCreate_HappyPath(t *testing.T) {
	svc, cr, er := newService(t)
	organizerID := uuid.New()
	eventID := uuid.New()
	req := model.CreateCollectionRequest{Title: "Venue", TargetAmount: 500}
	expected := model.Collection{ID: uuid.New(), EventID: eventID, Title: "Venue"}

	er.EXPECT().GetByID(gomock.Any(), eventID).Return(model.Event{OrganizerID: organizerID}, nil)
	cr.EXPECT().Create(gomock.Any(), eventID, req, organizerID).Return(expected, nil)

	result, err := svc.Create(context.Background(), eventID, organizerID, req)
	require.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
}

// ─── Delete ────────────────────────────────────────────────────────────────

func TestDelete_WithPaidContributions(t *testing.T) {
	svc, cr, er := newService(t)
	organizerID := uuid.New()
	eventID := uuid.New()
	collectionID := uuid.New()

	er.EXPECT().GetByID(gomock.Any(), eventID).Return(model.Event{OrganizerID: organizerID}, nil)
	cr.EXPECT().GetByID(gomock.Any(), collectionID).Return(model.Collection{ID: collectionID, EventID: eventID}, nil)
	cr.EXPECT().CountPaidContributions(gomock.Any(), collectionID).Return(2, nil)

	err := svc.Delete(context.Background(), eventID, collectionID, organizerID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "paid contributions")
}

func TestDelete_HappyPath(t *testing.T) {
	svc, cr, er := newService(t)
	organizerID := uuid.New()
	eventID := uuid.New()
	collectionID := uuid.New()

	er.EXPECT().GetByID(gomock.Any(), eventID).Return(model.Event{OrganizerID: organizerID}, nil)
	cr.EXPECT().GetByID(gomock.Any(), collectionID).Return(model.Collection{ID: collectionID, EventID: eventID}, nil)
	cr.EXPECT().CountPaidContributions(gomock.Any(), collectionID).Return(0, nil)
	cr.EXPECT().Delete(gomock.Any(), collectionID).Return(nil)

	require.NoError(t, svc.Delete(context.Background(), eventID, collectionID, organizerID))
}

// ─── ToggleContribution ────────────────────────────────────────────────────

func TestToggleContribution_NotConfirmed(t *testing.T) {
	svc, _, er := newService(t)
	eventID, collectionID, userID := uuid.New(), uuid.New(), uuid.New()
	er.EXPECT().IsConfirmedParticipant(gomock.Any(), eventID, userID).Return(false, nil)

	_, err := svc.ToggleContribution(context.Background(), eventID, collectionID, userID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestToggleContribution_WrongEvent(t *testing.T) {
	svc, cr, er := newService(t)
	eventID, collectionID, userID := uuid.New(), uuid.New(), uuid.New()
	er.EXPECT().IsConfirmedParticipant(gomock.Any(), eventID, userID).Return(true, nil)
	cr.EXPECT().GetByID(gomock.Any(), collectionID).Return(model.Collection{EventID: uuid.New()}, nil)

	_, err := svc.ToggleContribution(context.Background(), eventID, collectionID, userID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not belong")
}

// ─── Summary ───────────────────────────────────────────────────────────────

func TestSummary_Forbidden(t *testing.T) {
	svc, _, er := newService(t)
	eventID, userID := uuid.New(), uuid.New()
	er.EXPECT().IsParticipant(gomock.Any(), eventID, userID).Return(false, nil)

	_, err := svc.Summary(context.Background(), eventID, userID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestSummary_HappyPath(t *testing.T) {
	svc, cr, er := newService(t)
	eventID, userID := uuid.New(), uuid.New()
	colID := uuid.New()

	er.EXPECT().IsParticipant(gomock.Any(), eventID, userID).Return(true, nil)
	er.EXPECT().CountConfirmedParticipants(gomock.Any(), eventID).Return(4, nil)
	cr.EXPECT().ListByEvent(gomock.Any(), eventID).Return([]model.Collection{
		{ID: colID, EventID: eventID, Title: "Cake", TargetAmount: 200},
	}, nil)
	cr.EXPECT().CountPaidContributions(gomock.Any(), colID).Return(2, nil)

	result, err := svc.Summary(context.Background(), eventID, userID)
	require.NoError(t, err)
	assert.Len(t, result.Collections, 1)
	assert.Equal(t, float64(200), result.GrandTotal)
}
