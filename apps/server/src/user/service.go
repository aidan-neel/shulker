package user

import (
	"context"
)

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
	return s.repo.GetUser(ctx, id)
}

func (s *Service) GetAllUsers(ctx context.Context) ([]*User, error) {
	return s.repo.GetAllUsers(ctx)
}
