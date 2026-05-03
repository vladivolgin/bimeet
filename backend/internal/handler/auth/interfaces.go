package authhandler

//go:generate mockgen -destination=mock/mock.go -package=mock . AuthService

import (
	"context"

	"bimeet/internal/model"
)

type AuthService interface {
	Register(ctx context.Context, req model.RegisterRequest) (model.AuthResponse, error)
	Login(ctx context.Context, req model.LoginRequest) (model.AuthResponse, error)
	ForgotPassword(ctx context.Context, req model.ForgotPasswordRequest) error
	ResetPassword(ctx context.Context, req model.ResetPasswordRequest) error
}
