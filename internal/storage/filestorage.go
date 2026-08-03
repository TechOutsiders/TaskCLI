package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileStorage provides persistent storage backed by a file.
//
// FileStorage is not safe for concurrent use by multiple goroutines
// or processes. Synchronization, if required, should be handled by
// a higher layer or by a future storage implementation.
//
// TODO: Consider performing atomic writes using a temporary file
// and rename.
type FileStorage struct {
	path string
}

// NewFileStorage creates a new file storage.
func NewFileStorage(path string) (*FileStorage, error) {
	const fileExtJSON = ".json"

	ext := filepath.Ext(path)
	if ext != fileExtJSON {
		return nil, fmt.Errorf("invalid file extension, got: %s want: %s", ext, fileExtJSON)
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			return nil, fmt.Errorf("creating storage file: %w", err)
		}

		err = file.Close()
		if err != nil {
			return nil, fmt.Errorf("closing storage file: %w", err)
		}
	}

	return &FileStorage{path: path}, nil
}

// Save writes data to the storage file.
func (f *FileStorage) Save(data []byte) error {
	err := os.WriteFile(f.path, data, 0o644)
	if err != nil {
		return fmt.Errorf("writing to file: %w", err)
	}

	return nil
}

// Load loads data from the storage file.
func (f *FileStorage) Load() (data []byte, err error) {
	data, err = os.ReadFile(f.path)
	if err != nil {
		return nil, fmt.Errorf("reading from file: %w", err)
	}

	return data, nil
}
