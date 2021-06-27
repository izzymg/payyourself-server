package storage

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestFileSystemFetch(t *testing.T) {
	testFilePath := "test_file"
	testText := "aaccbb"

	testDir, err := os.MkdirTemp("", "fetch")
	if err != nil {
		t.Errorf("failed to make test dir: %w", err)
	}

	defer func() {
		err := os.RemoveAll(testDir)
		if err != nil {
			t.Errorf("failed to remove test dir: %w", err)
		}
	}()

	testFile, err := os.Create(filepath.Join(testDir, testFilePath))
	if err != nil {
		t.Errorf("failed to create test file: %w", err)
	}
	defer testFile.Close()

	_, err = testFile.WriteString(testText)
	if err != nil {
		t.Errorf("failed to write to test file: %w", err)
	}

	storer := FileSystemStorer{
		root: testDir,
	}

	reader, err := storer.Fetch(testFilePath)
	if err != nil {
		t.Error(err)
	}
	defer reader.Close()
	bytes, err := ioutil.ReadAll(reader)
	if err != nil {
		t.Error(err)
	}

	if read := string(bytes); read != testText {
		t.Errorf("expected to read %s, got %s", testText, read)
	}
}
