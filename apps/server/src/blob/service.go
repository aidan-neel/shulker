package blob

import (
	"context"

	"github.com/aidan-neel/shulker/apps/server/pkg/encryption"
	filesystem "github.com/aidan-neel/shulker/apps/server/pkg/filesystem"
)

type Service struct {
	repo       Repository
	filesystem filesystem.Filesystem
	encryption encryption.Encryption
}

func NewService(repo Repository) *Service {
	fs := filesystem.NewLocalFileSystem("./blobs")
	enc := encryption.NewAESEncryption()

	return &Service{
		filesystem: fs,
		encryption: enc,
		repo:       repo,
	}
}

func (s *Service) Put(ctx context.Context, data []byte, userID string, filePath string, mimeType string) (*Blob, error) {
	user_id := "019cc122-be83-74bd-a725-6aae5019bbd1"

	encrypted, err := s.encryption.Encrypt(ctx, data, mimeType)
	if err != nil {
		return nil, err
	}

	hash, err := s.filesystem.WriteFile(ctx, data, encrypted)
	if err != nil {
		return nil, err
	}

	blob, err := s.repo.UpsertBlob(ctx, user_id, hash, filePath, mimeType, int64(len(data)))
	if err != nil {
		return nil, err
	}

	return blob, err
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

	decrypted, err := s.encryption.Decrypt(ctx, data, blob.MimeType)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}
