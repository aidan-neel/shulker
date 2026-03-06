package auth

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

var JWT_SECRET = []byte(os.Getenv("JWT_SECRET"))

func generateToken(userID string, duration time.Duration, tokenType string) (string, error) {
	claims := jwt.MapClaims{
		"sub": userID,
		"typ": tokenType,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWT_SECRET)
}

func (s *Service) Register(ctx context.Context, email string, password string) (*Token, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	user, err := s.repo.CreateUser(ctx, email, string(hashed))
	if err != nil {
		return nil, fmt.Errorf("failed to register user: %w", err)
	}

	accessToken, err := generateToken(user.ID, 60*time.Minute, "access")
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := generateToken(user.ID, 24*7*time.Hour, "refresh")
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Service) Login(ctx context.Context, email string, password string) (*Token, error) {
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid email or password")
	}

	accessToken, err := generateToken(user.ID, 60*time.Minute, "access")
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	refreshToken, err := generateToken(user.ID, 24*7*time.Hour, "refresh")
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &Token{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
