package eventlinkhandler_test

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

	eventlinkhandler "bimeet/internal/handler/eventlink"
	"bimeet/internal/handler/eventlink/mock"
	"bimeet/internal/middleware"
	"bimeet/internal/model"
)

func newRouter(t *testing.T) (http.Handler, *mock.MockEventLinkService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	svc := mock.NewMockEventLinkService(ctrl)
	h := eventlinkhandler.New(svc)
	r := chi.NewRouter()
	r.Get("/events/{id}/links", h.List)
	r.Post("/events/{id}/links", h.Create)
	r.Delete("/events/{id}/links/{linkId}", h.Delete)
	return r, svc
}

func withAuth(r *http.Request, userID uuid.UUID) *http.Request {
	return r.WithContext(middleware.WithUserID(r.Context(), userID))
}

// ─── List ──────────────────────────────────────────────────────────────────

func TestList_NoAuth(t *testing.T) {
	router, _ := newRouter(t)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/events/"+uuid.New().String()+"/links", nil))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestList_Forbidden(t *testing.T) {
	router, svc := newRouter(t)
	eventID, userID := uuid.New(), uuid.New()
	svc.EXPECT().List(gomock.Any(), eventID, userID).Return(nil, errors.New("forbidden"))

	req := withAuth(httptest.NewRequest(http.MethodGet, "/events/"+eventID.String()+"/links", nil), userID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestList_HappyPath(t *testing.T) {
	router, svc := newRouter(t)
	eventID, userID := uuid.New(), uuid.New()
	svc.EXPECT().List(gomock.Any(), eventID, userID).Return([]model.EventLink{}, nil)

	req := withAuth(httptest.NewRequest(http.MethodGet, "/events/"+eventID.String()+"/links", nil), userID)
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
	svc.EXPECT().Create(gomock.Any(), eventID, userID, gomock.Any()).Return(model.EventLink{}, errors.New("forbidden"))

	body, _ := json.Marshal(model.CreateEventLinkRequest{Title: "Zoom", URL: "https://zoom.us"})
	req := withAuth(httptest.NewRequest(http.MethodPost, "/events/"+eventID.String()+"/links", bytes.NewBuffer(body)), userID)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCreate_BusinessOnly(t *testing.T) {
	router, svc := newRouter(t)
	eventID, userID := uuid.New(), uuid.New()
	svc.EXPECT().Create(gomock.Any(), eventID, userID, gomock.Any()).
		Return(model.EventLink{}, errors.New("links are only available for business events"))

	body, _ := json.Marshal(model.CreateEventLinkRequest{Title: "Zoom", URL: "https://zoom.us"})
	req := withAuth(httptest.NewRequest(http.MethodPost, "/events/"+eventID.String()+"/links", bytes.NewBuffer(body)), userID)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestCreate_HappyPath(t *testing.T) {
	router, svc := newRouter(t)
	eventID, userID := uuid.New(), uuid.New()
	linkID := uuid.New()
	svc.EXPECT().Create(gomock.Any(), eventID, userID, gomock.Any()).
		Return(model.EventLink{ID: linkID, EventID: eventID}, nil)

	body, _ := json.Marshal(model.CreateEventLinkRequest{Title: "Zoom", URL: "https://zoom.us"})
	req := withAuth(httptest.NewRequest(http.MethodPost, "/events/"+eventID.String()+"/links", bytes.NewBuffer(body)), userID)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

// ─── Delete ────────────────────────────────────────────────────────────────

func TestDelete_InvalidLinkID(t *testing.T) {
	router, _ := newRouter(t)
	req := withAuth(httptest.NewRequest(http.MethodDelete, "/events/"+uuid.New().String()+"/links/bad", nil), uuid.New())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestDelete_Forbidden(t *testing.T) {
	router, svc := newRouter(t)
	eventID, linkID, userID := uuid.New(), uuid.New(), uuid.New()
	svc.EXPECT().Delete(gomock.Any(), eventID, linkID, userID).Return(errors.New("forbidden"))

	req := withAuth(httptest.NewRequest(http.MethodDelete, "/events/"+eventID.String()+"/links/"+linkID.String(), nil), userID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestDelete_HappyPath(t *testing.T) {
	router, svc := newRouter(t)
	eventID, linkID, userID := uuid.New(), uuid.New(), uuid.New()
	svc.EXPECT().Delete(gomock.Any(), eventID, linkID, userID).Return(nil)

	req := withAuth(httptest.NewRequest(http.MethodDelete, "/events/"+eventID.String()+"/links/"+linkID.String(), nil), userID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
