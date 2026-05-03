package pollhandler_test

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

	pollhandler "bimeet/internal/handler/poll"
	"bimeet/internal/handler/poll/mock"
	"bimeet/internal/middleware"
	"bimeet/internal/model"
)

func newRouter(t *testing.T) (http.Handler, *mock.MockPollService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	svc := mock.NewMockPollService(ctrl)
	h := pollhandler.New(svc)
	r := chi.NewRouter()
	r.Get("/events/{id}/polls", h.List)
	r.Post("/events/{id}/polls", h.Create)
	r.Post("/events/{id}/polls/{pollId}/vote", h.Vote)
	return r, svc
}

func withAuth(r *http.Request, userID uuid.UUID) *http.Request {
	return r.WithContext(middleware.WithUserID(r.Context(), userID))
}

// ─── List ──────────────────────────────────────────────────────────────────

func TestList_NoAuth(t *testing.T) {
	router, _ := newRouter(t)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, httptest.NewRequest(http.MethodGet, "/events/"+uuid.New().String()+"/polls", nil))
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestList_Forbidden(t *testing.T) {
	router, svc := newRouter(t)
	eventID, userID := uuid.New(), uuid.New()
	svc.EXPECT().List(gomock.Any(), eventID, userID).Return(nil, errors.New("forbidden"))

	req := withAuth(httptest.NewRequest(http.MethodGet, "/events/"+eventID.String()+"/polls", nil), userID)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestList_ReturnsEmptySlice(t *testing.T) {
	router, svc := newRouter(t)
	eventID, userID := uuid.New(), uuid.New()
	svc.EXPECT().List(gomock.Any(), eventID, userID).Return([]model.Poll{}, nil)

	req := withAuth(httptest.NewRequest(http.MethodGet, "/events/"+eventID.String()+"/polls", nil), userID)
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
	svc.EXPECT().Create(gomock.Any(), eventID, userID, gomock.Any()).Return(model.Poll{}, errors.New("forbidden"))

	body, _ := json.Marshal(model.CreatePollRequest{Question: "Q?", Options: []string{"A", "B"}})
	req := withAuth(httptest.NewRequest(http.MethodPost, "/events/"+eventID.String()+"/polls", bytes.NewBuffer(body)), userID)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestCreate_HappyPath(t *testing.T) {
	router, svc := newRouter(t)
	eventID, userID := uuid.New(), uuid.New()
	pollID := uuid.New()
	svc.EXPECT().Create(gomock.Any(), eventID, userID, gomock.Any()).
		Return(model.Poll{ID: pollID, EventID: eventID}, nil)

	body, _ := json.Marshal(model.CreatePollRequest{Question: "Where?", Options: []string{"Paris", "Berlin"}})
	req := withAuth(httptest.NewRequest(http.MethodPost, "/events/"+eventID.String()+"/polls", bytes.NewBuffer(body)), userID)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)
}

// ─── Vote ──────────────────────────────────────────────────────────────────

func TestVote_InvalidPollID(t *testing.T) {
	router, _ := newRouter(t)
	req := withAuth(httptest.NewRequest(http.MethodPost, "/events/"+uuid.New().String()+"/polls/bad/vote", nil), uuid.New())
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestVote_Forbidden(t *testing.T) {
	router, svc := newRouter(t)
	eventID, pollID, userID := uuid.New(), uuid.New(), uuid.New()
	svc.EXPECT().Vote(gomock.Any(), eventID, pollID, userID, gomock.Any()).Return(errors.New("forbidden"))

	body, _ := json.Marshal(model.VoteRequest{OptionID: uuid.New()})
	req := withAuth(httptest.NewRequest(http.MethodPost, "/events/"+eventID.String()+"/polls/"+pollID.String()+"/vote", bytes.NewBuffer(body)), userID)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusForbidden, w.Code)
}

func TestVote_HappyPath(t *testing.T) {
	router, svc := newRouter(t)
	eventID, pollID, userID := uuid.New(), uuid.New(), uuid.New()
	svc.EXPECT().Vote(gomock.Any(), eventID, pollID, userID, gomock.Any()).Return(nil)

	body, _ := json.Marshal(model.VoteRequest{OptionID: uuid.New()})
	req := withAuth(httptest.NewRequest(http.MethodPost, "/events/"+eventID.String()+"/polls/"+pollID.String()+"/vote", bytes.NewBuffer(body)), userID)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusNoContent, w.Code)
}
