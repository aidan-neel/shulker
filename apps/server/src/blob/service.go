package blob

import (
	"context"
	"fmt"

	blobpb "github.com/aidan-neel/shulker/apps/proto/gen/go/blob"
	"github.com/aidan-neel/shulker/apps/server/pkg/encryption"
	filesystem "github.com/aidan-neel/shulker/apps/server/pkg/filesystem"
)

type Service struct {
	repo         Repository
	userBlobRepo UserBlobRepository
	filesystem   filesystem.Filesystem
	encryption   encryption.Encryption
}

func NewService(repo Repository, userBlobRepo UserBlobRepository) *Service {
	fs := filesystem.NewLocalFileSystem("./blobs")
	enc := encryption.NewAESEncryption()

	return &Service{
		filesystem:   fs,
		encryption:   enc,
		repo:         repo,
		userBlobRepo: userBlobRepo,
	}
}

func (s *Service) uniquePath(ctx context.Context, userID, path string) string {
	candidate := path
	for i := 1; ; i++ {
		_, err := s.userBlobRepo.GetUserBlobByPath(ctx, userID, candidate)
		if err != nil {
			return candidate
		}
		candidate = fmt.Sprintf("%s (%d)", path, i)
	}
}

func (s *Service) Put(ctx context.Context, data []byte, userID string, path string, mimeType string) (*blobpb.Blob, error) {
	encrypted, err := s.encryption.Encrypt(ctx, data, mimeType)
	if err != nil {
		return nil, err
	}

	hash, err := s.filesystem.WriteFile(ctx, data, encrypted)
	if err != nil {
		return nil, err
	}

	blob, err := s.repo.UpsertBlob(ctx, hash, mimeType, int64(len(data)))
	if err != nil {
		return nil, err
	}

	path = s.uniquePath(ctx, userID, path)

	_, err = s.userBlobRepo.CreateUserBlob(ctx, userID, blob.Id, path)
	if err != nil {
		return nil, err
	}

	return blob, nil
}

func (s *Service) Get(ctx context.Context, hash string) ([]byte, error) {
	blob, err := s.repo.GetBlobByHash(ctx, hash)
	if err != nil {
		return nil, err
	}

	data, err := s.filesystem.GetFile(ctx, blob.Hash)
	if err != nil {
		return nil, err
	}

	return s.encryption.Decrypt(ctx, data, blob.MimeType)
}
