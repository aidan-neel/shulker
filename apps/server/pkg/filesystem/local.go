package filesystem

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/zeebo/blake3"
)

type LocalFileSystem struct {
	root string
}

func hashContent(r io.Reader) (string, error) {
	h := blake3.New()
	if _, err := io.Copy(h, r); err != nil {
		return "", err
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func NewLocalFileSystem(root string) *LocalFileSystem {
	newpath := filepath.Join(root)
	os.MkdirAll(newpath, os.ModePerm)
	return &LocalFileSystem{root: root}
}

func (s *LocalFileSystem) WriteFile(ctx context.Context, decryptedData []byte, encryptedData []byte) (string, error) {
	r := bytes.NewReader(decryptedData)
	hash, err := hashContent(r)
	if err != nil {
		return "", fmt.Errorf("failed to blake3 hash: %w", err)
	}
	folderPath := filepath.Join(s.root, hash[:2])
	if err := os.MkdirAll(folderPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create data folder: %w", err)
	}

	path := filepath.Join(folderPath, hash)
	err = os.WriteFile(path, encryptedData, 0644)
	return hash, err
}

func (s *LocalFileSystem) GetFile(ctx context.Context, hash string) ([]byte, error) {
	path := filepath.Join(s.root, hash[:2], hash)
	return os.ReadFile(path)
}
