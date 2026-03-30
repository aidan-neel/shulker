package routes

import (
	"context"

	"connectrpc.com/connect"
	userv1 "github.com/aidan-neel/shulker/apps/proto/gen/go/user"
	"github.com/aidan-neel/shulker/apps/proto/gen/go/user/userconnect"
	db "github.com/aidan-neel/shulker/apps/server/postgres/gen"
	"github.com/google/uuid"
)

type UsersHandler struct {
	db *db.Queries
}

func NewUsersHandler(queries *db.Queries) *UsersHandler {
	return &UsersHandler{db: queries}
}

var _ userconnect.UserServiceHandler = (*UsersHandler)(nil)

func (h *UsersHandler) GetUser(
	ctx context.Context,
	req *connect.Request[userv1.GetUserRequest],
) (*connect.Response[userv1.GetUserResponse], error) {
	uid, err := uuid.Parse(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	user, err := h.db.GetUser(ctx, uid)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	return connect.NewResponse(&userv1.GetUserResponse{
		User: &userv1.User{
			Id:        user.ID.String(),
			Email:     user.Email,
			CreatedAt: user.CreatedAt.String(),
		},
	}), nil
}

func (h *UsersHandler) GetAllUsers(
	ctx context.Context,
	req *connect.Request[userv1.GetAllUsersRequest],
) (*connect.Response[userv1.GetAllUsersResponse], error) {
	rows, err := h.db.GetAllUsers(ctx)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	users := make([]*userv1.User, len(rows))
	for i, row := range rows {
		users[i] = &userv1.User{
			Id:        row.ID.String(),
			Email:     row.Email,
			CreatedAt: row.CreatedAt.String(),
		}
	}

	return connect.NewResponse(&userv1.GetAllUsersResponse{Users: users}), nil
}
