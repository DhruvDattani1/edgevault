package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type FileInfo struct {
	Name     string
	Size     int64
	Modified time.Time
}

func Delete(objectName string) error {
	path := filepath.Join(storageDir, objectName+".crypta")
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

func List() ([]FileInfo, error) {
	files, err := os.ReadDir(storageDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read storage dir: %w", err)
	}

	var list []FileInfo
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".crypta" {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		list = append(list, FileInfo{
			Name:     file.Name()[:len(file.Name())-7], // strip .crypta
			Size:     info.Size(),
			Modified: info.ModTime(),
		})
	}

	return list, nil
}
