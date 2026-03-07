package blob

import (
	"context"

	blobpb "github.com/aidan-neel/shulker/apps/proto/gen/go/blob"
)

type Repository interface {
	UpsertBlob(ctx context.Context, hash, mimeType string, size int64) (*blobpb.Blob, error)
	GetBlob(ctx context.Context, id string) (*blobpb.Blob, error)
	GetBlobByHash(ctx context.Context, hash string) (*blobpb.Blob, error)
	DeleteBlob(ctx context.Context, id string) error
	DecrementBlob(ctx context.Context, hash string) (*blobpb.Blob, error)
}

type UserBlobRepository interface {
	CreateUserBlob(ctx context.Context, userID, blobID, path string) (*blobpb.UserBlob, error)
	GetUserBlobs(ctx context.Context, userID string) ([]*blobpb.Blob, error)
	GetUserBlobByPath(ctx context.Context, userID, path string) (*blobpb.Blob, error)
	DeleteUserBlob(ctx context.Context, userID, blobID string) error
}
