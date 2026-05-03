package profilehandler

import (
	"context"
	"mime/multipart"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

type ProfileService interface {
	GetMe(ctx context.Context, userID uuid.UUID) (model.User, error)
	UpdateMe(ctx context.Context, userID uuid.UUID, req model.UpdateProfileRequest) (model.User, error)
	GetStats(ctx context.Context, userID uuid.UUID) (model.ProfileStats, error)
	UploadAvatar(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader) (model.User, error)
	DeleteAvatar(ctx context.Context, userID uuid.UUID) (model.User, error)
}
