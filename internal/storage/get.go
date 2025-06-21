package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/DhruvDattani1/edgevault/internal/crypto"
	"golang.org/x/crypto/chacha20poly1305"
)

func Get(objectName string, destPath string, masterKey []byte) error {
	srcPath := filepath.Join(storageDir, objectName+".crypta")

	inFile, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open encrypted file: %w", err)
	}
	defer inFile.Close()

	header := make([]byte, 4)
	if _, err := io.ReadFull(inFile, header); err != nil {
		return fmt.Errorf("failed to read file header: %w", err)
	}

	if string(header) == "EV1\x00" {
		inFile.Close()
		return crypto.DecryptLargeFile(srcPath, destPath, masterKey)
	}

	if _, err := inFile.Seek(0, io.SeekStart); err != nil {
		return fmt.Errorf("seek failed: %w", err)
	}

	data, err := io.ReadAll(inFile)
	if err != nil {
		return fmt.Errorf("read failed: %w", err)
	}

	if len(data) < chacha20poly1305.NonceSize {
		return fmt.Errorf("invalid file format (too small)")
	}

	nonce := data[:chacha20poly1305.NonceSize]
	ciphertext := data[chacha20poly1305.NonceSize:]

	plaintext, err := crypto.Decrypt(masterKey, nonce, ciphertext)
	if err != nil {
		return fmt.Errorf("decryption failed: %w", err)
	}

	outFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("can't create output file: %w", err)
	}
	defer outFile.Close()

	if _, err := outFile.Write(plaintext); err != nil {
		return fmt.Errorf("failed to write decrypted file: %w", err)
	}

	return nil
}
