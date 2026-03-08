package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/aidan-neel/shulker/apps/server/pkg/middleware"
	"github.com/aidan-neel/shulker/apps/server/pkg/utils"
	"github.com/aidan-neel/shulker/apps/server/src/user"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo user.Repository
}

func NewService(repo user.Repository) *Service {
	return &Service{repo: repo}
}

var JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

func (s *Service) Register(ctx context.Context, email string, password string) (*Token, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.repo.CreateUser(ctx, email, string(hashed))
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	accessToken, err := utils.GenerateToken(user.ID, 60*time.Minute, "access")
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateToken(user.ID, 24*7*time.Hour, "refresh")
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) CurrentUser(ctx context.Context) (*user.User, error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, fmt.Errorf("unauthorized")
	}

	user, err := s.repo.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	return user, nil
}

func (s *Service) Login(ctx context.Context, email string, password string) (*Token, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	accessToken, err := utils.GenerateToken(user.ID, 60*time.Minute, "access")
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := utils.GenerateToken(user.ID, 24*7*time.Hour, "refresh")
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
