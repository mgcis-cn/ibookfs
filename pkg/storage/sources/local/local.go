package local

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mgcis-cn/ibookfs/pkg/options"
	"github.com/mgcis-cn/ibookfs/pkg/storage/sources"
)

const SourceKind sources.Kind = "local"

func init() {
	sources.Register(SourceKind, func(cfg sources.Config) (sources.Storage, error) {
		opts, ok := cfg.(*options.LocalOSOptions)
		if !ok {
			return nil, fmt.Errorf("local: invalid config type %T", cfg)
		}
		bucket := opts.Bucket
		if bucket == "" {
			bucket = "./storages"
		}
		return New(opts.Endpoint, bucket)
	})
}

type Storage struct {
	Endpoint string `json:"endpoint"`
	Bucket   string `json:"bucket"`
}

func New(endpoint string, bucket string) (*Storage, error) {
	if bucket == "" {
		bucket = "./storages"
	}
	if err := os.MkdirAll(bucket, 0755); err != nil {
		return nil, err
	}
	return &Storage{
		Endpoint: endpoint,
		Bucket:   bucket,
	}, nil
}

func (s *Storage) Kind() sources.Kind {
	return SourceKind
}

// Upload saves a file to the local filesystem.
func (s *Storage) Upload(ctx context.Context, reader sources.Reader, path string) error {
	fullPath := filepath.Join(s.Bucket, path)

	// Create directory if not exists
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	// Create file
	file, err := os.Create(fullPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Copy content
	if _, err := io.Copy(file, reader); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}

// Download retrieves a file from the local filesystem.
func (s *Storage) Download(ctx context.Context, path string) (sources.ReadCloser, error) {
	fullPath := filepath.Join(s.Bucket, path)

	file, err := os.Open(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	return file, nil
}

// Delete removes a file from the local filesystem.
func (s *Storage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.Bucket, path)

	if err := os.Remove(fullPath); err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}

	return nil
}

// Exists checks if a file exists in the local filesystem.
func (s *Storage) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(s.Bucket, path)

	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetURL returns the publicly accessible URL for the file.
func (s *Storage) GetURL(path string) string {
	if s.Bucket[0] == '.' {
		return s.Endpoint + s.Bucket[1:] + "/" + path
	}
	if s.Bucket[0] == '/' {
		return s.Endpoint + s.Bucket + "/" + path
	}
	return s.Endpoint + "/" + s.Bucket + "/" + path
}
