package auth

import (
	"context"

	"connectrpc.com/connect"
	authv1 "github.com/aidan-neel/shulker/apps/proto/gen/go/auth"
	"github.com/aidan-neel/shulker/apps/proto/gen/go/auth/authconnect"
)

type Handler struct {
	service *Service
}

func New(service *Service) *Handler {
	return &Handler{service: service}
}

var _ authconnect.AuthServiceHandler = (*Handler)(nil)

func (h *Handler) Register(
	ctx context.Context,
	req *connect.Request[authv1.RegisterRequest],
) (*connect.Response[authv1.RegisterResponse], error) {
	u, err := h.service.Register(ctx, req.Msg.Email, req.Msg.Password)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return connect.NewResponse(&authv1.RegisterResponse{
		Token: &authv1.Token{
			AccessToken:  u.AccessToken,
			RefreshToken: u.RefreshToken,
		},
	}), nil
}

func (h *Handler) Login(
	ctx context.Context,
	req *connect.Request[authv1.LoginRequest],
) (*connect.Response[authv1.LoginResponse], error) {
	u, err := h.service.Login(ctx, req.Msg.Email, req.Msg.Password)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return connect.NewResponse(&authv1.LoginResponse{
		Token: &authv1.Token{
			AccessToken:  u.AccessToken,
			RefreshToken: u.RefreshToken,
		},
	}), nil
}
