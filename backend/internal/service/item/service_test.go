package itemsvc_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"bimeet/internal/model"
	itemsvc "bimeet/internal/service/item"
	"bimeet/internal/service/item/mock"
)

func newService(t *testing.T) (*itemsvc.Service, *mock.MockItemRepo, *mock.MockEventRepo) {
	t.Helper()
	ctrl := gomock.NewController(t)
	ir := mock.NewMockItemRepo(ctrl)
	er := mock.NewMockEventRepo(ctrl)
	return itemsvc.New(ir, er), ir, er
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
	svc, ir, er := newService(t)
	eventID, userID := uuid.New(), uuid.New()
	er.EXPECT().IsParticipant(gomock.Any(), eventID, userID).Return(true, nil)
	ir.EXPECT().ListByEvent(gomock.Any(), eventID).Return([]model.Item{}, nil)

	_, err := svc.List(context.Background(), eventID, userID)
	require.NoError(t, err)
}

// ─── Create ────────────────────────────────────────────────────────────────

func TestCreate_NotParticipant(t *testing.T) {
	svc, _, er := newService(t)
	eventID, userID := uuid.New(), uuid.New()
	er.EXPECT().IsParticipant(gomock.Any(), eventID, userID).Return(false, nil)

	_, err := svc.Create(context.Background(), eventID, userID, model.CreateItemRequest{Name: "Tent"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestCreate_MissingName(t *testing.T) {
	svc, _, er := newService(t)
	eventID, userID := uuid.New(), uuid.New()
	er.EXPECT().IsParticipant(gomock.Any(), eventID, userID).Return(true, nil)

	_, err := svc.Create(context.Background(), eventID, userID, model.CreateItemRequest{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "name is required")
}

func TestCreate_HappyPath(t *testing.T) {
	svc, ir, er := newService(t)
	eventID, userID := uuid.New(), uuid.New()
	req := model.CreateItemRequest{Name: "Tent"}
	expected := model.Item{ID: uuid.New(), EventID: eventID, Name: "Tent"}

	er.EXPECT().IsParticipant(gomock.Any(), eventID, userID).Return(true, nil)
	ir.EXPECT().Create(gomock.Any(), eventID, req).Return(expected, nil)

	result, err := svc.Create(context.Background(), eventID, userID, req)
	require.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
}

// ─── UpdateAssignment ──────────────────────────────────────────────────────

func TestUpdateAssignment_NotConfirmed(t *testing.T) {
	svc, _, er := newService(t)
	eventID, itemID, userID := uuid.New(), uuid.New(), uuid.New()
	er.EXPECT().IsConfirmedParticipant(gomock.Any(), eventID, userID).Return(false, nil)

	_, err := svc.UpdateAssignment(context.Background(), eventID, itemID, userID, model.UpdateItemRequest{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestUpdateAssignment_ItemWrongEvent(t *testing.T) {
	svc, ir, er := newService(t)
	eventID, itemID, userID := uuid.New(), uuid.New(), uuid.New()
	er.EXPECT().IsConfirmedParticipant(gomock.Any(), eventID, userID).Return(true, nil)
	ir.EXPECT().GetByID(gomock.Any(), itemID).Return(model.Item{EventID: uuid.New()}, nil)

	_, err := svc.UpdateAssignment(context.Background(), eventID, itemID, userID, model.UpdateItemRequest{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not belong")
}

func TestUpdateAssignment_AssignToSelf(t *testing.T) {
	svc, ir, er := newService(t)
	eventID, itemID, userID := uuid.New(), uuid.New(), uuid.New()
	er.EXPECT().IsConfirmedParticipant(gomock.Any(), eventID, userID).Return(true, nil)
	ir.EXPECT().GetByID(gomock.Any(), itemID).Return(model.Item{ID: itemID, EventID: eventID}, nil)
	ir.EXPECT().UpdateAssignment(gomock.Any(), itemID, &userID).Return(model.Item{ID: itemID, AssignedTo: &userID}, nil)

	result, err := svc.UpdateAssignment(context.Background(), eventID, itemID, userID, model.UpdateItemRequest{AssignedTo: &userID})
	require.NoError(t, err)
	assert.Equal(t, userID, *result.AssignedTo)
}
