package carpoolsvc_test

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"bimeet/internal/model"
	carpoolsvc "bimeet/internal/service/carpool"
	"bimeet/internal/service/carpool/mock"
)

func newService(t *testing.T) (*carpoolsvc.Service, *mock.MockCarpoolRepo, *mock.MockEventRepo) {
	t.Helper()
	ctrl := gomock.NewController(t)
	cr := mock.NewMockCarpoolRepo(ctrl)
	er := mock.NewMockEventRepo(ctrl)
	return carpoolsvc.New(cr, er), cr, er
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

// ─── Create ────────────────────────────────────────────────────────────────

func TestCreate_NotConfirmed(t *testing.T) {
	svc, _, er := newService(t)
	eventID, driverID := uuid.New(), uuid.New()
	er.EXPECT().IsConfirmedParticipant(gomock.Any(), eventID, driverID).Return(false, nil)

	_, err := svc.Create(context.Background(), eventID, driverID, model.CreateCarpoolRequest{SeatsAvailable: 3})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestCreate_NoSeats(t *testing.T) {
	svc, _, er := newService(t)
	eventID, driverID := uuid.New(), uuid.New()
	er.EXPECT().IsConfirmedParticipant(gomock.Any(), eventID, driverID).Return(true, nil)

	_, err := svc.Create(context.Background(), eventID, driverID, model.CreateCarpoolRequest{SeatsAvailable: 0})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "seats_available must be at least 1")
}

func TestCreate_HappyPath(t *testing.T) {
	svc, cr, er := newService(t)
	eventID, driverID := uuid.New(), uuid.New()
	req := model.CreateCarpoolRequest{SeatsAvailable: 3}
	expected := model.Carpool{ID: uuid.New(), EventID: eventID, DriverID: driverID}

	er.EXPECT().IsConfirmedParticipant(gomock.Any(), eventID, driverID).Return(true, nil)
	cr.EXPECT().Create(gomock.Any(), eventID, driverID, req).Return(expected, nil)

	result, err := svc.Create(context.Background(), eventID, driverID, req)
	require.NoError(t, err)
	assert.Equal(t, expected.ID, result.ID)
}

// ─── Join ──────────────────────────────────────────────────────────────────

func TestJoin_NotConfirmed(t *testing.T) {
	svc, _, er := newService(t)
	eventID, carpoolID, userID := uuid.New(), uuid.New(), uuid.New()
	er.EXPECT().IsConfirmedParticipant(gomock.Any(), eventID, userID).Return(false, nil)

	err := svc.Join(context.Background(), eventID, carpoolID, userID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "forbidden")
}

func TestJoin_WrongEvent(t *testing.T) {
	svc, cr, er := newService(t)
	eventID, carpoolID, userID := uuid.New(), uuid.New(), uuid.New()
	er.EXPECT().IsConfirmedParticipant(gomock.Any(), eventID, userID).Return(true, nil)
	cr.EXPECT().GetByID(gomock.Any(), carpoolID).Return(model.Carpool{EventID: uuid.New()}, nil)

	err := svc.Join(context.Background(), eventID, carpoolID, userID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "does not belong")
}

func TestJoin_IsDriver(t *testing.T) {
	svc, cr, er := newService(t)
	eventID, carpoolID, userID := uuid.New(), uuid.New(), uuid.New()
	er.EXPECT().IsConfirmedParticipant(gomock.Any(), eventID, userID).Return(true, nil)
	cr.EXPECT().GetByID(gomock.Any(), carpoolID).Return(model.Carpool{EventID: eventID, DriverID: userID, SeatsAvailable: 3}, nil)

	err := svc.Join(context.Background(), eventID, carpoolID, userID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "driver")
}

func TestJoin_Full(t *testing.T) {
	svc, cr, er := newService(t)
	eventID, carpoolID, userID := uuid.New(), uuid.New(), uuid.New()
	er.EXPECT().IsConfirmedParticipant(gomock.Any(), eventID, userID).Return(true, nil)
	cr.EXPECT().GetByID(gomock.Any(), carpoolID).Return(model.Carpool{EventID: eventID, DriverID: uuid.New(), SeatsAvailable: 2}, nil)
	cr.EXPECT().CountPassengers(gomock.Any(), carpoolID).Return(2, nil)

	err := svc.Join(context.Background(), eventID, carpoolID, userID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "full")
}

func TestJoin_AlreadyJoined(t *testing.T) {
	svc, cr, er := newService(t)
	eventID, carpoolID, userID := uuid.New(), uuid.New(), uuid.New()
	er.EXPECT().IsConfirmedParticipant(gomock.Any(), eventID, userID).Return(true, nil)
	cr.EXPECT().GetByID(gomock.Any(), carpoolID).Return(model.Carpool{EventID: eventID, DriverID: uuid.New(), SeatsAvailable: 3}, nil)
	cr.EXPECT().CountPassengers(gomock.Any(), carpoolID).Return(1, nil)
	cr.EXPECT().IsPassenger(gomock.Any(), carpoolID, userID).Return(true, nil)

	err := svc.Join(context.Background(), eventID, carpoolID, userID)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "already joined")
}

func TestJoin_HappyPath(t *testing.T) {
	svc, cr, er := newService(t)
	eventID, carpoolID, userID := uuid.New(), uuid.New(), uuid.New()
	er.EXPECT().IsConfirmedParticipant(gomock.Any(), eventID, userID).Return(true, nil)
	cr.EXPECT().GetByID(gomock.Any(), carpoolID).Return(model.Carpool{EventID: eventID, DriverID: uuid.New(), SeatsAvailable: 3}, nil)
	cr.EXPECT().CountPassengers(gomock.Any(), carpoolID).Return(1, nil)
	cr.EXPECT().IsPassenger(gomock.Any(), carpoolID, userID).Return(false, nil)
	cr.EXPECT().AddPassenger(gomock.Any(), carpoolID, userID).Return(nil)

	require.NoError(t, svc.Join(context.Background(), eventID, carpoolID, userID))
}
