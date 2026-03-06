package filesystem

import (
	"context"
)

type Filesystem interface {
	WriteFile(ctx context.Context, decryptedData []byte, encryptedData []byte) (string, error)
	GetFile(ctx context.Context, hash string) ([]byte, error)
}
