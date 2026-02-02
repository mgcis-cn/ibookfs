// Package storage provides abstracted file storage interfaces.
package storage

import "context"

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

	// GetType returns the storage type identifier.
	GetType() string
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

// Config holds storage configuration.
type Config struct {
	// Default is the default storage type to use.
	Default string

	// Local holds local filesystem storage configuration.
	Local LocalConfig

	// OSS holds Alibaba Cloud OSS configuration.
	OSS OSSConfig

	// S3 holds AWS S3 configuration.
	S3 S3Config

	// COS holds Tencent Cloud COS configuration.
	COS COSConfig
}

// LocalConfig defines local filesystem storage configuration.
type LocalConfig struct {
	BasePath string // Base directory for file storage
	BaseURL  string // Base URL for file access
}

// OSSConfig defines Alibaba Cloud OSS configuration.
type OSSConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
}

// S3Config defines AWS S3 configuration.
type S3Config struct {
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	Endpoint        string // Optional, for S3-compatible services
}

// COSConfig defines Tencent Cloud COS configuration.
type COSConfig struct {
	SecretID     string
	SecretKey    string
	BucketName   string
	BucketRegion string
}
