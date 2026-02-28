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
	GetByAccessToken(ctx context.Context, token string) (*model.Image, error)
	ListByOwner(ctx context.Context, ownerID uint, limit int, offset int) ([]*model.Image, int64, error)
	Update(ctx context.Context, image *model.Image) error
	Delete(ctx context.Context, id uint) error
	CreateVariant(ctx context.Context, variant *model.ImageVariant) error
	CreateGroup(ctx context.Context, group *model.ImageGroup) error
	GetGroupByID(ctx context.Context, id uint) (*model.ImageGroup, error)
	GetGroupByOwnerAndRef(ctx context.Context, ownerID uint, groupType model.ImageGroupType, refID uint, refType string) (*model.ImageGroup, error)
	AddImagesToGroup(ctx context.Context, groupID uint, imageIDs []uint) error
	RemoveImageFromGroup(ctx context.Context, groupID uint, imageID uint) error
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
		Preload("Groups").
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
		Preload("Groups").
		Where("id = ? AND owner_id = ?", id, ownerID).
		First(&image).Error
	if err != nil {
		return nil, err
	}
	return &image, nil
}

// GetByAccessToken retrieves an image using its access token.
func (r *ImageStore) GetByAccessToken(ctx context.Context, token string) (*model.Image, error) {
	var image model.Image
	err := r.db.Conn(ctx).
		Preload("Variants").
		Where("access_token = ?", token).
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

// CreateGroup creates a new image group.
func (r *ImageStore) CreateGroup(ctx context.Context, group *model.ImageGroup) error {
	return r.db.Conn(ctx).Create(group).Error
}

// GetGroupByID retrieves a group by ID.
func (r *ImageStore) GetGroupByID(ctx context.Context, id uint) (*model.ImageGroup, error) {
	var group model.ImageGroup
	err := r.db.Conn(ctx).
		Preload("Images").
		Preload("Images.Variants").
		First(&group, id).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// GetGroupByOwnerAndRef retrieves a group by owner and reference.
func (r *ImageStore) GetGroupByOwnerAndRef(ctx context.Context, ownerID uint, groupType model.ImageGroupType, refID uint, refType string) (*model.ImageGroup, error) {
	var group model.ImageGroup
	err := r.db.Conn(ctx).
		Where("owner_id = ? AND group_type = ? AND ref_id = ? AND ref_type = ?", ownerID, groupType, refID, refType).
		First(&group).Error
	if err != nil {
		return nil, err
	}
	return &group, nil
}

// AddImagesToGroup adds images to a group.
func (r *ImageStore) AddImagesToGroup(ctx context.Context, groupID uint, imageIDs []uint) error {
	return r.db.TX(ctx, func(txCtx context.Context) error {
		for _, imageID := range imageIDs {
			member := model.ImageGroupMember{
				GroupID: groupID,
				ImageID: imageID,
			}
			// Use OnConflict to handle duplicate entries
			if err := r.db.Conn(txCtx).Create(&member).Error; err != nil {
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
func (r *ImageStore) RemoveImageFromGroup(ctx context.Context, groupID uint, imageID uint) error {
	return r.db.Conn(ctx).
		Where("group_id = ? AND image_id = ?", groupID, imageID).
		Delete(&model.ImageGroupMember{}).Error
}

// isDuplicateKeyError checks if an error is a duplicate key error.
func isDuplicateKeyError(err error) bool {
	return err != nil && (err.Error() == "Error 1062: Duplicate entry" || err.Error() == "UNIQUE constraint failed: image_group_members.group_id, image_group_members.image_id")
}
