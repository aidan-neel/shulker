package routes

import (
	"context"
	"fmt"
	"time"

	"log"

	"connectrpc.com/connect"
	authv1 "github.com/aidan-neel/shulker/apps/proto/gen/go/auth"
	"github.com/aidan-neel/shulker/apps/proto/gen/go/auth/authconnect"
	commonv1 "github.com/aidan-neel/shulker/apps/proto/gen/go/common"
	userv1 "github.com/aidan-neel/shulker/apps/proto/gen/go/user"
	"github.com/aidan-neel/shulker/apps/server/jwt"
	"github.com/aidan-neel/shulker/apps/server/middleware"
	db "github.com/aidan-neel/shulker/apps/server/postgres/gen"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	db *db.Queries
}

func NewAuthHandler(queries *db.Queries) *AuthHandler {
	return &AuthHandler{db: queries}
}

var _ authconnect.AuthServiceHandler = (*AuthHandler)(nil)

func (h *AuthHandler) Register(
	ctx context.Context,
	req *connect.Request[authv1.RegisterRequest],
) (*connect.Response[authv1.RegisterResult], error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Msg.Password), bcrypt.DefaultCost)
	if err != nil {
		return errRegister("failed to hash password"), nil
	}

	user, err := h.db.CreateUser(ctx, db.CreateUserParams{
		Email:        req.Msg.Email,
		PasswordHash: string(hashed),
	})
	if err != nil {
		return errRegister(err.Error()), nil
	}

	accessToken, refreshToken, err := generateTokenPair(user.ID.String())
	if err != nil {
		return errRegister("failed to generate tokens"), nil
	}

	middleware.SetAuthCookies(ctx, accessToken, refreshToken)

	return connect.NewResponse(&authv1.RegisterResult{
		Result: &authv1.RegisterResult_Success{
			Success: &authv1.RegisterResponse{},
		},
	}), nil
}

func (h *AuthHandler) Login(
	ctx context.Context,
	req *connect.Request[authv1.LoginRequest],
) (*connect.Response[authv1.LoginResult], error) {
	user, err := h.db.GetUserByEmail(ctx, req.Msg.Email)
	if err != nil {
		return errLogin(err.Error()), nil
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Msg.Password)); err != nil {
		return errLogin(err.Error()), nil
	}

	accessToken, refreshToken, err := generateTokenPair(user.ID.String())
	if err != nil {
		return errLogin(err.Error()), nil
	}

	middleware.SetAuthCookies(ctx, accessToken, refreshToken)

	return connect.NewResponse(&authv1.LoginResult{
		Result: &authv1.LoginResult_Success{
			Success: &authv1.LoginResponse{},
		},
	}), nil
}

func (h *AuthHandler) CurrentUser(
	ctx context.Context,
	req *connect.Request[authv1.CurrentUserRequest],
) (*connect.Response[authv1.CurrentUserResult], error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return errCurrentUser("unauthorized"), nil
	}

	uid, err := uuid.Parse(userID)
	if err != nil {
		return errCurrentUser(err.Error()), nil
	}

	user, err := h.db.GetUser(ctx, uid)
	if err != nil {
		return errCurrentUser(err.Error()), nil
	}

	return connect.NewResponse(&authv1.CurrentUserResult{
		Result: &authv1.CurrentUserResult_Success{
			Success: &authv1.CurrentUserResponse{
				Me: &userv1.User{
					Id:        user.ID.String(),
					Email:     user.Email,
					CreatedAt: user.CreatedAt.String(),
				},
			},
		},
	}), nil
}

func generateTokenPair(userID string) (access, refresh string, err error) {
	access, err = jwt.GenerateToken(userID, 60*time.Minute, "access")
	if err != nil {
		return "", "", fmt.Errorf("failed to generate access token: %w", err)
	}
	refresh, err = jwt.GenerateToken(userID, 7*24*time.Hour, "refresh")
	if err != nil {
		return "", "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return access, refresh, nil
}

func errRegister(msg string) *connect.Response[authv1.RegisterResult] {
	log.Println("Error: ", msg)
	return connect.NewResponse(&authv1.RegisterResult{
		Result: &authv1.RegisterResult_Error{
			Error: &commonv1.ErrorResponse{Code: "INVALID_ARGUMENT", Message: msg},
		},
	})
}

func errLogin(msg string) *connect.Response[authv1.LoginResult] {
	log.Println("Error: ", msg)
	return connect.NewResponse(&authv1.LoginResult{
		Result: &authv1.LoginResult_Error{
			Error: &commonv1.ErrorResponse{Code: "INVALID_ARGUMENT", Message: msg},
		},
	})
}

func errCurrentUser(msg string) *connect.Response[authv1.CurrentUserResult] {
	log.Println("Error: ", msg)
	return connect.NewResponse(&authv1.CurrentUserResult{
		Result: &authv1.CurrentUserResult_Error{
			Error: &commonv1.ErrorResponse{Code: "INVALID_ARGUMENT", Message: msg},
		},
	})
}
