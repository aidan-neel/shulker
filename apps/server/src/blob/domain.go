package blob

import "context"

type Blob struct {
	ID        string
	UserID    string
	Hash      string
	Filepath  string
	MimeType  string
	Size      int64
	CreatedAt string
}

type Repository interface {
	UpsertBlob(ctx context.Context, userID, hash, filepath, mimeType string, size int64) (*Blob, error)
	GetBlob(ctx context.Context, id string) (*Blob, error)
	GetBlobByHash(ctx context.Context, hash string) (*Blob, error)
	GetBlobsByUser(ctx context.Context, userID string) ([]*Blob, error)
	DeleteBlob(ctx context.Context, id string) error
}
