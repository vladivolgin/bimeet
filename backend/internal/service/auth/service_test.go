package authsvc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"

	"bimeet/internal/model"
	authsvc "bimeet/internal/service/auth"
	"bimeet/internal/service/auth/mock"
)

func newService(t *testing.T) (*authsvc.Service, *mock.MockUserRepo) {
	t.Helper()
	ctrl := gomock.NewController(t)
	repo := mock.NewMockUserRepo(ctrl)
	tokens := mock.NewMockTokenRepo(ctrl)
	mailer := mock.NewMockMailer(ctrl)
	return authsvc.New(repo, tokens, mailer, "secret", 72), repo
}

func newServiceFull(t *testing.T) (*authsvc.Service, *mock.MockUserRepo, *mock.MockTokenRepo, *mock.MockMailer) {
	t.Helper()
	ctrl := gomock.NewController(t)
	repo := mock.NewMockUserRepo(ctrl)
	tokens := mock.NewMockTokenRepo(ctrl)
	mailer := mock.NewMockMailer(ctrl)
	return authsvc.New(repo, tokens, mailer, "secret", 72), repo, tokens, mailer
}

// ─── Register ──────────────────────────────────────────────────────────────

func TestRegister_MissingFields(t *testing.T) {
	svc, _ := newService(t)
	_, err := svc.Register(context.Background(), model.RegisterRequest{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "required")
}

func TestRegister_HappyPath(t *testing.T) {
	svc, repo := newService(t)

	userID := uuid.New()
	repo.EXPECT().Create(gomock.Any(), "Alice", "alice@example.com", gomock.Any()).
		Return(model.User{ID: userID, Name: "Alice", Email: "alice@example.com"}, nil)

	resp, err := svc.Register(context.Background(), model.RegisterRequest{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "pass123",
	})
	require.NoError(t, err)
	assert.Equal(t, userID, resp.User.ID)
	assert.NotEmpty(t, resp.Token)
}

func TestRegister_RepoError(t *testing.T) {
	svc, repo := newService(t)

	repo.EXPECT().Create(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
		Return(model.User{}, errors.New("duplicate key"))

	_, err := svc.Register(context.Background(), model.RegisterRequest{
		Name:     "Alice",
		Email:    "alice@example.com",
		Password: "pass123",
	})
	require.Error(t, err)
}

// ─── Login ─────────────────────────────────────────────────────────────────

func TestLogin_MissingFields(t *testing.T) {
	svc, _ := newService(t)
	_, err := svc.Login(context.Background(), model.LoginRequest{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "required")
}

func TestLogin_UserNotFound(t *testing.T) {
	svc, repo := newService(t)

	repo.EXPECT().GetByEmail(gomock.Any(), "ghost@x.com").Return(model.User{}, errors.New("not found"))

	_, err := svc.Login(context.Background(), model.LoginRequest{Email: "ghost@x.com", Password: "pass"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestLogin_WrongPassword(t *testing.T) {
	svc, repo := newService(t)

	hash, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.MinCost)
	repo.EXPECT().GetByEmail(gomock.Any(), "u@x.com").
		Return(model.User{Email: "u@x.com", PasswordHash: string(hash)}, nil)

	_, err := svc.Login(context.Background(), model.LoginRequest{Email: "u@x.com", Password: "wrong"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid credentials")
}

func TestLogin_HappyPath(t *testing.T) {
	svc, repo := newService(t)

	hash, _ := bcrypt.GenerateFromPassword([]byte("pass123"), bcrypt.MinCost)
	userID := uuid.New()
	repo.EXPECT().GetByEmail(gomock.Any(), "alice@example.com").
		Return(model.User{ID: userID, Email: "alice@example.com", PasswordHash: string(hash)}, nil)

	resp, err := svc.Login(context.Background(), model.LoginRequest{Email: "alice@example.com", Password: "pass123"})
	require.NoError(t, err)
	assert.Equal(t, userID, resp.User.ID)
	assert.NotEmpty(t, resp.Token)
}

// ─── ForgotPassword ────────────────────────────────────────────────────────

func TestForgotPassword_MissingEmail(t *testing.T) {
	svc, _, _, _ := newServiceFull(t)
	err := svc.ForgotPassword(context.Background(), model.ForgotPasswordRequest{})
	require.Error(t, err)
}

func TestForgotPassword_UserNotFound_ReturnsNil(t *testing.T) {
	svc, repo, _, _ := newServiceFull(t)
	repo.EXPECT().GetByEmail(gomock.Any(), "ghost@x.com").
		Return(model.User{}, errors.New("not found"))

	err := svc.ForgotPassword(context.Background(), model.ForgotPasswordRequest{Email: "ghost@x.com"})
	require.NoError(t, err)
}

func TestForgotPassword_HappyPath(t *testing.T) {
	svc, repo, tokens, mailer := newServiceFull(t)

	userID := uuid.New()
	tokenID := uuid.New()
	repo.EXPECT().GetByEmail(gomock.Any(), "alice@example.com").
		Return(model.User{ID: userID, Email: "alice@example.com"}, nil)
	tokens.EXPECT().Create(gomock.Any(), userID, gomock.Any()).
		Return(model.PasswordResetToken{Token: tokenID, UserID: userID}, nil)

	done := make(chan struct{})
	mailer.EXPECT().SendPasswordReset("alice@example.com", tokenID.String()).
		DoAndReturn(func(string, string) error {
			close(done)
			return nil
		})

	err := svc.ForgotPassword(context.Background(), model.ForgotPasswordRequest{Email: "alice@example.com"})
	require.NoError(t, err)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("SendPasswordReset was not called")
	}
}

// ─── ResetPassword ─────────────────────────────────────────────────────────

func TestResetPassword_MissingFields(t *testing.T) {
	svc, _, _, _ := newServiceFull(t)
	err := svc.ResetPassword(context.Background(), model.ResetPasswordRequest{})
	require.Error(t, err)
}

func TestResetPassword_InvalidTokenFormat(t *testing.T) {
	svc, _, _, _ := newServiceFull(t)
	err := svc.ResetPassword(context.Background(), model.ResetPasswordRequest{Token: "not-a-uuid", Password: "new"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestResetPassword_TokenUsed(t *testing.T) {
	svc, _, tokens, _ := newServiceFull(t)
	tokenID := uuid.New()
	tokens.EXPECT().GetByToken(gomock.Any(), tokenID).
		Return(model.PasswordResetToken{Token: tokenID, Used: true, ExpiresAt: time.Now().Add(time.Hour)}, nil)

	err := svc.ResetPassword(context.Background(), model.ResetPasswordRequest{Token: tokenID.String(), Password: "new"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestResetPassword_TokenExpired(t *testing.T) {
	svc, _, tokens, _ := newServiceFull(t)
	tokenID := uuid.New()
	tokens.EXPECT().GetByToken(gomock.Any(), tokenID).
		Return(model.PasswordResetToken{Token: tokenID, ExpiresAt: time.Now().Add(-time.Hour)}, nil)

	err := svc.ResetPassword(context.Background(), model.ResetPasswordRequest{Token: tokenID.String(), Password: "new"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "invalid token")
}

func TestResetPassword_HappyPath(t *testing.T) {
	svc, repo, tokens, _ := newServiceFull(t)

	tokenID := uuid.New()
	userID := uuid.New()
	tokens.EXPECT().GetByToken(gomock.Any(), tokenID).
		Return(model.PasswordResetToken{Token: tokenID, UserID: userID, ExpiresAt: time.Now().Add(time.Hour)}, nil)
	repo.EXPECT().UpdatePassword(gomock.Any(), userID, gomock.Any()).Return(nil)
	tokens.EXPECT().MarkUsed(gomock.Any(), tokenID).Return(nil)

	err := svc.ResetPassword(context.Background(), model.ResetPasswordRequest{Token: tokenID.String(), Password: "newpass123"})
	require.NoError(t, err)
}
