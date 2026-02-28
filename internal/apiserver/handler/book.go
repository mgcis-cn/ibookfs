// Package handler provides HTTP handlers for book management.
package handler

import (
	"context"
	"errors"

	"github.com/mgcis-cn/ibookfs/internal/apiserver/model"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/store"
	v1 "github.com/mgcis-cn/ibookfs/pkg/api/apiserver/v1"
	contextx "github.com/mgcis-cn/ibookfs/pkg/context"
)

type BookRouter interface {
	ListBooks(ctx context.Context, req *v1.ListBooksRequest) (*v1.ListBooksResponse, error)
	GetBook(ctx context.Context, req *v1.GetBookRequest) (*v1.GetBookResponse, error)
	CreateBook(ctx context.Context, req *v1.CreateBookRequest) (*v1.CreateBookResponse, error)
	UpdateBook(ctx context.Context, req *v1.UpdateBookRequest) (*v1.UpdateBookResponse, error)
	DeleteBook(ctx context.Context, req *v1.DeleteBookRequest) (*v1.DeleteBookResponse, error)
}

// ListBooks returns all books for the current user.
func (h *handler) ListBooks(ctx context.Context, req *v1.ListBooksRequest) (*v1.ListBooksResponse, error) {
	userID := contextx.UserId(ctx)
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}

	query := store.ListQuery{
		Page:     req.Page,
		PageSize: req.PageSize,
		Status:   req.Status,
		Search:   req.Search,
		UserID:   uint(userID),
	}

	books, total, err := h.biz.Book().List(ctx, query)
	if err != nil {
		return &v1.ListBooksResponse{}, err
	}

	return &v1.ListBooksResponse{
		Data:    books,
		Total:   total,
		Message: "success",
	}, nil
}

// GetBook returns a single book by ID.
func (h *handler) GetBook(ctx context.Context, req *v1.GetBookRequest) (*v1.GetBookResponse, error) {
	userID := contextx.UserId(ctx)
	book, err := h.biz.Book().GetByID(ctx, uint(req.Id), uint(userID))
	if err != nil {
		return &v1.GetBookResponse{}, errors.New("book not found")
	}

	return &v1.GetBookResponse{
		Data:    book,
		Message: "success",
	}, nil
}

// CreateBook creates a new book.
func (h *handler) CreateBook(ctx context.Context, req *v1.CreateBookRequest) (*v1.CreateBookResponse, error) {
	userID := uint(contextx.UserId(ctx))

	book := &model.Book{
		Title:     req.Title,
		Author:    req.Author,
		ISBN:      req.ISBN,
		Publisher: req.Publisher,
		Year:      req.Year,
		UserID:    userID,
	}

	if err := h.biz.Book().Create(ctx, book); err != nil {
		return &v1.CreateBookResponse{}, err
	}

	return &v1.CreateBookResponse{
		Data:    book,
		Message: "success",
	}, nil
}

// UpdateBook updates an existing book.
func (h *handler) UpdateBook(ctx context.Context, req *v1.UpdateBookRequest) (*v1.UpdateBookResponse, error) {
	userID := uint(contextx.UserId(ctx))

	data := &model.Book{
		Title:     req.Title,
		Author:    req.Author,
		ISBN:      req.ISBN,
		Publisher: req.Publisher,
		Year:      req.Year,
	}
	data.ID = uint(req.Id)
	data.UserID = userID

	if err := h.biz.Book().Update(ctx, data); err != nil {
		return &v1.UpdateBookResponse{}, err
	}

	return &v1.UpdateBookResponse{
		Data:    data,
		Message: "success",
	}, nil
}

// DeleteBook deletes a book by ID.
func (h *handler) DeleteBook(ctx context.Context, req *v1.DeleteBookRequest) (*v1.DeleteBookResponse, error) {
	if err := h.biz.Book().Delete(ctx, uint(req.Id), uint(contextx.UserId(ctx))); err != nil {
		return &v1.DeleteBookResponse{}, err
	}

	return &v1.DeleteBookResponse{
		Message: "success",
	}, nil
}
