// Package model defines image management data structures.
package model

import "time"

// Image represents an image entity with progressive loading support.
type Image struct {
	ID           uint           `json:"id" gorm:"primaryKey;comment:主键"`
	OwnerID      uint           `json:"owner_id" gorm:"not null;index:idx_images_owner_id;comment:所有者ID"`
	Filename     string         `json:"filename" gorm:"size:255;not null;comment:随机文件名"`
	OriginalName string         `json:"original_name" gorm:"size:255;not null;comment:原始文件名"`
	MimeType     string         `json:"mime_type" gorm:"size:100;not null;comment:MIME类型"`
	Size         int64          `json:"size" gorm:"not null;comment:文件大小(字节)"`
	Width        int            `json:"width" gorm:"not null;comment:原始宽度"`
	Height       int            `json:"height" gorm:"not null;comment:原始高度"`
	Blurhash     string         `json:"blurhash,omitempty" gorm:"size:40;comment:BlurHash字符串"`
	StorageType  string         `json:"storage_type" gorm:"size:20;default:'local';comment:存储类型"`
	StoragePath  string         `json:"storage_path" gorm:"size:500;comment:存储路径"`
	IsPublic     bool           `json:"is_public" gorm:"default:false;comment:是否公开"`
	AccessToken  string         `json:"access_token" gorm:"size:64;comment:访问令牌"`
	Status       ImageStatus    `json:"status" gorm:"size:20;default:'processing';index:idx_images_status;comment:状态"`
	CreatedAt    time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt    time.Time      `json:"updated_at" gorm:"comment:更新时间"`

	// Associations
	Variants []ImageVariant  `json:"variants,omitempty" gorm:"foreignKey:ImageID;constraint:OnDelete:CASCADE"`
	Groups   []ImageGroup    `json:"groups,omitempty" gorm:"many2many:image_group_members;"`
}

// TableName returns the table name for Image.
func (Image) TableName() string {
	return "images"
}

// ImageStatus represents the processing status of an image.
type ImageStatus string

const (
	ImageStatusProcessing ImageStatus = "processing" // 处理中
	ImageStatusReady      ImageStatus = "ready"       // 就绪
	ImageStatusFailed     ImageStatus = "failed"      // 失败
)

// ImageVariantType represents the type of image variant.
type ImageVariantType string

const (
	VariantOriginal ImageVariantType = "original" // 原图
	VariantMedium   ImageVariantType = "medium"   // 中图
	VariantSmall    ImageVariantType = "small"    // 小图
)

// StorageType represents the storage backend type.
type StorageType string

const (
	StorageTypeLocal StorageType = "local" // 本地存储
	StorageTypeOSS   StorageType = "oss"   // 阿里云OSS
	StorageTypeS3    StorageType = "s3"    // AWS S3
	StorageTypeCOS   StorageType = "cos"   // 腾讯云COS
)

// ImageVariant represents a resized version of an image.
type ImageVariant struct {
	ID        uint             `json:"id" gorm:"primaryKey;comment:主键"`
	ImageID   uint             `json:"image_id" gorm:"not null;index:idx_variants_image_id;comment:图片ID"`
	Variant   ImageVariantType `json:"variant" gorm:"size:20;not null;index:idx_variants_variant;comment:变体类型"`
	Width     int              `json:"width" gorm:"not null;comment:宽度"`
	Height    int              `json:"height" gorm:"not null;comment:高度"`
	FileSize  int64            `json:"file_size" gorm:"not null;comment:文件大小"`
	FilePath  string           `json:"file_path" gorm:"size:500;not null;comment:文件路径"`
	CreatedAt time.Time        `json:"created_at" gorm:"comment:创建时间"`

	// Associations
	Image *Image `json:"-" gorm:"foreignKey:ImageID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for ImageVariant.
func (ImageVariant) TableName() string {
	return "image_variants"
}

// ImageGroupType represents the type of image grouping.
type ImageGroupType string

const (
	GroupTypeBook   ImageGroupType = "book"   // 图书分组
	GroupTypeAlbum  ImageGroupType = "album"  // 相册（预留）
	GroupTypeCustom ImageGroupType = "custom" // 自定义（预留）
)

// ImageGroup represents a collection of images.
type ImageGroup struct {
	ID        uint           `json:"id" gorm:"primaryKey;comment:主键"`
	OwnerID   uint           `json:"owner_id" gorm:"not null;index:idx_groups_owner_id;comment:所有者ID"`
	GroupType ImageGroupType `json:"group_type" gorm:"size:20;not null;index:idx_groups_type;comment:分组类型"`
	GroupName string         `json:"group_name" gorm:"size:255;comment:分组名称"`
	RefID     *uint          `json:"ref_id,omitempty" gorm:"index:idx_groups_ref_id;comment:关联ID(如book_id)"`
	RefType   *string        `json:"ref_type,omitempty" gorm:"size:50;comment:关联类型(如book)"`
	CreatedAt time.Time      `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"comment:更新时间"`

	// Associations
	Images []Image `json:"images,omitempty" gorm:"many2many:image_group_members;"`
}

// TableName returns the table name for ImageGroup.
func (ImageGroup) TableName() string {
	return "image_groups"
}

// ImageGroupMember represents the many-to-many relationship between images and groups.
type ImageGroupMember struct {
	GroupID   uint      `json:"group_id" gorm:"primaryKey;comment:分组ID"`
	ImageID   uint      `json:"image_id" gorm:"primaryKey;comment:图片ID"`
	SortOrder int       `json:"sort_order" gorm:"default:0;comment:排序"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`

	// Associations
	Group *ImageGroup `json:"-" gorm:"foreignKey:GroupID;constraint:OnDelete:CASCADE"`
	Image *Image      `json:"-" gorm:"foreignKey:ImageID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for ImageGroupMember.
func (ImageGroupMember) TableName() string {
	return "image_group_members"
}
