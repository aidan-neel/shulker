package blob

import (
	"context"

	blobpb "github.com/aidan-neel/shulker/apps/proto/gen/go/blob"
	db "github.com/aidan-neel/shulker/apps/server/postgres/gen"
	"github.com/google/uuid"
)

type UserBlobRepo struct {
	queries *db.Queries
}

func NewUserBlobRepo(queries *db.Queries) *UserBlobRepo {
	return &UserBlobRepo{queries: queries}
}

func (r *UserBlobRepo) CreateUserBlob(ctx context.Context, userID, blobID, path string) (*blobpb.UserBlob, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	bid, err := uuid.Parse(blobID)
	if err != nil {
		return nil, err
	}
	row, err := r.queries.CreateUserBlob(ctx, db.CreateUserBlobParams{
		UserID: uid,
		BlobID: bid,
		Path:   path,
	})
	if err != nil {
		return nil, err
	}
	return rowToUserBlob(row), nil
}

func (r *UserBlobRepo) GetUserBlobs(ctx context.Context, userID string) ([]*blobpb.Blob, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	rows, err := r.queries.GetUserBlobs(ctx, uid)
	if err != nil {
		return nil, err
	}
	blobs := make([]*blobpb.Blob, len(rows))
	for i, row := range rows {
		blobs[i] = rowToBlob(db.Blob(row))
	}
	return blobs, nil
}

func (r *UserBlobRepo) GetUserBlobByPath(ctx context.Context, userID, path string) (*blobpb.Blob, error) {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}
	row, err := r.queries.GetUserBlobByPath(ctx, db.GetUserBlobByPathParams{
		UserID: uid,
		Path:   path,
	})
	if err != nil {
		return nil, err
	}
	return rowToBlob(db.Blob(row)), nil
}

func (r *UserBlobRepo) DeleteUserBlob(ctx context.Context, userID, blobID string) error {
	uid, err := uuid.Parse(userID)
	if err != nil {
		return err
	}
	bid, err := uuid.Parse(blobID)
	if err != nil {
		return err
	}
	return r.queries.DeleteUserBlob(ctx, db.DeleteUserBlobParams{
		UserID: uid,
		BlobID: bid,
	})
}

func rowToUserBlob(row db.UserBlob) *blobpb.UserBlob {
	return &blobpb.UserBlob{
		Id:        row.ID.String(),
		UserId:    row.UserID.String(),
		BlobId:    row.BlobID.String(),
		Path:      row.Path,
		CreatedAt: row.CreatedAt.String(),
	}
}
