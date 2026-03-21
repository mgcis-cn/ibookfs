// Package repository provides image data access layer.
package store

import (
	"context"

	"github.com/mgcis-cn/ibookfs/internal/apiserver/model"
	"github.com/mgcis-cn/ibookfs/pkg/database"
)

// IImageStore defines the interface for image data operations.
type IImageStore interface {
	Create(ctx context.Context, image *model.Image) error
	GetByID(ctx context.Context, id uint) (*model.Image, error)
	GetByIDAndOwner(ctx context.Context, id uint, ownerID uint) (*model.Image, error)
	ListByOwner(ctx context.Context, ownerID uint, limit int, offset int) ([]*model.Image, int64, error)
	Update(ctx context.Context, image *model.Image) error
	Delete(ctx context.Context, id uint) error
	CreateVariant(ctx context.Context, variant *model.ImageVariant) error
	AddImageToBook(ctx context.Context, imageID uint, ownerID uint, bookID uint) error
	RemoveImageFromBook(ctx context.Context, imageID uint, bookID uint) error
	ListByBookID(ctx context.Context, bookID uint, limit int, offset int) ([]*model.Image, int64, error)
}

// ImageStore handles image data operations.
type ImageStore struct {
	db database.Database
}

// NewImageStore creates a new image store with injected DatabaseOptions.
func NewImageStore(db database.Database) *ImageStore {
	return &ImageStore{db: db}
}

// Create creates a new image record.
func (r *ImageStore) Create(ctx context.Context, image *model.Image) error {
	return r.db.Conn(ctx).Create(image).Error
}

// GetByID retrieves an image by ID.
func (r *ImageStore) GetByID(ctx context.Context, id uint) (*model.Image, error) {
	var image model.Image
	err := r.db.Conn(ctx).
		Preload("Variants").
		First(&image, id).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// GetByIDAndOwner retrieves an image by ID and owner ID.
func (r *ImageStore) GetByIDAndOwner(ctx context.Context, id uint, ownerID uint) (*model.Image, error) {
	var image model.Image
	err := r.db.Conn(ctx).
		Preload("Variants").
		Where("id = ? AND owner_id = ?", id, ownerID).
		First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// ListByOwner retrieves images for an owner with pagination.
func (r *ImageStore) ListByOwner(ctx context.Context, ownerID uint, limit int, offset int) ([]*model.Image, int64, error) {
	var images []*model.Image
	var total int64

	db := r.db.Conn(ctx).Model(&model.Image{})

	// Filter by owner
	db = db.Where("owner_id = ?", ownerID)

	// Count total
	db.Count(&total)

	// Fetch with pagination
	err := db.
		Preload("Variants").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&images).Error

	if err != nil {
		return nil, 0, err
	}

	return images, total, nil
}

// Update updates an existing image record.
func (r *ImageStore) Update(ctx context.Context, image *model.Image) error {
	return r.db.Conn(ctx).Save(image).Error
}

// Delete deletes an image by ID.
func (r *ImageStore) Delete(ctx context.Context, id uint) error {
	return r.db.Conn(ctx).Delete(&model.Image{}, id).Error
}

// CreateVariant creates a new image variant record.
func (r *ImageStore) CreateVariant(ctx context.Context, variant *model.ImageVariant) error {
	return r.db.Conn(ctx).Create(variant).Error
}

// AddImageToBook adds an image to a book.
func (r *ImageStore) AddImageToBook(ctx context.Context, imageID uint, ownerID uint, bookID uint) error {
	bookImage := model.BookImage{
		ImageID: imageID,
		BookID:  bookID,
		OwnerID: ownerID,
	}
	return r.db.Conn(ctx).Create(&bookImage).Error
}

// RemoveImageFromBook removes an image from a book.
func (r *ImageStore) RemoveImageFromBook(ctx context.Context, imageID uint, bookID uint) error {
	return r.db.Conn(ctx).
		Where("image_id = ? AND book_id = ?", imageID, bookID).
		Delete(&model.BookImage{}).Error
}

// ListByBookID retrieves images associated with a book via book_image.
func (r *ImageStore) ListByBookID(ctx context.Context, bookID uint, limit int, offset int) ([]*model.Image, int64, error) {
	var images []*model.Image
	var total int64

	// Count images in the book
	err := r.db.Conn(ctx).
		Model(&model.Image{}).
		Joins("JOIN book_image ON book_image.image_id = images.id").
		Where("book_image.book_id = ?", bookID).
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	if total == 0 {
		return []*model.Image{}, 0, nil
	}

	// Fetch images
	err = r.db.Conn(ctx).
		Joins("JOIN book_image ON book_image.image_id = images.id").
		Where("book_image.book_id = ?", bookID).
		Preload("Variants").
		Order("book_image.sort_order ASC, images.created_at ASC").
		Limit(limit).
		Offset(offset).
		Find(&images).Error

	if err != nil {
		return nil, 0, err
	}

	return images, total, nil
}
