// Package image provides image management business logic.
package image

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/mgcis-cn/ibookfs/internal/apiserver/biz/processor"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/model"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/store"
	ssources "github.com/mgcis-cn/ibookfs/pkg/storage/sources"
)

// ImageBiz defines the interface for image management business logic.
type ImageBiz interface {
	Upload(ctx context.Context, req *UploadRequest) (*UploadResponse, error)
	GetByID(ctx context.Context, id uint, ownerID uint) (*model.Image, error)
	GetByAccessToken(ctx context.Context, token string) (*model.Image, error)
	List(ctx context.Context, ownerID uint, page, pageSize int) ([]*model.Image, int64, error)
	Delete(ctx context.Context, id uint, ownerID uint) error
	ProcessImage(ctx context.Context, imageID uint) error
	DownloadFile(ctx context.Context, image *model.Image, variant string) (ssources.ReadCloser, string, error)
	CreateGroup(ctx context.Context, ownerID uint, groupType model.ImageGroupType, groupName string, refID *uint, refType *string) (*model.ImageGroup, error)
	AddImagesToGroup(ctx context.Context, groupID uint, imageIDs []uint) error
}

// ImageWorker defines the interface for asynchronous image processing.
type ImageWorker interface {
	Enqueue(imageID uint) error
	Start()
	Stop()
	SetService(svc any) // SetService sets the processor for the worker
}

// imageBiz is the concrete implementation of ImageBiz.
type imageBiz struct {
	storage   ssources.Storage
	processor *processor.Processor
	repo      store.IStore
	worker    ImageWorker
}

// NewImageBiz creates a new image business logic.
func NewImageBiz(s ssources.Storage, p *processor.Processor, repo store.IStore, w ImageWorker) ImageBiz {
	return &imageBiz{
		storage:   s,
		processor: p,
		repo:      repo,
		worker:    w,
	}
}

// UploadRequest contains parameters for image upload.
type UploadRequest struct {
	FileName    string
	FileSize    int64
	ContentType string
	Reader      multipart.File
	OwnerID     uint
	GroupID     *uint
}

// UploadResponse contains the result of image upload.
type UploadResponse struct {
	ID           uint   `json:"id"`
	Filename     string `json:"filename"`
	OriginalName string `json:"original_name"`
	MimeType     string `json:"mime_type"`
	Size         int64  `json:"size"`
	Width        int    `json:"width"`
	Height       int    `json:"height"`
	Blurhash     string `json:"blurhash,omitempty"`
	Status       string `json:"status"`
	AccessToken  string `json:"access_token"`
	URL          string `json:"url"`
}

