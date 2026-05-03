package authsvc

//go:generate mockgen -destination=mock/mock.go -package=mock . UserRepo,TokenRepo,Mailer

import (
	"context"
	"time"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

type UserRepo interface {
	Create(ctx context.Context, name, email, passwordHash string) (model.User, error)
	GetByEmail(ctx context.Context, email string) (model.User, error)
	GetByID(ctx context.Context, id uuid.UUID) (model.User, error)
	UpdatePassword(ctx context.Context, id uuid.UUID, passwordHash string) error
}

type TokenRepo interface {
	Create(ctx context.Context, userID uuid.UUID, expiresAt time.Time) (model.PasswordResetToken, error)
	GetByToken(ctx context.Context, token uuid.UUID) (model.PasswordResetToken, error)
	MarkUsed(ctx context.Context, token uuid.UUID) error
}

type Mailer interface {
	SendPasswordReset(toEmail, resetToken string) error
}
