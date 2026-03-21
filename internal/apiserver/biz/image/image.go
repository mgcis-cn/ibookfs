// Package image provides image management business logic.
package image

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/mgcis-cn/ibookfs/internal/apiserver/biz/processor"
	apierr "github.com/mgcis-cn/ibookfs/internal/apiserver/errors"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/model"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/store"
	pkgerr "github.com/mgcis-cn/ibookfs/pkg/errors"
	ssources "github.com/mgcis-cn/ibookfs/pkg/storage/sources"
)

// ImageBiz defines the interface for image management business logic.
type ImageBiz interface {
	Upload(ctx context.Context, req *UploadRequest) (*UploadResponse, error)
	GetByID(ctx context.Context, id uint, ownerID uint) (*model.Image, error)
	List(ctx context.Context, ownerID uint, page, pageSize int) ([]*model.Image, int64, error)
	Delete(ctx context.Context, id uint, ownerID uint) error
	ProcessImage(ctx context.Context, imageID uint) error
	DownloadFile(ctx context.Context, image *model.Image, variant string) (ssources.ReadCloser, string, error)
	ListByBookID(ctx context.Context, bookID uint, page, pageSize int) ([]*model.Image, int64, error)
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
	BookID      *uint
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
	URL          string `json:"url"`
}

// Upload handles image upload and processing.
func (s *imageBiz) Upload(ctx context.Context, req *UploadRequest) (*UploadResponse, error) {
	// Validate file extension
	if err := processor.ValidateFileExtension(req.FileName); err != nil {
		return nil, err
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
		return nil, pkgerr.Wrap(err, "读取文件失败")
	}

	// Validate image
	info, err := s.processor.Validate(ctx, strings.NewReader(string(content)), mimeType, req.FileSize)
	if err != nil {
		return nil, err
	}

	// Generate random filename
	ext := filepath.Ext(req.FileName)
	randomFilename, err := generateRandomFilename(ext)
	if err != nil {
		return nil, pkgerr.Wrap(err, "生成文件名失败")
	}

	// Generate storage path (code plans the structure)
	storagePath := fmt.Sprintf("images/book/%d/original/%s", req.OwnerID, randomFilename)

	// Save original image
	fullPath := s.getStorageBasePath()
	destPath := filepath.Join(fullPath, storagePath)
	if err := s.processor.SaveOriginal(strings.NewReader(string(content)), destPath); err != nil {
		return nil, apierr.ImageUploadFailed(err)
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
	}

	if err := s.repo.Image().Create(ctx, image); err != nil {
		// Cleanup file on database error
		_ = s.storage.Delete(ctx, storagePath)
		return nil, pkgerr.Wrap(err, "创建图片记录失败")
	}

	// Associate image with book if BookID is provided
	if req.BookID != nil {
		if err := s.repo.Image().AddImageToBook(ctx, image.ID, req.OwnerID, *req.BookID); err != nil {
			fmt.Printf("warning: failed to add image %d to book %d: %v\n", image.ID, *req.BookID, err)
		} else {
			// Update book uploaded pages count
			if err := s.repo.Book().UpdateUploadedPages(ctx, *req.BookID); err != nil {
				fmt.Printf("warning: failed to update uploaded pages for book %d: %v\n", *req.BookID, err)
			}
		}

		// Set book cover if empty (first uploaded image becomes cover)
		coverURL := s.storage.GetURL(storagePath)
		if err := s.repo.Book().SetCoverIfEmpty(ctx, *req.BookID, req.OwnerID, coverURL); err != nil {
			fmt.Printf("warning: failed to set cover for book %d: %v\n", *req.BookID, err)
		}
	}

	// Enqueue for async processing
	if err := s.worker.Enqueue(image.ID); err != nil {
		// Worker not available, log warning but don't fail the upload
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
		return nil, apierr.ErrImageAccessDenied
	}

	return image, nil
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
		return apierr.ErrImageAccessDenied
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
		return pkgerr.Wrap(err, "获取图片失败")
	}

	// Get full path for processing
	fullPath := s.getStorageBasePath()
	sourcePath := filepath.Join(fullPath, image.StoragePath)

	// Process the image
	result, err := s.processor.Process(ctx, sourcePath, filepath.Dir(sourcePath))
	if err != nil {
		return pkgerr.Wrap(err, "图片处理失败")
	}

	// Update image with blurhash
	image.Blurhash = result.Blurhash

	// Create variant records
	for _, v := range result.Variants {
		// Convert absolute path to relative storage path
		relPath, _ := filepath.Rel(fullPath, v.Path)

		image.Variants = append(image.Variants, model.ImageVariant{
			ImageID:  image.ID,
			Variant:  v.Name,
			Width:    v.Width,
			Height:   v.Height,
			FileSize: v.FileSize,
			FilePath: relPath,
		})
	}

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
		return nil, "", apierr.ImageDownloadFailed(err)
	}

	return reader, image.MimeType, nil
}

// ListByBookID retrieves images associated with a book.
func (s *imageBiz) ListByBookID(ctx context.Context, bookID uint, page, pageSize int) ([]*model.Image, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.Image().ListByBookID(ctx, bookID, pageSize, offset)
}

// getStorageBasePath returns the base path for local storage.
func (s *imageBiz) getStorageBasePath() string {
	return s.storage.GetBasePath()
}

// generateRandomFilename generates a random filename with the given extension.
func generateRandomFilename(ext string) (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b) + ext, nil
}
