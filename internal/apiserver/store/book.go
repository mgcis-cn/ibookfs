// Package repository provides data access layer.
package store

import (
	"context"

	"github.com/mgcis-cn/ibookfs/internal/apiserver/model"
	"github.com/mgcis-cn/ibookfs/pkg/database"
)

// IBookStore defines the interface for book data operations.
type IBookStore interface {
	Create(ctx context.Context, book *model.Book) error
	GetByID(ctx context.Context, id uint, userID uint) (*model.Book, error)
	Update(ctx context.Context, book *model.Book) error
	Delete(ctx context.Context, id uint, userID uint) error
	List(ctx context.Context, query ListQuery) ([]model.Book, int64, error)
}

// BookStore handles book data operations.
type BookStore struct {
	db database.Database
}

// NewBookStore creates a new BookStore with injected DatabaseOptions.
func NewBookStore(db database.Database) *BookStore {
	return &BookStore{db: db}
}

// Create creates a new book.
func (r *BookStore) Create(ctx context.Context, book *model.Book) error {
	return r.db.Conn(ctx).Create(book).Error
}

// GetByID retrieves a book by ID and user ID.
func (r *BookStore) GetByID(ctx context.Context, id uint, userID uint) (*model.Book, error) {
	var book model.Book
	err := r.db.Conn(ctx).Where("id = ? AND user_id = ?", id, userID).First(&book).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// Update updates an existing book.
func (r *BookStore) Update(ctx context.Context, book *model.Book) error {
	return r.db.Conn(ctx).Save(book).Error
}

// Delete deletes a book by ID and user ID.
func (r *BookStore) Delete(ctx context.Context, id uint, userID uint) error {
	return r.db.Conn(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Book{}).Error
}

// ListQuery defines query parameters for listing books.
type ListQuery struct {
	Page     int
	PageSize int
	Status   string
	Search   string
	UserID   uint
}

// List retrieves all books with pagination and filtering.
func (r *BookStore) List(ctx context.Context, query ListQuery) ([]model.Book, int64, error) {
	var books []model.Book
	var total int64

	db := r.db.Conn(ctx).Model(&model.Book{})

	// Filter by user ID
	db = db.Where("user_id = ?", query.UserID)

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
