package passwordresetrepo

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"bimeet/internal/model"
)

type Repository struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

func (r *Repository) Create(ctx context.Context, userID uuid.UUID, expiresAt time.Time) (model.PasswordResetToken, error) {
	var t model.PasswordResetToken
	err := r.pool.QueryRow(ctx,
		`INSERT INTO password_reset_tokens (user_id, expires_at)
		 VALUES ($1, $2)
		 RETURNING token, user_id, expires_at, used, created_at`,
		userID, expiresAt,
	).Scan(&t.Token, &t.UserID, &t.ExpiresAt, &t.Used, &t.CreatedAt)
	if err != nil {
		return model.PasswordResetToken{}, fmt.Errorf("create password reset token: %w", err)
	}
	return t, nil
}

func (r *Repository) GetByToken(ctx context.Context, token uuid.UUID) (model.PasswordResetToken, error) {
	var t model.PasswordResetToken
	err := r.pool.QueryRow(ctx,
		`SELECT token, user_id, expires_at, used, created_at
		 FROM password_reset_tokens WHERE token=$1`,
		token,
	).Scan(&t.Token, &t.UserID, &t.ExpiresAt, &t.Used, &t.CreatedAt)
	if err != nil {
		return model.PasswordResetToken{}, fmt.Errorf("get password reset token: %w", err)
	}
	return t, nil
}

func (r *Repository) MarkUsed(ctx context.Context, token uuid.UUID) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE password_reset_tokens SET used=TRUE WHERE token=$1`,
		token,
	)
	if err != nil {
		return fmt.Errorf("mark password reset token used: %w", err)
	}
	return nil
}
