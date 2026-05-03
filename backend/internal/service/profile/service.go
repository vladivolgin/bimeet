package profilesvc

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"

	"github.com/google/uuid"

	"bimeet/internal/model"
)

var allowedMIME = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
}

type Service struct {
	userRepo UserRepo
	storage  AvatarStorage
}

func New(userRepo UserRepo, storage AvatarStorage) *Service {
	return &Service{userRepo: userRepo, storage: storage}
}

func (s *Service) GetMe(ctx context.Context, userID uuid.UUID) (model.User, error) {
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return model.User{}, fmt.Errorf("get me: %w", err)
	}
	return user, nil
}

func (s *Service) UpdateMe(ctx context.Context, userID uuid.UUID, req model.UpdateProfileRequest) (model.User, error) {
	user, err := s.userRepo.UpdateProfile(ctx, userID, req)
	if err != nil {
		return model.User{}, fmt.Errorf("update me: %w", err)
	}
	return user, nil
}

func (s *Service) GetStats(ctx context.Context, userID uuid.UUID) (model.ProfileStats, error) {
	stats, err := s.userRepo.GetStats(ctx, userID)
	if err != nil {
		return model.ProfileStats{}, fmt.Errorf("get stats: %w", err)
	}
	return stats, nil
}

func (s *Service) DeleteAvatar(ctx context.Context, userID uuid.UUID) (model.User, error) {
	user, err := s.userRepo.ClearAvatar(ctx, userID)
	if err != nil {
		return model.User{}, fmt.Errorf("delete avatar: %w", err)
	}
	return user, nil
}

func (s *Service) UploadAvatar(ctx context.Context, userID uuid.UUID, file multipart.File, header *multipart.FileHeader) (model.User, error) {
	// Detect content type from file bytes (reliable — ignores browser-supplied header)
	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil {
		return model.User{}, fmt.Errorf("read file: %w", err)
	}
	contentType := http.DetectContentType(buf[:n])

	// Trim any parameters (e.g. "image/jpeg; charset=utf-8" → "image/jpeg")
	if idx := len(contentType); idx > 0 {
		for i, c := range contentType {
			if c == ';' {
				contentType = contentType[:i]
				break
			}
		}
	}

	ext, ok := allowedMIME[contentType]
	if !ok {
		return model.User{}, fmt.Errorf("unsupported image type: %s", contentType)
	}

	// Seek back so S3 gets the full file
	if _, err := file.Seek(0, 0); err != nil {
		return model.User{}, fmt.Errorf("seek file: %w", err)
	}

	key := "avatars/" + userID.String() + ext

	url, err := s.storage.Upload(ctx, key, contentType, file)
	if err != nil {
		return model.User{}, fmt.Errorf("upload avatar: %w", err)
	}

	user, err := s.userRepo.UpdateAvatar(ctx, userID, url)
	if err != nil {
		return model.User{}, fmt.Errorf("save avatar url: %w", err)
	}
	return user, nil
}
