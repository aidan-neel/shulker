package blob

import (
	"context"

	blobpb "github.com/aidan-neel/shulker/apps/proto/gen/go/blob"
)

type Repository interface {
	UpsertBlob(ctx context.Context, userID, hash, filepath, mimeType string, size int64) (*blobpb.Blob, error)
	GetBlob(ctx context.Context, id string) (*blobpb.Blob, error)
	GetBlobByHash(ctx context.Context, hash string) (*blobpb.Blob, error)
	GetBlobsByUser(ctx context.Context, userID string) ([]*blobpb.Blob, error)
	DeleteBlob(ctx context.Context, id string) error
}
