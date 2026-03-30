package routes

import (
	"context"

	"connectrpc.com/connect"
	healthv1 "github.com/aidan-neel/shulker/apps/proto/gen/go/health"
	"github.com/aidan-neel/shulker/apps/proto/gen/go/health/healthconnect"
)

type HealthHandler struct{}

func NewHealthHandler() *HealthHandler {
	return &HealthHandler{}
}

var _ healthconnect.HealthServiceHandler = (*HealthHandler)(nil)

func (h *HealthHandler) Check(
	ctx context.Context,
	req *connect.Request[healthv1.HealthCheckRequest],
) (*connect.Response[healthv1.HealthCheckResponse], error) {
	return connect.NewResponse(&healthv1.HealthCheckResponse{Status: "ok"}), nil
}
