package profilesvc

import (
	"context"
	"io"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

type UserRepo interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.User, error)
	UpdateProfile(ctx context.Context, id uuid.UUID, req model.UpdateProfileRequest) (model.User, error)
	GetStats(ctx context.Context, id uuid.UUID) (model.ProfileStats, error)
	UpdateAvatar(ctx context.Context, id uuid.UUID, url string) (model.User, error)
	ClearAvatar(ctx context.Context, id uuid.UUID) (model.User, error)
}

type AvatarStorage interface {
	Upload(ctx context.Context, key, contentType string, body io.Reader) (string, error)
}
