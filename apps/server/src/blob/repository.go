package blob

import (
	"context"

	db "github.com/aidan-neel/shulker/apps/server/postgres/gen"
	"github.com/google/uuid"
)

type BlobRepo struct {
	queries *db.Queries
}

func NewBlobRepo(queries *db.Queries) *BlobRepo {
	return &BlobRepo{queries: queries}
}

func (r *BlobRepo) UpsertBlob(ctx context.Context, userID, hash, filepath, mimeType string, size int64) (*Blob, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	row, err := r.queries.UpsertBlob(ctx, db.UpsertBlobParams{
		UserID:   uid,
		Hash:     hash,
		Filepath: filepath,
		MimeType: mimeType,
		Size:     size,
	})
	if err != nil {
		return nil, err
	}
	return rowToBlob(row), nil
}

func (r *BlobRepo) GetBlob(ctx context.Context, id string) (*Blob, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	row, err := r.queries.GetBlob(ctx, uid)
	if err != nil {
		return nil, err
	}
	return rowToBlob(row), nil
}

func (r *BlobRepo) GetBlobByHash(ctx context.Context, hash string) (*Blob, error) {
	row, err := r.queries.GetBlobByHash(ctx, hash)
	if err != nil {
		return nil, err
	}
	return rowToBlob(row), nil
}

func (r *BlobRepo) GetBlobsByUser(ctx context.Context, userID string) ([]*Blob, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	rows, err := r.queries.GetBlobsByUser(ctx, uid)
	if err != nil {
		return nil, err
	}
	blobs := make([]*Blob, len(rows))
	for i, row := range rows {
		blobs[i] = rowToBlob(row)
	}
	return blobs, nil
}

func (r *BlobRepo) DeleteBlob(ctx context.Context, id string) error {
	uid, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return r.queries.DeleteBlob(ctx, uid)
}

func (r *BlobRepo) DecrementBlob(ctx context.Context, hash string) (*Blob, error) {
	row, err := r.queries.DecrementBlob(ctx, hash)
	if err != nil {
		return nil, err
	}
	return rowToBlob(row), nil
}

func rowToBlob(row db.Blob) *Blob {
	return &Blob{
		ID:        row.ID.String(),
		UserID:    row.UserID.String(),
		Hash:      row.Hash,
		Filepath:  row.Filepath,
		MimeType:  row.MimeType,
		Size:      row.Size,
		CreatedAt: row.CreatedAt.String(),
	}
}
