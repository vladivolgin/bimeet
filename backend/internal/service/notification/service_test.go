package notificationsvc_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"bimeet/internal/model"
	notificationsvc "bimeet/internal/service/notification"
	"bimeet/internal/service/notification/mock"
)

func newService(t *testing.T) (*notificationsvc.Service, *mock.MockNotificationRepo) {
	t.Helper()
	ctrl := gomock.NewController(t)
	repo := mock.NewMockNotificationRepo(ctrl)
	return notificationsvc.New(repo), repo
}

func TestList_HappyPath(t *testing.T) {
	svc, repo := newService(t)
	userID := uuid.New()
	expected := []model.Notification{{ID: uuid.New(), UserID: userID}}
	repo.EXPECT().ListForUser(gomock.Any(), userID).Return(expected, nil)

	result, err := svc.List(context.Background(), userID)
	require.NoError(t, err)
	assert.Len(t, result, 1)
}

func TestMarkRead_NotFound(t *testing.T) {
	svc, repo := newService(t)
	notifID := uuid.New()
	repo.EXPECT().GetByID(gomock.Any(), notifID).Return(model.Notification{}, errors.New("not found"))

	_, err := svc.MarkRead(context.Background(), notifID, uuid.New())
	require.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestMarkRead_Forbidden(t *testing.T) {
	svc, repo := newService(t)
	notifID := uuid.New()
	ownerID := uuid.New()
	callerID := uuid.New()
	repo.EXPECT().GetByID(gomock.Any(), notifID).Return(model.Notification{ID: notifID, UserID: ownerID}, nil)

	_, err := svc.MarkRead(context.Background(), notifID, callerID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestMarkRead_HappyPath(t *testing.T) {
	svc, repo := newService(t)
	notifID := uuid.New()
	userID := uuid.New()
	n := model.Notification{ID: notifID, UserID: userID, IsRead: true}
	repo.EXPECT().GetByID(gomock.Any(), notifID).Return(model.Notification{ID: notifID, UserID: userID}, nil)
	repo.EXPECT().MarkRead(gomock.Any(), notifID).Return(n, nil)

	result, err := svc.MarkRead(context.Background(), notifID, userID)
	require.NoError(t, err)
	assert.True(t, result.IsRead)
}
