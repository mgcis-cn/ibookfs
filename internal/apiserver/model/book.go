// Package model defines data structures.
package model

import "time"

// Book represents a book entity.
type Book struct {
	ID            uint      `json:"id" gorm:"primaryKey;comment:主键"`
	UserID        uint      `json:"user_id" gorm:"not null;index:idx_books_user_id;comment:用户ID"`
	Title         string    `json:"title" gorm:"size:255;not null;comment:书名"`
	Author        string    `json:"author" gorm:"size:255;comment:作者"`
	ISBN          string    `json:"isbn" gorm:"size:20;uniqueIndex:idx_books_isbn;comment:ISBN编号"`
	Publisher     string    `json:"publisher" gorm:"size:255;comment:出版社"`
	Year          int       `json:"year" gorm:"comment:出版年份"`
	Pages         int       `json:"pages" gorm:"comment:总页数"`
	UploadedPages int       `json:"uploaded_pages" gorm:"default:0;comment:已上传页数"`
	Status        string    `json:"status" gorm:"size:20;default:'draft';comment:状态"`
	Cover         string    `json:"cover" gorm:"size:500;comment:封面图片路径"`
	CreatedAt     time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"comment:更新时间"`

	// Associations
	Images []BookImage `json:"images,omitempty" gorm:"foreignKey:BookID;constraint:OnDelete:CASCADE"`
}

// TableName returns the table name for Book.
func (*Book) TableName() string {
	return "books"
}
