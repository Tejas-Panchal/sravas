package service

import (
	"io"
	"os"
	"path/filepath"
)

// Storage defines the interface for video file storage (local or S3)
type Storage interface {
	Save(filename string, reader io.Reader) error
	Get(path string) (io.ReadCloser, error)
	Delete(path string) error
}

// LocalStorage stores video files on the local filesystem
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a LocalStorage rooted at the given directory
func NewLocalStorage(basePath string) *LocalStorage {
	os.MkdirAll(basePath, 0755)
	return &LocalStorage{basePath: basePath}
}

// Save writes an uploaded file to the local filesystem
func (s *LocalStorage) Save(filename string, reader io.Reader) error {
	path := filepath.Join(s.basePath, filename)
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(f, reader)
	return err
}

// Get opens a file from the local filesystem for reading
func (s *LocalStorage) Get(path string) (io.ReadCloser, error) {
	return os.Open(filepath.Join(s.basePath, path))
}

// Delete removes a file from the local filesystem
func (s *LocalStorage) Delete(path string) error {
	return os.Remove(filepath.Join(s.basePath, path))
}
