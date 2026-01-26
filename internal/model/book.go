// Package model defines data structures.
package model

import "time"

// Book represents a book entity.
type Book struct {
	ID            uint      `json:"id" gorm:"primaryKey"`
	Title         string    `json:"title" gorm:"size:255;not null"`
	Author        string    `json:"author" gorm:"size:255"`
	ISBN          string    `json:"isbn" gorm:"size:20;uniqueIndex"`
	Publisher     string    `json:"publisher" gorm:"size:255"`
	Year          int       `json:"year"`
	Pages         int       `json:"pages"`
	UploadedPages int       `json:"uploaded_pages" gorm:"default:0"`
	Status        string    `json:"status" gorm:"size:20;default:'draft'"`
	FilePath      string    `json:"file_path" gorm:"size:500"`
	FileSize      int64     `json:"file_size"`
	Format        string    `json:"format" gorm:"size:20"`
	Cover         string    `json:"cover" gorm:"size:500"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// TableName returns the table name for Book.
func (Book) TableName() string {
	return "books"
}
