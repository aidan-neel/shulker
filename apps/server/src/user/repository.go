package user

import (
	"context"

	db "github.com/aidan-neel/shulker/apps/server/postgres/gen"
	"github.com/google/uuid"
)

type UserRepo struct {
	queries *db.Queries
}

func NewUserRepo(queries *db.Queries) *UserRepo {
	return &UserRepo{queries: queries}
}

func (r *UserRepo) CreateUser(ctx context.Context, email, passwordHash string) (*User, error) {
	row, err := r.queries.CreateUser(ctx, db.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
	})
	if err != nil {
		return nil, err
	}
	return &User{
		ID:        row.ID.String(),
		Email:     row.Email,
		CreatedAt: row.CreatedAt.String(),
	}, nil
}

func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*User, error) {
	row, err := r.queries.GetUserByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return &User{
		ID:           row.ID.String(),
		Email:        row.Email,
		CreatedAt:    row.CreatedAt.String(),
		PasswordHash: row.PasswordHash,
	}, nil
}

func (r *UserRepo) GetUser(ctx context.Context, id string) (*User, error) {
	uuid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	row, err := r.queries.GetUser(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return &User{
		ID:        row.ID.String(),
		Email:     row.Email,
		CreatedAt: row.CreatedAt.String(),
	}, nil
}

func (r *UserRepo) GetAllUsers(ctx context.Context) ([]*User, error) {
	rows, err := r.queries.GetAllUsers(ctx)
	if err != nil {
		return nil, err
	}

	users := make([]*User, len(rows))
	for i, row := range rows {
		users[i] = &User{
			ID:        row.ID.String(),
			Email:     row.Email,
			CreatedAt: row.CreatedAt.String(),
		}
	}

	return users, nil
}
