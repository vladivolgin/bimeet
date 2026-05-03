package notificationhandler_test

import (
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

	notificationhandler "bimeet/internal/handler/notification"
	"bimeet/internal/handler/notification/mock"
	"bimeet/internal/middleware"
	"bimeet/internal/model"
)

func newRouter(t *testing.T) (http.Handler, *mock.MockNotificationService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	svc := mock.NewMockNotificationService(ctrl)
	h := notificationhandler.New(svc)
	r := chi.NewRouter()
	r.Get("/notifications", h.List)
	r.Patch("/notifications/{id}/read", h.MarkRead)
	return r, svc
}

func withAuth(r *http.Request, userID uuid.UUID) *http.Request {
	return r.WithContext(middleware.WithUserID(r.Context(), userID))
}

// ─── List ──────────────────────────────────────────────────────────────────

func TestList_NoAuth(t *testing.T) {
	router, _ := newRouter(t)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/notifications", nil))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestList_ReturnsEmptySlice(t *testing.T) {
	router, svc := newRouter(t)
	userID := uuid.New()
	svc.EXPECT().List(gomock.Any(), userID).Return([]model.Notification{}, nil)

	req := withAuth(httptest.NewRequest(http.MethodGet, "/notifications", nil), userID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var body []interface{}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&body))
	assert.NotNil(t, body)
}

// ─── MarkRead ──────────────────────────────────────────────────────────────

func TestMarkRead_NoAuth(t *testing.T) {
	router, _ := newRouter(t)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodPatch, "/notifications/"+uuid.New().String()+"/read", nil))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestMarkRead_InvalidUUID(t *testing.T) {
	router, _ := newRouter(t)
	req := withAuth(httptest.NewRequest(http.MethodPatch, "/notifications/bad-id/read", nil), uuid.New())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestMarkRead_Forbidden(t *testing.T) {
	router, svc := newRouter(t)
	userID := uuid.New()
	notifID := uuid.New()
	svc.EXPECT().MarkRead(gomock.Any(), notifID, userID).Return(model.Notification{}, errors.New("forbidden"))

	req := withAuth(httptest.NewRequest(http.MethodPatch, "/notifications/"+notifID.String()+"/read", nil), userID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestMarkRead_NotFound(t *testing.T) {
	router, svc := newRouter(t)
	userID := uuid.New()
	notifID := uuid.New()
	svc.EXPECT().MarkRead(gomock.Any(), notifID, userID).Return(model.Notification{}, errors.New("not found"))

	req := withAuth(httptest.NewRequest(http.MethodPatch, "/notifications/"+notifID.String()+"/read", nil), userID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestMarkRead_HappyPath(t *testing.T) {
	router, svc := newRouter(t)
	userID := uuid.New()
	notifID := uuid.New()
	n := model.Notification{ID: notifID, UserID: userID, IsRead: true}
	svc.EXPECT().MarkRead(gomock.Any(), notifID, userID).Return(n, nil)

	req := withAuth(httptest.NewRequest(http.MethodPatch, "/notifications/"+notifID.String()+"/read", nil), userID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
