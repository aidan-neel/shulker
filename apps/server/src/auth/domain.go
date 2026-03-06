package auth

import (
	"context"

	"github.com/aidan-neel/shulker/apps/server/src/user"
)

type Token struct {
	AccessToken  string
	RefreshToken string
}

type Repository interface {
	CreateUser(ctx context.Context, email, passwordHash string) (*user.User, error)
	GetUserByEmail(ctx context.Context, email string) (*user.User, error)
}