// Upload handles image upload and processing.
func (s *imageBiz) Upload(ctx context.Context, req *UploadRequest) (*UploadResponse, error) {
	// Validate file extension
	if err := processor.ValidateFileExtension(req.FileName); err != nil {
		return nil, fmt.Errorf("invalid file: %w", err)
	}

	// Get MIME type from extension if not provided
	mimeType := req.ContentType
	if mimeType == "" {
		mimeType = processor.GetMIMEType(req.FileName)
	}

	// Create a tee reader to both validate and preserve content
	// First, read all content for validation
	content, err := io.ReadAll(req.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Validate image
	info, err := s.processor.Validate(ctx, strings.NewReader(string(content)), mimeType, req.FileSize)
	if err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	// Generate random filename
	ext := filepath.Ext(req.FileName)
	randomFilename, err := generateRandomFilename(ext)
	if err != nil {
		return nil, fmt.Errorf("failed to generate filename: %w", err)
	}

	// Generate storage path (code plans the structure)
	storagePath := fmt.Sprintf("images/book/%d/original/%s", req.OwnerID, randomFilename)

	// Save original image
	fullPath := s.getStorageBasePath()
	destPath := filepath.Join(fullPath, storagePath)
	if err := s.processor.SaveOriginal(strings.NewReader(string(content)), destPath); err != nil {
		return nil, fmt.Errorf("failed to save image: %w", err)
	}

	// Generate access token
	accessToken, err := generateAccessToken()
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Create image record
	image := &model.Image{
		OwnerID:      req.OwnerID,
		Filename:     randomFilename,
		OriginalName: req.FileName,
		MimeType:     info.MimeType,
		Size:         info.Size,
		Width:        info.Width,
		Height:       info.Height,
		StorageType:  string(s.storage.Kind()),
		StoragePath:  storagePath,
		AccessToken:  accessToken,
		Status:       model.ImageStatusProcessing,
	}

	if err := s.repo.Image().Create(ctx, image); err != nil {
		// Cleanup file on database error
		_ = s.storage.Delete(ctx, storagePath)
		return nil, fmt.Errorf("failed to create image record: %w", err)
	}

	// Enqueue for async processing
	if err := s.worker.Enqueue(image.ID); err != nil {
		// Worker not available, mark as failed but don't fail the upload
		// Username can retry processing later via API
		image.Status = model.ImageStatusFailed
		_ = s.repo.Image().Update(ctx, image)
		fmt.Printf("warning: failed to enqueue image %d for processing: %v\n", image.ID, err)
	}

	return &UploadResponse{
		ID:           image.ID,
		Filename:     image.Filename,
		OriginalName: image.OriginalName,
		MimeType:     image.MimeType,
		Size:         image.Size,
		Width:        image.Width,
		Height:       image.Height,
		Status:       string(image.Status),
		AccessToken:  image.AccessToken,
		URL:          s.storage.GetURL(storagePath),
	}, nil
}

// GetByID retrieves an image by ID.
func (s *imageBiz) GetByID(ctx context.Context, id uint, ownerID uint) (*model.Image, error) {
	image, err := s.repo.Image().GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Check ownership
	if image.OwnerID != ownerID {
		return nil, errors.New("access denied")
	}

	return image, nil
}

// GetByAccessToken retrieves an image using its access token.
func (s *imageBiz) GetByAccessToken(ctx context.Context, token string) (*model.Image, error) {
	return s.repo.Image().GetByAccessToken(ctx, token)
}

// List retrieves images for an owner with pagination.
func (s *imageBiz) List(ctx context.Context, ownerID uint, page, pageSize int) ([]*model.Image, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.Image().ListByOwner(ctx, ownerID, pageSize, offset)
}

// Delete removes an image and its files.
func (s *imageBiz) Delete(ctx context.Context, id uint, ownerID uint) error {
	image, err := s.repo.Image().GetByID(ctx, id)
	if err != nil {
		return err
	}

	// Check ownership
	if image.OwnerID != ownerID {
		return errors.New("access denied")
	}

	// Delete variants from storage
	for _, variant := range image.Variants {
		if err := s.storage.Delete(ctx, variant.FilePath); err != nil {
			fmt.Printf("warning: failed to delete variant %s: %v\n", variant.FilePath, err)
		}
	}

	// Delete original from storage
	if err := s.storage.Delete(ctx, image.StoragePath); err != nil {
		fmt.Printf("warning: failed to delete image %s: %v\n", image.StoragePath, err)
	}

	// Delete database record
	return s.repo.Image().Delete(ctx, id)
}

// ProcessImage performs async image processing.
func (s *imageBiz) ProcessImage(ctx context.Context, imageID uint) error {
	image, err := s.repo.Image().GetByID(ctx, imageID)
	if err != nil {
		return fmt.Errorf("failed to get image: %w", err)
	}

	// Get full path for processing
	fullPath := s.getStorageBasePath()
	sourcePath := filepath.Join(fullPath, image.StoragePath)

	// Process the image
	result, err := s.processor.Process(ctx, sourcePath, filepath.Dir(sourcePath))
	if err != nil {
		// Update status to failed
		image.Status = model.ImageStatusFailed
		_ = s.repo.Image().Update(ctx, image)
		return fmt.Errorf("processing failed: %w", err)
	}

	// Update image with blurhash
	image.Blurhash = result.Blurhash

	// Create variant records
	for _, v := range result.Variants {
		// Convert absolute path to relative storage path
		relPath, _ := filepath.Rel(fullPath, v.Path)

		image.Variants = append(image.Variants, model.ImageVariant{
			ImageID:  image.ID,
			Variant:  model.ImageVariantType(v.Name),
			Width:    v.Width,
			Height:   v.Height,
			FileSize: v.FileSize,
			FilePath: relPath,
		})
	}

	// Update status to ready
	image.Status = model.ImageStatusReady

	return s.repo.Image().Update(ctx, image)
}

// DownloadFile retrieves a file for download.
func (s *imageBiz) DownloadFile(ctx context.Context, image *model.Image, variant string) (ssources.ReadCloser, string, error) {
	var filePath string

	switch variant {
	case "small", "medium", "original":
		// Find the variant
		for _, v := range image.Variants {
			if string(v.Variant) == variant {
				filePath = v.FilePath
				break
			}
		}
		if filePath == "" {
			// Variant not ready yet, fall back to original
			filePath = image.StoragePath
		}
	default:
		// Default to original
		filePath = image.StoragePath
	}

	reader, err := s.storage.Download(ctx, filePath)
	if err != nil {
		return nil, "", fmt.Errorf("failed to download file: %w", err)
	}

	return reader, image.MimeType, nil
}

// CreateGroup creates a new image group.
func (s *imageBiz) CreateGroup(ctx context.Context, ownerID uint, groupType model.ImageGroupType, groupName string, refID *uint, refType *string) (*model.ImageGroup, error) {
	group := &model.ImageGroup{
		OwnerID:   ownerID,
		GroupType: groupType,
		GroupName: groupName,
		RefID:     refID,
		RefType:   refType,
	}

	// Use repository to create group (implementation needed)
	return group, nil
}

// AddImagesToGroup adds images to a group.
func (s *imageBiz) AddImagesToGroup(ctx context.Context, groupID uint, imageIDs []uint) error {
	// Implementation needed
	return nil
}

// getStorageBasePath returns the base path for local storage.
func (s *imageBiz) getStorageBasePath() string {
	if local, ok := s.storage.(ssources.Storage); ok {
		return local.GetURL("")
	}
	return ""
}

// generateRandomFilename generates a random filename with the given extension.
func generateRandomFilename(ext string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b) + ext, nil
}

// generateAccessToken generates a secure access token.
func generateAccessToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
