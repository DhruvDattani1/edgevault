package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/DhruvDattani1/edgevault/internal/crypto"
)

const (
	storageDir          = "crypta"
	LargeFileThreshold  = 10 * 1024 * 1024 // 10 MB threshold
)

func Put(sourceFile string, masterKey []byte) error {
	inFile, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("can't open source file: %w", err)
	}
	defer inFile.Close()

	info, err := inFile.Stat()
	if err != nil {
		return fmt.Errorf("can't stat source file: %w", err)
	}

	destFilename := filepath.Base(sourceFile) + ".crypta"
	partialPath := filepath.Join(storageDir, destFilename+".partial")
	finalPath := filepath.Join(storageDir, destFilename)

	err = os.MkdirAll(storageDir, 0700)
	if err != nil {
		return fmt.Errorf("no dir created: %w", err)
	}

	if info.Size() > LargeFileThreshold {

		fmt.Println("Encrypting in chunks (large file)...")


		err = crypto.EncryptLargeFile(sourceFile, partialPath, masterKey)
		if err != nil {
			os.Remove(partialPath)
			return fmt.Errorf("chunked encryption failed: %w", err)
		}
	} else {

		fmt.Println("Encrypting small file...")

		plaintext, err := io.ReadAll(inFile)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		nonce, ciphertext, err := crypto.Encrypt(masterKey, plaintext)
		if err != nil {
			return fmt.Errorf("encryption failed: %w", err)
		}

		outFile, err := os.OpenFile(partialPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
		if err != nil {
			return fmt.Errorf("can't make partial file: %w", err)
		}
		defer func() {
			outFile.Sync()
			outFile.Close()
		}()

		if _, err := outFile.Write(nonce); err != nil {
			os.Remove(partialPath)
			return fmt.Errorf("failed to write nonce: %w", err)
		}

		if _, err := outFile.Write(ciphertext); err != nil {
			os.Remove(partialPath)
			return fmt.Errorf("failed to write ciphertext: %w", err)
		}
	}

	if err := os.Rename(partialPath, finalPath); err != nil {
		os.Remove(partialPath)
		return fmt.Errorf("partial not renamed: %w", err)
	}

	return nil
}
