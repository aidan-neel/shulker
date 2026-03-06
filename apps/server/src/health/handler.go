package healthhandler

import (
	"context"

	"connectrpc.com/connect"
	healthv1 "github.com/aidan-neel/shulker/apps/proto/gen/go/health"
	"github.com/aidan-neel/shulker/apps/proto/gen/go/health/healthconnect"
)

type Handler struct{}

func New() *Handler {
	return &Handler{}
}

var _ healthconnect.HealthServiceHandler = (*Handler)(nil)

// Check godoc
// @Summary      Health check
// @Description  Returns server status
// @Tags         health
// @Produce      json
// @Success      200  {object}  healthv1.HealthCheckResponse
// @Router       /health [get]
func (h *Handler) Check(
	ctx context.Context,
	req *connect.Request[healthv1.HealthCheckRequest],
) (*connect.Response[healthv1.HealthCheckResponse], error) {
	return connect.NewResponse(&healthv1.HealthCheckResponse{
		Status: "ok",
	}), nil
}
