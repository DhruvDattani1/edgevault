package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	
	"github.com/DhruvDattani1/edgevault/internal/crypto"

)

const (
	storageDir = "crypta"
	bufferSize = 128 * 1024
)

func Put(sourceFile string, masterKey []byte) error {

	inFile, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("can't open source file: %w", err)
	}
	defer inFile.Close()

	plaintext, err := io.ReadAll(inFile)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	nonce, ciphertext, err := crypto.Encrypt(masterKey, plaintext)
	if err != nil {
		return fmt.Errorf("encryption failed: %w", err)
	}

	destFilename := filepath.Base(sourceFile) + ".crypta"
	partialPath := filepath.Join(storageDir, destFilename+".partial")
	finalPath := filepath.Join(storageDir, destFilename)

	err = os.MkdirAll(storageDir, 0700)
	if err != nil {
		return fmt.Errorf("no dir created: %w", err)
	}

	outFile, err := os.OpenFile(partialPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("can't make partial file: %w", err)
	}

	defer func() {
		for i := range plaintext {
			plaintext[i] = 0
		}
		outFile.Close()
	}()

	if _, err := outFile.Write(nonce); err != nil {
		return fmt.Errorf("failed to write nonce: %w", err)
	}

	if _, err := outFile.Write(ciphertext); err != nil {
		return fmt.Errorf("failed to write ciphertext: %w", err)
	}

	if err := outFile.Sync(); err != nil {
		return fmt.Errorf("partial didn't sync: %w", err)
	}

	if err := os.Rename(partialPath, finalPath); err != nil {
		return fmt.Errorf("partial not renamed: %w", err)
	}

	return nil
}