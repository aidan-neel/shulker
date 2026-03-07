package auth

import (
	"context"

	"connectrpc.com/connect"
	authv1 "github.com/aidan-neel/shulker/apps/proto/gen/go/auth"
	"github.com/aidan-neel/shulker/apps/proto/gen/go/auth/authconnect"
	commonv1 "github.com/aidan-neel/shulker/apps/proto/gen/go/common"
	"github.com/aidan-neel/shulker/apps/server/pkg/middleware"
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
) (*connect.Response[authv1.RegisterResult], error) {

	u, err := h.service.Register(ctx, req.Msg.Email, req.Msg.Password)
	if err != nil {
		return connect.NewResponse(&authv1.RegisterResult{
			Result: &authv1.RegisterResult_Error{
				Error: &commonv1.ErrorResponse{
					Code:    "INVALID_ARGUMENT",
					Message: err.Error(),
				},
			},
		}), nil
	}

	middleware.SetAuthCookies(ctx, u.AccessToken, u.RefreshToken)

	return connect.NewResponse(&authv1.RegisterResult{
		Result: &authv1.RegisterResult_Success{
			Success: &authv1.RegisterResponse{},
		},
	}), nil
}

func (h *Handler) Login(
	ctx context.Context,
	req *connect.Request[authv1.LoginRequest],
) (*connect.Response[authv1.LoginResult], error) {

	u, err := h.service.Login(ctx, req.Msg.Email, req.Msg.Password)
	if err != nil {
		return connect.NewResponse(&authv1.LoginResult{
			Result: &authv1.LoginResult_Error{
				Error: &commonv1.ErrorResponse{
					Code:    "INVALID_ARGUMENT",
					Message: err.Error(),
				},
			},
		}), nil
	}

	middleware.SetAuthCookies(ctx, u.AccessToken, u.RefreshToken)

	return connect.NewResponse(&authv1.LoginResult{
		Result: &authv1.LoginResult_Success{
			Success: &authv1.LoginResponse{},
		},
	}), nil
}
