package encryption

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"os"
)

func shouldCompress(mimeType string) bool {
	switch mimeType {
	case "text/plain", "text/html", "application/json", "text/csv":
		return true
	}
	return false
}

func compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write(data); err != nil {
		return nil, err
	}
	if err := w.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decompress(data []byte) ([]byte, error) {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer r.Close()
	return io.ReadAll(r)
}

type AESEncryption struct {
	key []byte
}

func NewAESEncryption() *AESEncryption {
	keyHex := os.Getenv("AES_256_KEY")

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		panic(fmt.Sprintf("failed to decode AES key: %v", err))
	}

	if len(key) != 32 {
		panic(fmt.Sprintf("invalid key size: expected 32 bytes, got %d", len(key)))
	}

	return &AESEncryption{
		key: key,
	}
}

func (s *AESEncryption) Encrypt(ctx context.Context, plaintext []byte, mimeType string) ([]byte, error) {
	data := plaintext
	if shouldCompress(mimeType) {
		compressed, err := compress(data)
		if err != nil {
			return nil, err
		}
		data = compressed
	}

	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, fmt.Errorf("failed to aes encrypt: %w", err)
	}

	nonce := make([]byte, 12)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to nonce: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	ciphertext := aesgcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func (s *AESEncryption) Decrypt(ctx context.Context, data []byte, mimeType string) ([]byte, error) {
	block, err := aes.NewCipher(s.key)
	if err != nil {
		return nil, fmt.Errorf("failed to create cipher: %w", err)
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("failed to create gcm: %w", err)
	}

	nonceSize := aesgcm.NonceSize()
	if len(data) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to decrypt: %w", err)
	}

	if shouldCompress(mimeType) {
		decompressed, err := decompress(plaintext)
		if err != nil {
			return nil, err
		}
		plaintext = decompressed
	}

	return plaintext, nil
}
