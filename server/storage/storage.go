package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
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
		return nil, fmt.Errorf("failed to open file for user ID %s: %w", userID, err)
	}

	return file, nil
}
