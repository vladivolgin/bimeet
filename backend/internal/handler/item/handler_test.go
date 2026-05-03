package itemhandler_test

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

	itemhandler "bimeet/internal/handler/item"
	"bimeet/internal/handler/item/mock"
	"bimeet/internal/middleware"
	"bimeet/internal/model"
)

func newRouter(t *testing.T) (http.Handler, *mock.MockItemService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	svc := mock.NewMockItemService(ctrl)
	h := itemhandler.New(svc)
	r := chi.NewRouter()
	r.Get("/events/{id}/items", h.List)
	r.Post("/events/{id}/items", h.Create)
	r.Patch("/events/{id}/items/{itemId}", h.UpdateAssignment)
	return r, svc
}

func withAuth(r *http.Request, userID uuid.UUID) *http.Request {
	return r.WithContext(middleware.WithUserID(r.Context(), userID))
}

// ─── List ──────────────────────────────────────────────────────────────────

func TestList_NoAuth(t *testing.T) {
	router, _ := newRouter(t)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/events/"+uuid.New().String()+"/items", nil))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestList_Forbidden(t *testing.T) {
	router, svc := newRouter(t)
	eventID, userID := uuid.New(), uuid.New()
	svc.EXPECT().List(gomock.Any(), eventID, userID).Return(nil, errors.New("forbidden"))

	req := withAuth(httptest.NewRequest(http.MethodGet, "/events/"+eventID.String()+"/items", nil), userID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestList_ReturnsEmptySlice(t *testing.T) {
	router, svc := newRouter(t)
	eventID, userID := uuid.New(), uuid.New()
	svc.EXPECT().List(gomock.Any(), eventID, userID).Return([]model.Item{}, nil)

	req := withAuth(httptest.NewRequest(http.MethodGet, "/events/"+eventID.String()+"/items", nil), userID)
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
	svc.EXPECT().Create(gomock.Any(), eventID, userID, gomock.Any()).Return(model.Item{}, errors.New("forbidden"))

	body, _ := json.Marshal(model.CreateItemRequest{Name: "Tent"})
	req := withAuth(httptest.NewRequest(http.MethodPost, "/events/"+eventID.String()+"/items", bytes.NewBuffer(body)), userID)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCreate_HappyPath(t *testing.T) {
	router, svc := newRouter(t)
	eventID, userID := uuid.New(), uuid.New()
	itemID := uuid.New()
	svc.EXPECT().Create(gomock.Any(), eventID, userID, gomock.Any()).
		Return(model.Item{ID: itemID, EventID: eventID, Name: "Tent"}, nil)

	body, _ := json.Marshal(model.CreateItemRequest{Name: "Tent"})
	req := withAuth(httptest.NewRequest(http.MethodPost, "/events/"+eventID.String()+"/items", bytes.NewBuffer(body)), userID)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

// ─── UpdateAssignment ──────────────────────────────────────────────────────

func TestUpdateAssignment_InvalidItemID(t *testing.T) {
	router, _ := newRouter(t)
	req := withAuth(httptest.NewRequest(http.MethodPatch, "/events/"+uuid.New().String()+"/items/bad", nil), uuid.New())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestUpdateAssignment_Forbidden(t *testing.T) {
	router, svc := newRouter(t)
	eventID, itemID, userID := uuid.New(), uuid.New(), uuid.New()
	svc.EXPECT().UpdateAssignment(gomock.Any(), eventID, itemID, userID, gomock.Any()).
		Return(model.Item{}, errors.New("forbidden"))

	body, _ := json.Marshal(model.UpdateItemRequest{})
	req := withAuth(httptest.NewRequest(http.MethodPatch, "/events/"+eventID.String()+"/items/"+itemID.String(), bytes.NewBuffer(body)), userID)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestUpdateAssignment_HappyPath(t *testing.T) {
	router, svc := newRouter(t)
	eventID, itemID, userID := uuid.New(), uuid.New(), uuid.New()
	svc.EXPECT().UpdateAssignment(gomock.Any(), eventID, itemID, userID, gomock.Any()).
		Return(model.Item{ID: itemID, AssignedTo: &userID}, nil)

	body, _ := json.Marshal(model.UpdateItemRequest{AssignedTo: &userID})
	req := withAuth(httptest.NewRequest(http.MethodPatch, "/events/"+eventID.String()+"/items/"+itemID.String(), bytes.NewBuffer(body)), userID)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
