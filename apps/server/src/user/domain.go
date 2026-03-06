package user

import (
	"context"
)

type User struct {
	ID           string
	Email        string
	CreatedAt    string
	PasswordHash string
}

type Repository interface {
	GetUser(ctx context.Context, id string) (*User, error)
	GetAllUsers(ctx context.Context) ([]*User, error)
	CreateUser(ctx context.Context, email, passwordHash string) (*User, error)
	GetUserByEmail(ctx context.Context, email string) (*User, error)
}
