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
	UpdateFields(ctx context.Context, bookID uint, userID uint, fields map[string]any) error
	Delete(ctx context.Context, id uint, userID uint) error
	List(ctx context.Context, query ListQuery) ([]model.Book, int64, error)
	SetCoverIfEmpty(ctx context.Context, bookID uint, userID uint, coverURL string) error
	UpdateUploadedPages(ctx context.Context, bookID uint) error
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

// UpdateFields updates specific fields of a book.
func (r *BookStore) UpdateFields(ctx context.Context, bookID uint, userID uint, fields map[string]any) error {
	return r.db.Conn(ctx).Model(&model.Book{}).
		Where("id = ? AND user_id = ?", bookID, userID).
		Updates(fields).Error
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

// SetCoverIfEmpty sets the book cover if it's currently empty.
func (r *BookStore) SetCoverIfEmpty(ctx context.Context, bookID uint, userID uint, coverURL string) error {
	return r.db.Conn(ctx).Model(&model.Book{}).
		Where("id = ? AND user_id = ? AND (cover = '' OR cover IS NULL)", bookID, userID).
		Update("cover", coverURL).Error
}

// UpdateUploadedPages recalculates and updates the uploaded_pages count for a book.
func (r *BookStore) UpdateUploadedPages(ctx context.Context, bookID uint) error {
	var count int64
	err := r.db.Conn(ctx).Model(&model.BookImage{}).
		Where("book_id = ?", bookID).
		Count(&count).Error
	if err != nil {
		return err
	}

	return r.db.Conn(ctx).Model(&model.Book{}).
		Where("id = ?", bookID).
		Update("uploaded_pages", count).Error
}
