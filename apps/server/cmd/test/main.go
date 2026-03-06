package main

import (
	"context"
	"fmt"

	"github.com/aidan-neel/shulker/apps/server/pkg/encryption"
	filesystem "github.com/aidan-neel/shulker/apps/server/pkg/filesystem"
)

func main() {
	ctx := context.Background()
	fs := filesystem.NewLocalFileSystem("./blobs")
	aes := encryption.NewAESEncryption()

	plaintext := []byte("exampleplaintext")
	encryptedData, err := aes.Encrypt(ctx, plaintext)
	if err != nil {
		panic(err)
	}
	fmt.Println("Encrypted data: ", encryptedData, "\n\n")

	decryptedData, err := aes.Decrypt(ctx, encryptedData)
	if err != nil {
		panic(err)
	}
	fmt.Println("Decrypted data: ", decryptedData, "\n\n")
	hash, _ := fs.WriteFile(ctx, decryptedData, encryptedData)
	fmt.Println(fs.GetFile(ctx, hash))
}
