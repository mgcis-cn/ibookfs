// Package model defines image management data structures.
package model

import "time"

// Image represents an image entity with progressive loading support.
type Image struct {
	ID           uint      `json:"id" gorm:"primaryKey;comment:主键"`
	OwnerID      uint      `json:"owner_id" gorm:"not null;index:idx_images_owner_id;comment:所有者ID"`
	Filename     string    `json:"filename" gorm:"size:255;not null;comment:随机文件名"`
	OriginalName string    `json:"original_name" gorm:"size:255;not null;comment:原始文件名"`
	MimeType     string    `json:"mime_type" gorm:"size:100;not null;comment:MIME类型"`
	Size         int64     `json:"size" gorm:"not null;comment:文件大小(字节)"`
	Width        int       `json:"width" gorm:"not null;comment:原始宽度"`
	Height       int       `json:"height" gorm:"not null;comment:原始高度"`
	Blurhash     string    `json:"blurhash,omitempty" gorm:"size:40;comment:BlurHash字符串"`
	StorageType  string    `json:"storage_type" gorm:"size:20;default:'local';comment:存储类型"`
	StoragePath  string    `json:"storage_path" gorm:"size:500;comment:存储路径"`
	CreatedAt    time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"comment:更新时间"`

	// Associations
	Variants []ImageVariant `json:"variants,omitempty" gorm:"foreignKey:ImageID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for Image.
func (*Image) TableName() string {
	return "images"
}

const (
	VariantOriginal string = "original" // 原图
	VariantMedium   string = "medium"   // 中图
	VariantSmall    string = "small"    // 小图
)

const (
	StorageTypeLocal string = "local" // 本地存储
	StorageTypeOSS   string = "oss"   // 阿里云OSS
	StorageTypeS3    string = "s3"    // AWS S3
	StorageTypeCOS   string = "cos"   // 腾讯云COS
)

// ImageVariant represents a resized version of an image.
type ImageVariant struct {
	ID        uint      `json:"id" gorm:"primaryKey;comment:主键"`
	ImageID   uint      `json:"image_id" gorm:"not null;index:idx_variants_image_id;comment:图片ID"`
	Variant   string    `json:"variant" gorm:"size:20;not null;index:idx_variants_variant;comment:变体类型"`
	Width     int       `json:"width" gorm:"not null;comment:宽度"`
	Height    int       `json:"height" gorm:"not null;comment:高度"`
	FileSize  int64     `json:"file_size" gorm:"not null;comment:文件大小"`
	FilePath  string    `json:"file_path" gorm:"size:500;not null;comment:文件路径"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`

	// Associations
	Image *Image `json:"-" gorm:"foreignKey:ImageID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for ImageVariant.
func (*ImageVariant) TableName() string {
	return "image_variants"
}

// BookImage represents the association between a book and an image.
type BookImage struct {
	ID        uint      `json:"id" gorm:"primaryKey;comment:主键"`
	BookID    uint      `json:"book_id" gorm:"not null;index:idx_book_image_book_id;comment:图书ID"`
	ImageID   uint      `json:"image_id" gorm:"not null;index:idx_book_image_image_id;comment:图片ID"`
	OwnerID   uint      `json:"owner_id" gorm:"not null;index:idx_book_image_owner_id;comment:所有者ID"`
	SortOrder int       `json:"sort_order" gorm:"default:0;comment:排序"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:更新时间"`

	// Associations
	Image *Image `json:"-" gorm:"foreignKey:ImageID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for BookImage.
func (BookImage) TableName() string {
	return "book_image"
}
