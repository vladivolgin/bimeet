package collectionhandler_test

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

	collectionhandler "bimeet/internal/handler/collection"
	"bimeet/internal/handler/collection/mock"
	"bimeet/internal/middleware"
	"bimeet/internal/model"
)

func newRouter(t *testing.T) (http.Handler, *mock.MockCollectionService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	svc := mock.NewMockCollectionService(ctrl)
	h := collectionhandler.New(svc)
	r := chi.NewRouter()
	r.Get("/events/{id}/collections", h.List)
	r.Post("/events/{id}/collections", h.Create)
	r.Get("/events/{id}/collections/summary", h.Summary)
	r.Delete("/events/{id}/collections/{collectionId}", h.Delete)
	r.Patch("/events/{id}/collections/{collectionId}/contribute", h.ToggleContribution)
	return r, svc
}

func withAuth(r *http.Request, userID uuid.UUID) *http.Request {
	return r.WithContext(middleware.WithUserID(r.Context(), userID))
}

// ─── List ──────────────────────────────────────────────────────────────────

func TestList_NoAuth(t *testing.T) {
	router, _ := newRouter(t)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/events/"+uuid.New().String()+"/collections", nil))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestList_InvalidEventID(t *testing.T) {
	router, _ := newRouter(t)
	req := withAuth(httptest.NewRequest(http.MethodGet, "/events/bad/collections", nil), uuid.New())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestList_Forbidden(t *testing.T) {
	router, svc := newRouter(t)
	eventID, userID := uuid.New(), uuid.New()
	svc.EXPECT().List(gomock.Any(), eventID, userID).Return(nil, errors.New("forbidden"))

	req := withAuth(httptest.NewRequest(http.MethodGet, "/events/"+eventID.String()+"/collections", nil), userID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestList_ReturnsEmptySlice(t *testing.T) {
	router, svc := newRouter(t)
	eventID, userID := uuid.New(), uuid.New()
	svc.EXPECT().List(gomock.Any(), eventID, userID).Return([]model.Collection{}, nil)

	req := withAuth(httptest.NewRequest(http.MethodGet, "/events/"+eventID.String()+"/collections", nil), userID)
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
	svc.EXPECT().Create(gomock.Any(), eventID, userID, gomock.Any()).Return(model.Collection{}, errors.New("forbidden"))

	body, _ := json.Marshal(model.CreateCollectionRequest{Title: "T", TargetAmount: 100})
	req := withAuth(httptest.NewRequest(http.MethodPost, "/events/"+eventID.String()+"/collections", bytes.NewBuffer(body)), userID)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCreate_HappyPath(t *testing.T) {
	router, svc := newRouter(t)
	eventID, userID := uuid.New(), uuid.New()
	colID := uuid.New()
	svc.EXPECT().Create(gomock.Any(), eventID, userID, gomock.Any()).
		Return(model.Collection{ID: colID, EventID: eventID, Title: "Cake"}, nil)

	body, _ := json.Marshal(model.CreateCollectionRequest{Title: "Cake", TargetAmount: 200})
	req := withAuth(httptest.NewRequest(http.MethodPost, "/events/"+eventID.String()+"/collections", bytes.NewBuffer(body)), userID)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

// ─── Delete ────────────────────────────────────────────────────────────────

func TestDelete_HappyPath(t *testing.T) {
	router, svc := newRouter(t)
	eventID, colID, userID := uuid.New(), uuid.New(), uuid.New()
	svc.EXPECT().Delete(gomock.Any(), eventID, colID, userID).Return(nil)

	req := withAuth(httptest.NewRequest(http.MethodDelete, "/events/"+eventID.String()+"/collections/"+colID.String(), nil), userID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}

// ─── Summary ───────────────────────────────────────────────────────────────

func TestSummary_Forbidden(t *testing.T) {
	router, svc := newRouter(t)
	eventID, userID := uuid.New(), uuid.New()
	svc.EXPECT().Summary(gomock.Any(), eventID, userID).Return(model.EventCollectionsSummaryResponse{}, errors.New("forbidden"))

	req := withAuth(httptest.NewRequest(http.MethodGet, "/events/"+eventID.String()+"/collections/summary", nil), userID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}
