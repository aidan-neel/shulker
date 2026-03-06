package user

import (
	"context"

	"connectrpc.com/connect"
	userv1 "github.com/aidan-neel/shulker/apps/proto/gen/go/user"
	"github.com/aidan-neel/shulker/apps/proto/gen/go/user/userconnect"
)

type Handler struct {
	service *Service
}

func New(service *Service) *Handler {
	return &Handler{service: service}
}

var _ userconnect.UserServiceHandler = (*Handler)(nil)

func (h *Handler) GetUser(
	ctx context.Context,
	req *connect.Request[userv1.GetUserRequest],
) (*connect.Response[userv1.GetUserResponse], error) {
	u, err := h.service.GetUser(ctx, req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	return connect.NewResponse(&userv1.GetUserResponse{
		User: &userv1.User{
			Id:        u.ID,
			Email:     u.Email,
			CreatedAt: u.CreatedAt,
		},
	}), nil
}

func (h *Handler) GetAllUsers(
	ctx context.Context,
	req *connect.Request[userv1.GetAllUsersRequest],
) (*connect.Response[userv1.GetAllUsersResponse], error) {
	u, err := h.service.GetAllUsers(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	users := make([]*userv1.User, len(u))
	for i, user := range u {
		users[i] = &userv1.User{
			Id:        user.ID,
			Email:     user.Email,
			CreatedAt: user.CreatedAt,
		}
	}

	return connect.NewResponse(&userv1.GetAllUsersResponse{
		Users: users,
	}), nil
}
