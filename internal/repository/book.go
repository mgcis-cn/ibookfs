// Package repository provides data access layer.
package repository

import (
	"context"

	"github.com/mgcis/ibookfs/internal/database"
	"github.com/mgcis/ibookfs/internal/model"
)

// BookRepository handles book data operations.
type BookRepository struct{}

// Create creates a new book.
func (r *BookRepository) Create(ctx context.Context, book *model.Book) error {
	return database.Default().WithContext(ctx).Create(book).Error
}

// GetByID retrieves a book by ID.
func (r *BookRepository) GetByID(ctx context.Context, id uint) (*model.Book, error) {
	var book model.Book
	err := database.Default().WithContext(ctx).First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// Update updates an existing book.
func (r *BookRepository) Update(ctx context.Context, book *model.Book) error {
	return database.Default().WithContext(ctx).Save(book).Error
}

// Delete deletes a book by ID.
func (r *BookRepository) Delete(ctx context.Context, id uint) error {
	return database.Default().WithContext(ctx).Delete(&model.Book{}, id).Error
}

// ListQuery defines query parameters for listing books.
type ListQuery struct {
	Page     int
	PageSize int
	Status   string
	Search   string
}

// List retrieves all books with pagination and filtering.
func (r *BookRepository) List(ctx context.Context, query ListQuery) ([]model.Book, int64, error) {
	var books []model.Book
	var total int64

	db := database.Default().WithContext(ctx).Model(&model.Book{})

	// Apply status filter
	if query.Status != "" {
		db = db.Where("status = ?", query.Status)
	}

	// Apply search filter
	if query.Search != "" {
		search := "%" + query.Search + "%"
		db = db.Where("title LIKE ? OR author LIKE ?", search, search)
	}

	db.Count(&total)

	offset := (query.Page - 1) * query.PageSize
	err := db.Offset(offset).Limit(query.PageSize).Find(&books).Error
	if err != nil {
		return nil, 0, err
	}

	return books, total, nil
}
