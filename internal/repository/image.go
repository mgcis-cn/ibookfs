// Package repository provides image data access layer.
package repository

import (
	"context"

	"github.com/mgcis-cn/ibookfs/internal/database"
	"github.com/mgcis-cn/ibookfs/internal/model"
	"gorm.io/gorm"
)

// ImageRepository handles image data operations.
type ImageRepository struct{}

// NewImageRepository creates a new image repository.
func NewImageRepository() *ImageRepository {
	return &ImageRepository{}
}

// Create creates a new image record.
func (r *ImageRepository) Create(ctx context.Context, image *model.Image) error {
	return database.Default().WithContext(ctx).Create(image).Error
}

// GetByID retrieves an image by ID.
func (r *ImageRepository) GetByID(ctx context.Context, id uint) (*model.Image, error) {
	var image model.Image
	err := database.Default().WithContext(ctx).
		Preload("Variants").
		Preload("Groups").
		First(&image, id).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// GetByIDAndOwner retrieves an image by ID and owner ID.
func (r *ImageRepository) GetByIDAndOwner(ctx context.Context, id uint, ownerID uint) (*model.Image, error) {
	var image model.Image
	err := database.Default().WithContext(ctx).
		Preload("Variants").
		Preload("Groups").
		Where("id = ? AND owner_id = ?", id, ownerID).
		First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// GetByAccessToken retrieves an image using its access token.
func (r *ImageRepository) GetByAccessToken(ctx context.Context, token string) (*model.Image, error) {
	var image model.Image
	err := database.Default().WithContext(ctx).
		Preload("Variants").
		Where("access_token = ?", token).
		First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// ListByOwner retrieves images for an owner with pagination.
func (r *ImageRepository) ListByOwner(ctx context.Context, ownerID uint, limit int, offset int) ([]*model.Image, int64, error) {
	var images []*model.Image
	var total int64

	db := database.Default().WithContext(ctx).Model(&model.Image{})

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
func (r *ImageRepository) Update(ctx context.Context, image *model.Image) error {
	return database.Default().WithContext(ctx).Save(image).Error
}

// Delete deletes an image by ID.
func (r *ImageRepository) Delete(ctx context.Context, id uint) error {
	return database.Default().WithContext(ctx).Delete(&model.Image{}, id).Error
}

// CreateVariant creates a new image variant record.
func (r *ImageRepository) CreateVariant(ctx context.Context, variant *model.ImageVariant) error {
	return database.Default().WithContext(ctx).Create(variant).Error
}

// CreateGroup creates a new image group.
func (r *ImageRepository) CreateGroup(ctx context.Context, group *model.ImageGroup) error {
	return database.Default().WithContext(ctx).Create(group).Error
}

// GetGroupByID retrieves a group by ID.
func (r *ImageRepository) GetGroupByID(ctx context.Context, id uint) (*model.ImageGroup, error) {
	var group model.ImageGroup
	err := database.Default().WithContext(ctx).
		Preload("Images").
		Preload("Images.Variants").
		First(&group, id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// GetGroupByOwnerAndRef retrieves a group by owner and reference.
func (r *ImageRepository) GetGroupByOwnerAndRef(ctx context.Context, ownerID uint, groupType model.ImageGroupType, refID uint, refType string) (*model.ImageGroup, error) {
	var group model.ImageGroup
	err := database.Default().WithContext(ctx).
		Where("owner_id = ? AND group_type = ? AND ref_id = ? AND ref_type = ?", ownerID, groupType, refID, refType).
		First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// AddImagesToGroup adds images to a group.
func (r *ImageRepository) AddImagesToGroup(ctx context.Context, groupID uint, imageIDs []uint) error {
	db := database.Default().WithContext(ctx)
	return db.Transaction(func(tx *gorm.DB) error {
		for _, imageID := range imageIDs {
			member := model.ImageGroupMember{
				GroupID: groupID,
				ImageID: imageID,
			}
			// Use OnConflict to handle duplicate entries
			if err := tx.Create(&member).Error; err != nil {
				// Check if it's a duplicate key error
				if isDuplicateKeyError(err) {
					continue
				}
				return err
			}
		}
		return nil
	})
}

// RemoveImageFromGroup removes an image from a group.
func (r *ImageRepository) RemoveImageFromGroup(ctx context.Context, groupID uint, imageID uint) error {
	return database.Default().WithContext(ctx).
		Where("group_id = ? AND image_id = ?", groupID, imageID).
		Delete(&model.ImageGroupMember{}).Error
}

// isDuplicateKeyError checks if an error is a duplicate key error.
func isDuplicateKeyError(err error) bool {
	return err != nil && (err.Error() == "Error 1062: Duplicate entry" || err.Error() == "UNIQUE constraint failed: image_group_members.group_id, image_group_members.image_id")
}
