package encryption

import (
	"context"
)

type Encryption interface {
	Encrypt(ctx context.Context, plaintext []byte) ([]byte, error)
	Decrypt(ctx context.Context, data []byte) ([]byte, error)
}
