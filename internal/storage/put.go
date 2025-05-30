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
	//don't end up using the buffer anyways since I am doing ReadAll(), this is also a Poly1305 constraint
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

	defer outFile.Close() //this way I know my file will be closed, more idiomatic, in case of function panic

	var writeErr error
	defer func() {
		if writeErr != nil {
			os.Remove(partialPath)
		}
	}()

	// using a defer to make sure that the file is discarded as well if there is an error (all or nothing)

	if _, writeErr = outFile.Write(nonce); writeErr != nil {
		return fmt.Errorf("failed to write nonce: %w", writeErr)
	}

	if _, writeErr = outFile.Write(ciphertext); writeErr != nil {
		return fmt.Errorf("failed to write ciphertext: %w", writeErr)
	}

	if writeErr = outFile.Sync(); writeErr != nil {
		return fmt.Errorf("partial didn't sync: %w", writeErr)
	}

	if writeErr = os.Rename(partialPath, finalPath); writeErr != nil {
		return fmt.Errorf("partial not renamed: %w", writeErr)
	}

	return nil
}