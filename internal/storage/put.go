package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

const (
	storageDir = "crypta"             
	bufferSize = 128 * 1024         
)

func Put(sourceFile string) error {

	err := os.MkdirAll(storageDir, 0700)
	if err != nil {
		return fmt.Errorf("no dir created: %w", err)
	}


	inFile, err := os.Open(sourceFile)
	if err != nil {
		return fmt.Errorf("can't open source file: %w", err)
	}
	defer inFile.Close()


	destFilename := filepath.Base(sourceFile) + ".crypta"
	partialPath := filepath.Join(storageDir, destFilename+".partial")
	finalPath := filepath.Join(storageDir, destFilename)


	outFile, err := os.Create(partialPath)
	if err != nil {
		return fmt.Errorf("can't make partial file: %w", err)
	}

	defer func() {
		if closeErr := outFile.Close(); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close partial file: %w", closeErr)
		}
	}()

	buf := make([]byte, bufferSize)
	_, err = io.CopyBuffer(outFile, inFile, buf)
	if err != nil {
		return fmt.Errorf("data couldn't be copied: %w", err)
	}


	err = outFile.Sync()
	if err != nil {
		return fmt.Errorf("partial didn't sync all the way: %w", err)
	}


	if err != nil {
		return fmt.Errorf("partial not closed: %w", err)
	}


	err = os.Rename(partialPath, finalPath)
	if err != nil {
		return fmt.Errorf("partial not renamed: %w", err)
	}

	return nil
}
