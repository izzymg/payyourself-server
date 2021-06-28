package storage

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"py-server/server"
)

// FileSystemStorer is a UserSaveStorer which writes and retrieves saves
// using the file system
type FileSystemStorer struct {
	root string
}

func (fss FileSystemStorer) Fetch(userID string) (io.ReadCloser, error) {
	fp := filepath.Join(fss.root, userID)
	file, err := os.Open(fp)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return nil, server.ErrNoUserSave
		}
		return nil, fmt.Errorf("failed to open file for user ID %s: %w", userID, err)
	}

	return file, nil
}

// MakeFileSystemStorer returns a new FileSystemStorer using the given root,
// trying open it for writing or create it if it doesn't exist
func MakeFileSystemStorer(root string) (*FileSystemStorer, error) {

	info, err := os.Stat(root)
	if err != nil {
		return nil, fmt.Errorf("failed to stat root: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("provided root is not a directory: %w", err)
	}

	return &FileSystemStorer{
		root: root,
	}, nil
}
