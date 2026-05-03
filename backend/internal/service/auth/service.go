package authsvc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"bimeet/internal/model"
)

const passwordResetTokenTTL = time.Hour

type Service struct {
	userRepo  UserRepo
	tokenRepo TokenRepo
	mailer    Mailer
	jwtSecret string
	jwtExpH   int
}

func New(userRepo UserRepo, tokenRepo TokenRepo, mailer Mailer, jwtSecret string, jwtExpHours int) *Service {
	return &Service{
		userRepo:  userRepo,
		tokenRepo: tokenRepo,
		mailer:    mailer,
		jwtSecret: jwtSecret,
		jwtExpH:   jwtExpHours,
	}
}

func (s *Service) Register(ctx context.Context, req model.RegisterRequest) (model.AuthResponse, error) {
	if req.Name == "" || req.Email == "" || req.Password == "" {
		return model.AuthResponse{}, fmt.Errorf("name, email and password are required")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.AuthResponse{}, fmt.Errorf("hash password: %w", err)
	}

	user, err := s.userRepo.Create(ctx, req.Name, req.Email, string(hash))
	if err != nil {
		return model.AuthResponse{}, fmt.Errorf("create user: %w", err)
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return model.AuthResponse{}, err
	}

	return model.AuthResponse{Token: token, User: user}, nil
}

func (s *Service) Login(ctx context.Context, req model.LoginRequest) (model.AuthResponse, error) {
	if req.Email == "" || req.Password == "" {
		return model.AuthResponse{}, fmt.Errorf("email and password are required")
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return model.AuthResponse{}, fmt.Errorf("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return model.AuthResponse{}, fmt.Errorf("invalid credentials")
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return model.AuthResponse{}, err
	}

	return model.AuthResponse{Token: token, User: user}, nil
}

// ForgotPassword always returns nil to avoid leaking which emails are registered.
// If the email exists, a reset token is generated and emailed.
func (s *Service) ForgotPassword(ctx context.Context, req model.ForgotPasswordRequest) error {
	if req.Email == "" {
		return fmt.Errorf("email is required")
	}

	user, err := s.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil
	}

	token, err := s.tokenRepo.Create(ctx, user.ID, time.Now().Add(passwordResetTokenTTL))
	if err != nil {
		log.Printf("forgot password: create token: %v", err)
		return nil
	}

	go func(email, tok string) {
		if err := s.mailer.SendPasswordReset(email, tok); err != nil {
			log.Printf("forgot password: send email to %s: %v", email, err)
		}
	}(user.Email, token.Token.String())

	return nil
}

func (s *Service) ResetPassword(ctx context.Context, req model.ResetPasswordRequest) error {
	if req.Token == "" || req.Password == "" {
		return fmt.Errorf("token and password are required")
	}

	tokenUUID, err := uuid.Parse(req.Token)
	if err != nil {
		return fmt.Errorf("invalid token")
	}

	t, err := s.tokenRepo.GetByToken(ctx, tokenUUID)
	if err != nil {
		return fmt.Errorf("invalid token")
	}
	if t.Used {
		return fmt.Errorf("invalid token")
	}
	if time.Now().After(t.ExpiresAt) {
		return fmt.Errorf("invalid token")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	if err := s.userRepo.UpdatePassword(ctx, t.UserID, string(hash)); err != nil {
		return fmt.Errorf("update password: %w", err)
	}

	if err := s.tokenRepo.MarkUsed(ctx, tokenUUID); err != nil {
		log.Printf("reset password: mark token used: %v", err)
	}

	return nil
}

func (s *Service) generateToken(userID uuid.UUID) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID.String(),
		"exp": time.Now().Add(time.Duration(s.jwtExpH) * time.Hour).Unix(),
		"iat": time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}
	return signed, nil
}
