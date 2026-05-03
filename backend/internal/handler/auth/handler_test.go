package authhandler_test

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

	authhandler "bimeet/internal/handler/auth"
	"bimeet/internal/handler/auth/mock"
	"bimeet/internal/model"
)

func newRouter(t *testing.T) (http.Handler, *mock.MockAuthService) {
	t.Helper()
	ctrl := gomock.NewController(t)
	svc := mock.NewMockAuthService(ctrl)
	h := authhandler.New(svc)
	r := chi.NewRouter()
	r.Post("/auth/register", h.Register)
	r.Post("/auth/login", h.Login)
	return r, svc
}

// ─── Register ──────────────────────────────────────────────────────────────

func TestRegister_InvalidBody(t *testing.T) {
	router, _ := newRouter(t)
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBufferString("bad"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_RequiredFields(t *testing.T) {
	router, svc := newRouter(t)
	svc.EXPECT().Register(gomock.Any(), gomock.Any()).Return(model.AuthResponse{}, errors.New("name, email and password are required"))

	body, _ := json.Marshal(model.RegisterRequest{})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestRegister_DuplicateEmail(t *testing.T) {
	router, svc := newRouter(t)
	svc.EXPECT().Register(gomock.Any(), gomock.Any()).Return(model.AuthResponse{}, errors.New("duplicate key"))

	body, _ := json.Marshal(model.RegisterRequest{Name: "A", Email: "a@b.com", Password: "p"})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusConflict, w.Code)
}

func TestRegister_HappyPath(t *testing.T) {
	router, svc := newRouter(t)
	userID := uuid.New()
	svc.EXPECT().Register(gomock.Any(), gomock.Any()).
		Return(model.AuthResponse{Token: "tok", User: model.User{ID: userID}}, nil)

	body, _ := json.Marshal(model.RegisterRequest{Name: "Alice", Email: "alice@x.com", Password: "pass"})
	req := httptest.NewRequest(http.MethodPost, "/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusCreated, w.Code)

	var resp model.AuthResponse
	require.NoError(t, json.NewDecoder(w.Body).Decode(&resp))
	assert.Equal(t, "tok", resp.Token)
}

// ─── Login ─────────────────────────────────────────────────────────────────

func TestLogin_InvalidCredentials(t *testing.T) {
	router, svc := newRouter(t)
	svc.EXPECT().Login(gomock.Any(), gomock.Any()).Return(model.AuthResponse{}, errors.New("invalid credentials"))

	body, _ := json.Marshal(model.LoginRequest{Email: "x@x.com", Password: "wrong"})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusUnauthorized, w.Code)
}

func TestLogin_HappyPath(t *testing.T) {
	router, svc := newRouter(t)
	userID := uuid.New()
	svc.EXPECT().Login(gomock.Any(), gomock.Any()).
		Return(model.AuthResponse{Token: "tok", User: model.User{ID: userID}}, nil)

	body, _ := json.Marshal(model.LoginRequest{Email: "alice@x.com", Password: "pass"})
	req := httptest.NewRequest(http.MethodPost, "/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}
