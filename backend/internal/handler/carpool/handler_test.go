package carpoolhandler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	carpoolhandler "bimeet/internal/handler/carpool"
	"bimeet/internal/handler/carpool/mock"
	"bimeet/internal/middleware"
	"bimeet/internal/model"
)

func newRouter(t *testing.T) (http.Handler, *mock.MockCarpoolService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	svc := mock.NewMockCarpoolService(ctrl)
	h := carpoolhandler.New(svc)
	r := chi.NewRouter()
	r.Get("/events/{id}/carpools", h.List)
	r.Post("/events/{id}/carpools", h.Create)
	r.Post("/events/{id}/carpools/{carpoolId}/join", h.Join)
	return r, svc
}

func withAuth(r *http.Request, userID uuid.UUID) *http.Request {
	return r.WithContext(middleware.WithUserID(r.Context(), userID))
}

// ─── List ──────────────────────────────────────────────────────────────────

func TestList_NoAuth(t *testing.T) {
	router, _ := newRouter(t)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/events/"+uuid.New().String()+"/carpools", nil))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestList_Forbidden(t *testing.T) {
	router, svc := newRouter(t)
	eventID, userID := uuid.New(), uuid.New()
	svc.EXPECT().List(gomock.Any(), eventID, userID).Return(nil, errors.New("forbidden"))

	req := withAuth(httptest.NewRequest(http.MethodGet, "/events/"+eventID.String()+"/carpools", nil), userID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestList_ReturnsEmptySlice(t *testing.T) {
	router, svc := newRouter(t)
	eventID, userID := uuid.New(), uuid.New()
	svc.EXPECT().List(gomock.Any(), eventID, userID).Return([]model.Carpool{}, nil)

	req := withAuth(httptest.NewRequest(http.MethodGet, "/events/"+eventID.String()+"/carpools", nil), userID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
	var body []interface{}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&body))
	assert.NotNil(t, body)
}

// ─── Create ────────────────────────────────────────────────────────────────

func TestCreate_Forbidden(t *testing.T) {
	router, svc := newRouter(t)
	eventID, userID := uuid.New(), uuid.New()
	svc.EXPECT().Create(gomock.Any(), eventID, userID, gomock.Any()).Return(model.Carpool{}, errors.New("forbidden"))

	body, _ := json.Marshal(model.CreateCarpoolRequest{SeatsAvailable: 3})
	req := withAuth(httptest.NewRequest(http.MethodPost, "/events/"+eventID.String()+"/carpools", bytes.NewBuffer(body)), userID)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCreate_HappyPath(t *testing.T) {
	router, svc := newRouter(t)
	eventID, userID := uuid.New(), uuid.New()
	carpoolID := uuid.New()
	svc.EXPECT().Create(gomock.Any(), eventID, userID, gomock.Any()).
		Return(model.Carpool{ID: carpoolID, EventID: eventID, DriverID: userID}, nil)

	body, _ := json.Marshal(model.CreateCarpoolRequest{SeatsAvailable: 3})
	req := withAuth(httptest.NewRequest(http.MethodPost, "/events/"+eventID.String()+"/carpools", bytes.NewBuffer(body)), userID)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

// ─── Join ──────────────────────────────────────────────────────────────────

func TestJoin_InvalidCarpoolID(t *testing.T) {
	router, _ := newRouter(t)
	req := withAuth(httptest.NewRequest(http.MethodPost, "/events/"+uuid.New().String()+"/carpools/bad/join", nil), uuid.New())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestJoin_Full(t *testing.T) {
	router, svc := newRouter(t)
	eventID, carpoolID, userID := uuid.New(), uuid.New(), uuid.New()
	svc.EXPECT().Join(gomock.Any(), eventID, carpoolID, userID).Return(errors.New("carpool is full"))

	req := withAuth(httptest.NewRequest(http.MethodPost, "/events/"+eventID.String()+"/carpools/"+carpoolID.String()+"/join", nil), userID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestJoin_HappyPath(t *testing.T) {
	router, svc := newRouter(t)
	eventID, carpoolID, userID := uuid.New(), uuid.New(), uuid.New()
	svc.EXPECT().Join(gomock.Any(), eventID, carpoolID, userID).Return(nil)

	req := withAuth(httptest.NewRequest(http.MethodPost, "/events/"+eventID.String()+"/carpools/"+carpoolID.String()+"/join", nil), userID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
