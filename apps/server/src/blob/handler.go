package blob

import (
	"context"

	"connectrpc.com/connect"
	blobv1 "github.com/aidan-neel/shulker/apps/proto/gen/go/blob"
	"github.com/aidan-neel/shulker/apps/proto/gen/go/blob/blobconnect"
	"github.com/aidan-neel/shulker/apps/server/pkg/middleware"
)

type Handler struct {
	service *Service
}

func New(service *Service) *Handler {
	return &Handler{service: service}
}

var _ blobconnect.BlobServiceHandler = (*Handler)(nil)

func (h *Handler) PutBlob(
	ctx context.Context,
	req *connect.Request[blobv1.PutBlobRequest],
) (*connect.Response[blobv1.PutBlobResponse], error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	blob, err := h.service.Put(ctx, req.Msg.Data, userID, req.Msg.Path, req.Msg.MimeType)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&blobv1.PutBlobResponse{
		Blob: blob,
	}), nil
}
