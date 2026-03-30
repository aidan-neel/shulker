package routes

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	blobv1 "github.com/aidan-neel/shulker/apps/proto/gen/go/blob"
	"github.com/aidan-neel/shulker/apps/proto/gen/go/blob/blobconnect"
	"github.com/aidan-neel/shulker/apps/server/middleware"
	db "github.com/aidan-neel/shulker/apps/server/postgres/gen"
	"github.com/aidan-neel/shulker/apps/server/storage"
	"github.com/google/uuid"
)

type BlobsHandler struct {
	db         *db.Queries
	filesystem storage.Filesystem
	encryption storage.Encryption
}

func NewBlobsHandler(queries *db.Queries, fs storage.Filesystem, enc storage.Encryption) *BlobsHandler {
	return &BlobsHandler{db: queries, filesystem: fs, encryption: enc}
}

var _ blobconnect.BlobServiceHandler = (*BlobsHandler)(nil)

func (h *BlobsHandler) PutBlob(
	ctx context.Context,
	req *connect.Request[blobv1.PutBlobRequest],
) (*connect.Response[blobv1.PutBlobResponse], error) {
	userID, ok := middleware.GetUserID(ctx)
	if !ok {
		return nil, connect.NewError(connect.CodeUnauthenticated, nil)
	}

	encrypted, err := h.encryption.Encrypt(ctx, req.Msg.Data, req.Msg.MimeType)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	hash, err := h.filesystem.WriteFile(ctx, req.Msg.Data, encrypted)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	blob, err := h.db.UpsertBlob(ctx, db.UpsertBlobParams{
		Hash:     hash,
		MimeType: req.Msg.MimeType,
		Size:     int64(len(req.Msg.Data)),
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	path := h.uniquePath(ctx, userID, req.Msg.Path)

	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	_, err = h.db.CreateUserBlob(ctx, db.CreateUserBlobParams{
		UserID: uid,
		BlobID: blob.ID,
		Path:   path,
	})
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&blobv1.PutBlobResponse{
		Blob: rowToBlob(blob),
	}), nil
}

func (h *BlobsHandler) uniquePath(ctx context.Context, userID, path string) string {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return path
	}
	candidate := path
	for i := 1; ; i++ {
		_, err := h.db.GetUserBlobByPath(ctx, db.GetUserBlobByPathParams{UserID: uid, Path: candidate})
		if err != nil {
			return candidate
		}
		candidate = fmt.Sprintf("%s (%d)", path, i)
	}
}

func rowToBlob(row db.Blob) *blobv1.Blob {
	return &blobv1.Blob{
		Id:        row.ID.String(),
		Hash:      row.Hash,
		MimeType:  row.MimeType,
		Size:      row.Size,
		CreatedAt: row.CreatedAt.String(),
	}
}
