package sources

import (
	"context"
	"fmt"
)

type Kind string

// Config is the interface that storage configuration structs must implement.
type Config interface {
	Kind() string
}

type Factory func(cfg Config) (Storage, error)

var registry = make(map[Kind]Factory)

func Register(kind Kind, factory Factory) bool {
	if _, exists := registry[kind]; exists {
		return false
	}
	registry[kind] = factory
	return true
}

func Create(cfg Config) (Storage, error) {
	kind := Kind(cfg.Kind())
	factory, found := registry[kind]
	if !found {
		return nil, fmt.Errorf("unknown storage kind: %q", kind)
	}
	return factory(cfg)
}

// Storage defines the interface for file storage operations.
// Implementations can support local filesystem, cloud storage (OSS, S3, COS, etc.)
type Storage interface {
	// Upload uploads a file to the storage backend.
	// The path is relative to the storage root.
	Upload(ctx context.Context, reader Reader, path string) error

	// Download retrieves a file from the storage backend.
	Download(ctx context.Context, path string) (ReadCloser, error)

	// Delete removes a file from the storage backend.
	Delete(ctx context.Context, path string) error

	// Exists checks if a file exists in the storage backend.
	Exists(ctx context.Context, path string) (bool, error)

	// GetURL returns a publicly accessible URL for the file.
	// For private storage, this should return a signed URL.
	GetURL(path string) string

	// GetBasePath returns the base path for local storage.
	// For cloud storage, this returns an empty string.
	GetBasePath() string

	// Kind returns the storage type identifier.
	Kind() Kind
}

// Reader wraps the io.Reader interface for type clarity.
type Reader interface {
	Read(p []byte) (n int, err error)
}

// ReadCloser wraps the io.ReadCloser interface for type clarity.
type ReadCloser interface {
	Read(p []byte) (n int, err error)
	Close() error
}
