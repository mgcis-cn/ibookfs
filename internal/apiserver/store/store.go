// Package store provides centralized data access layer.
package store

import (
	"github.com/mgcis-cn/ibookfs/pkg/database"
)

// IStore defines the interface for accessing all repositories.
type IStore interface {
	// Book returns the book store
	Book() IBookStore

	// Image returns the image store
	Image() IImageStore

	// AccountSecret returns the account secret store
	AccountSecret() IAccountSecretStore

	// User returns the user store
	User() IUserStore
}

// Store implements IStore interface with lazy initialization.
type Store struct {
	db database.Database
}

// New creates a new Store.
func New(db database.Database) *Store {
	return &Store{db: db}
}

// Book returns the book store (lazy initialization).
func (s *Store) Book() IBookStore {
	return NewBookStore(s.db)
}

// Image returns the image store (lazy initialization).
func (s *Store) Image() IImageStore {
	return NewImageStore(s.db)
}

// AccountSecret returns the account secret store (lazy initialization).
func (s *Store) AccountSecret() IAccountSecretStore {
	return NewAccountSecretStore(s.db)
}

// User returns the user store (lazy initialization).
func (s *Store) User() IUserStore {
	return NewUserStore(s.db)
}
