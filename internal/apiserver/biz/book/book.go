// Package book provides book business logic.
package book

import (
	"context"

	"github.com/mgcis-cn/ibookfs/internal/apiserver/model"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/store"
)

// BookBiz defines the interface for book business logic.
type BookBiz interface {
	Create(ctx context.Context, book *model.Book) error
	GetByID(ctx context.Context, id uint, userID uint) (*model.Book, error)
	Update(ctx context.Context, book *model.Book) error
	Delete(ctx context.Context, id uint, userID uint) error
	List(ctx context.Context, query store.ListQuery) ([]model.Book, int64, error)
}

// bookBiz is the concrete implementation of BookBiz.
type bookBiz struct {
	repo store.IStore
}

// NewBookBiz creates a new book business logic with injected store.
func NewBookBiz(repo store.IStore) BookBiz {
	return &bookBiz{
		repo: repo,
	}
}

// Create creates a new book.
func (s *bookBiz) Create(ctx context.Context, book *model.Book) error {
	return s.repo.Book().Create(ctx, book)
}

// GetByID retrieves a book by ID and user ID.
func (s *bookBiz) GetByID(ctx context.Context, id uint, userID uint) (*model.Book, error) {
	return s.repo.Book().GetByID(ctx, id, userID)
}

// Update updates an existing book.
func (s *bookBiz) Update(ctx context.Context, book *model.Book) error {
	return s.repo.Book().Update(ctx, book)
}

// Delete deletes a book by ID and user ID.
func (s *bookBiz) Delete(ctx context.Context, id uint, userID uint) error {
	return s.repo.Book().Delete(ctx, id, userID)
}

// List retrieves all books with pagination and filtering.
func (s *bookBiz) List(ctx context.Context, query store.ListQuery) ([]model.Book, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}
	return s.repo.Book().List(ctx, query)
}
