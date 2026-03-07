package encryption

import (
	"context"
)

type Encryption interface {
	Encrypt(ctx context.Context, plaintext []byte, mimeType string) ([]byte, error)
	Decrypt(ctx context.Context, data []byte, mimeType string) ([]byte, error)
}
