// Package repository provides account secret data access layer.
package store

import (
	"context"
	"time"

	"github.com/mgcis-cn/ibookfs/internal/apiserver/model"
	"github.com/mgcis-cn/ibookfs/pkg/database"
)

// IAccountSecretStore defines the interface for account secret data operations.
type IAccountSecretStore interface {
	Create(ctx context.Context, secret *model.AccountSecret) error
	GetByID(ctx context.Context, id uint) (*model.AccountSecret, error)
	GetByAccountKey(ctx context.Context, accountKey string) (*model.AccountSecret, error)
	GetByIDAndOwner(ctx context.Context, id uint, ownerID uint) (*model.AccountSecret, error)
	ListByOwner(ctx context.Context, ownerID uint, limit int, offset int) ([]*model.AccountSecret, int64, error)
	Update(ctx context.Context, secret *model.AccountSecret) error
	Delete(ctx context.Context, id uint) error
	UpdateLastUsedAt(ctx context.Context, id uint) error
	ListValidSecretsForAuth(ctx context.Context, accountKey string) (*model.AccountSecret, error)
}

// AccountSecretStore handles account secret data operations.
type AccountSecretStore struct {
	db database.Database
}

// NewAccountSecretStore creates a new account secret store with injected DatabaseOptions.
func NewAccountSecretStore(db database.Database) *AccountSecretStore {
	return &AccountSecretStore{db: db}
}

// Create creates a new account secret record.
func (r *AccountSecretStore) Create(ctx context.Context, secret *model.AccountSecret) error {
	return r.db.Conn(ctx).Create(secret).Error
}

// GetByID retrieves an account secret by ID.
func (r *AccountSecretStore) GetByID(ctx context.Context, id uint) (*model.AccountSecret, error) {
	var secret model.AccountSecret
	err := r.db.Conn(ctx).Where("id = ?", id).First(&secret).Error
	if err != nil {
		return nil, err
	}
	return &secret, nil
}

// GetByAccountKey retrieves an account secret by its account key.
func (r *AccountSecretStore) GetByAccountKey(ctx context.Context, accountKey string) (*model.AccountSecret, error) {
	var secret model.AccountSecret
	err := r.db.Conn(ctx).Where("account_key = ?", accountKey).First(&secret).Error
	if err != nil {
		return nil, err
	}
	return &secret, nil
}

// GetByIDAndOwner retrieves an account secret by ID and owner ID.
func (r *AccountSecretStore) GetByIDAndOwner(ctx context.Context, id uint, ownerID uint) (*model.AccountSecret, error) {
	var secret model.AccountSecret
	err := r.db.Conn(ctx).
		Where("id = ? AND owner_id = ?", id, ownerID).
		First(&secret).Error
	if err != nil {
		return nil, err
	}
	return &secret, nil
}

// ListByOwner retrieves all account secrets for an owner with pagination.
func (r *AccountSecretStore) ListByOwner(ctx context.Context, ownerID uint, limit, offset int) ([]*model.AccountSecret, int64, error) {
	var secrets []*model.AccountSecret
	var total int64

	db := r.db.Conn(ctx).Model(&model.AccountSecret{}).Where("owner_id = ?", ownerID)

	// Count total
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Fetch with pagination
	err := db.Limit(limit).Offset(offset).Order("created_at DESC").Find(&secrets).Error
	return secrets, total, err
}

// Update updates an account secret record.
func (r *AccountSecretStore) Update(ctx context.Context, secret *model.AccountSecret) error {
	return r.db.Conn(ctx).Save(secret).Error
}

// Delete deletes an account secret by ID.
func (r *AccountSecretStore) Delete(ctx context.Context, id uint) error {
	return r.db.Conn(ctx).Delete(&model.AccountSecret{}, id).Error
}

// UpdateLastUsedAt updates the last_used_at timestamp for a secret.
func (r *AccountSecretStore) UpdateLastUsedAt(ctx context.Context, id uint) error {
	now := time.Now()
	return r.db.Conn(ctx).
		Model(&model.AccountSecret{}).
		Where("id = ?", id).
		Update("last_used_at", now).Error
}

// ListValidSecretsForAuth retrieves valid secrets for authentication by account key.
// This excludes expired and disabled secrets.
func (r *AccountSecretStore) ListValidSecretsForAuth(ctx context.Context, accountKey string) (*model.AccountSecret, error) {
	var secret model.AccountSecret
	err := r.db.Conn(ctx).
		Where("account_key = ? AND is_enabled = ?", accountKey, true).
		Where("(expires_at IS NULL OR expires_at > ?)", time.Now()).
		First(&secret).Error
	if err != nil {
		return nil, err
	}
	return &secret, nil
}
