// Package service provides business logic.
package service

import (
	"context"

	"github.com/mgcis-cn/ibookfs/internal/model"
	"github.com/mgcis-cn/ibookfs/internal/repository"
)

var bookRepo repository.BookRepository

// BookService handles book business logic.
type BookService struct{}

// Create creates a new book.
func (s *BookService) Create(ctx context.Context, book *model.Book) error {
	return bookRepo.Create(ctx, book)
}

// GetByID retrieves a book by ID and user ID.
func (s *BookService) GetByID(ctx context.Context, id uint, userID uint) (*model.Book, error) {
	return bookRepo.GetByID(ctx, id, userID)
}

// Update updates an existing book.
func (s *BookService) Update(ctx context.Context, book *model.Book) error {
	return bookRepo.Update(ctx, book)
}

// Delete deletes a book by ID and user ID.
func (s *BookService) Delete(ctx context.Context, id uint, userID uint) error {
	return bookRepo.Delete(ctx, id, userID)
}

// List retrieves all books with pagination and filtering.
func (s *BookService) List(ctx context.Context, query repository.ListQuery) ([]model.Book, int64, error) {
	if query.Page < 1 {
		query.Page = 1
	}
	if query.PageSize < 1 {
		query.PageSize = 10
	}
	return bookRepo.List(ctx, query)
}
