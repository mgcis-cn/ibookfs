// Package accountsecret provides account secret business logic.
package accountsecret

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"

	apierr "github.com/mgcis-cn/ibookfs/internal/apiserver/errors"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/model"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/store"
	pkgerr "github.com/mgcis-cn/ibookfs/pkg/errors"
	"golang.org/x/crypto/bcrypt"
)

// AccountSecretBiz defines the interface for account secret business logic.
type AccountSecretBiz interface {
	Create(ctx context.Context, ownerID uint, req *CreateRequest) (*CreateResponse, error)
	GetByID(ctx context.Context, id uint, ownerID uint) (*model.AccountSecret, error)
	List(ctx context.Context, ownerID uint, page, pageSize int) ([]*model.AccountSecret, int64, error)
	Delete(ctx context.Context, id uint, ownerID uint) error
	Disable(ctx context.Context, id uint, ownerID uint) error
	Enable(ctx context.Context, id uint, ownerID uint) error
	Authenticate(ctx context.Context, accountKey string, plainSecretKey string) (*model.AccountSecret, error)
	GetByAccountKey(ctx context.Context, accountKey string) (*model.AccountSecret, error)
	UpdateLastUsedAt(ctx context.Context, id uint) error
}

// accountSecretBiz is the concrete implementation of AccountSecretBiz.
type accountSecretBiz struct {
	repo store.IStore
}

// NewAccountSecretBiz creates a new account secret business logic with injected store.
func NewAccountSecretBiz(repo store.IStore) AccountSecretBiz {
	return &accountSecretBiz{
		repo: repo,
	}
}

// CreateRequest contains parameters for creating an account secret.
type CreateRequest struct {
	Name         string                   `json:"name" binding:"required"`
	Scope        model.AccountSecretScope `json:"scope" binding:"required"`
	ResourceID   *uint                    `json:"resource_id,omitempty"`
	ResourceType *string                  `json:"resource_type,omitempty"`
	Permissions  []string                 `json:"permissions" binding:"required"`
	ExpiresIn    *int64                   `json:"expires_in,omitempty"` // TTL in seconds
}

// CreateResponse contains the result of creating an account secret.
type CreateResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	AccountKey  string     `json:"account_key"`
	SecretKey   string     `json:"secret_key"` // Only shown once
	Scope       string     `json:"scope"`
	Permissions []string   `json:"permissions"`
	ExpiresAt   *time.Time `json:"expires_at,omitempty"`
}

// Create creates a new account secret.
// The secret key is only returned once and cannot be retrieved again.
func (s *accountSecretBiz) Create(ctx context.Context, ownerID uint, req *CreateRequest) (*CreateResponse, error) {
	// Validate scope
	if req.Scope == model.AccountSecretScopeResource {
		if req.ResourceID == nil || req.ResourceType == nil {
			return nil, apierr.ErrAccountResourceRequired
		}
	}

	// Validate permissions
	if err := s.validatePermissions(req.Permissions); err != nil {
		return nil, err
	}

	// Generate account key (public identifier)
	accountKey, err := generateAccountKey()
	if err != nil {
		return nil, pkgerr.Wrap(err, "生成账户密钥失败")
	}

	// Generate secret key (private)
	plainSecretKey, err := generateSecretKey()
	if err != nil {
		return nil, pkgerr.Wrap(err, "生成密钥失败")
	}

	// Hash secret key for storage
	hashedSecretKey, err := bcrypt.GenerateFromPassword([]byte(plainSecretKey), bcrypt.DefaultCost)
	if err != nil {
		return nil, pkgerr.Wrap(err, "加密密钥失败")
	}

	// Calculate expiration
	var expiresAt *time.Time
	if req.ExpiresIn != nil && *req.ExpiresIn > 0 {
		expires := time.Now().Add(time.Duration(*req.ExpiresIn) * time.Second)
		expiresAt = &expires
	}

	// Create account secret
	secret := &model.AccountSecret{
		OwnerID:      ownerID,
		Name:         req.Name,
		AccountKey:   accountKey,
		SecretKey:    string(hashedSecretKey),
		Scope:        req.Scope,
		ResourceID:   req.ResourceID,
		ResourceType: req.ResourceType,
		Permissions:  req.Permissions,
		IsEnabled:    true,
		ExpiresAt:    expiresAt,
	}

	if err := s.repo.AccountSecret().Create(ctx, secret); err != nil {
		return nil, pkgerr.Wrap(err, "创建账户密钥失败")
	}

	return &CreateResponse{
		ID:          secret.ID,
		Name:        secret.Name,
		AccountKey:  secret.AccountKey,
		SecretKey:   plainSecretKey, // Only return once
		Scope:       string(secret.Scope),
		Permissions: secret.Permissions,
		ExpiresAt:   secret.ExpiresAt,
	}, nil
}

// GetByID retrieves an account secret by ID (owner must be the owner).
func (s *accountSecretBiz) GetByID(ctx context.Context, id uint, ownerID uint) (*model.AccountSecret, error) {
	return s.repo.AccountSecret().GetByIDAndOwner(ctx, id, ownerID)
}

// List retrieves account secrets for an owner with pagination.
func (s *accountSecretBiz) List(ctx context.Context, ownerID uint, page, pageSize int) ([]*model.AccountSecret, int64, error) {
	offset := (page - 1) * pageSize
	return s.repo.AccountSecret().ListByOwner(ctx, ownerID, pageSize, offset)
}

// Delete deletes an account secret by ID.
func (s *accountSecretBiz) Delete(ctx context.Context, id uint, ownerID uint) error {
	// Verify ownership
	secret, err := s.repo.AccountSecret().GetByIDAndOwner(ctx, id, ownerID)
	if err != nil {
		return err
	}
	if secret.OwnerID != ownerID {
		return apierr.ErrAccountAccessDenied
	}
	return s.repo.AccountSecret().Delete(ctx, id)
}

// Disable disables an account secret.
func (s *accountSecretBiz) Disable(ctx context.Context, id uint, ownerID uint) error {
	secret, err := s.repo.AccountSecret().GetByIDAndOwner(ctx, id, ownerID)
	if err != nil {
		return err
	}
	secret.IsEnabled = false
	return s.repo.AccountSecret().Update(ctx, secret)
}

// Enable enables an account secret.
func (s *accountSecretBiz) Enable(ctx context.Context, id uint, ownerID uint) error {
	secret, err := s.repo.AccountSecret().GetByIDAndOwner(ctx, id, ownerID)
	if err != nil {
		return err
	}
	secret.IsEnabled = true
	return s.repo.AccountSecret().Update(ctx, secret)
}

// Authenticate validates credentials and returns the account secret if valid.
func (s *accountSecretBiz) Authenticate(ctx context.Context, accountKey string, plainSecretKey string) (*model.AccountSecret, error) {
	// Get secret by account key
	secret, err := s.repo.AccountSecret().ListValidSecretsForAuth(ctx, accountKey)
	if err != nil {
		return nil, apierr.ErrAccountInvalidCredentials
	}

	// Verify secret key
	if err := bcrypt.CompareHashAndPassword([]byte(secret.SecretKey), []byte(plainSecretKey)); err != nil {
		return nil, apierr.ErrAccountInvalidCredentials
	}

	// Update last used at
	_ = s.repo.AccountSecret().UpdateLastUsedAt(ctx, secret.ID)

	return secret, nil
}

// GetByAccountKey retrieves an account secret by account key (without auth check).
func (s *accountSecretBiz) GetByAccountKey(ctx context.Context, accountKey string) (*model.AccountSecret, error) {
	return s.repo.AccountSecret().GetByAccountKey(ctx, accountKey)
}

// UpdateLastUsedAt updates the last used timestamp for an account secret.
func (s *accountSecretBiz) UpdateLastUsedAt(ctx context.Context, id uint) error {
	return s.repo.AccountSecret().UpdateLastUsedAt(ctx, id)
}

// validatePermissions checks if the permissions are valid.
func (s *accountSecretBiz) validatePermissions(permissions []string) error {
	validPermissions := map[string]bool{
		model.PermissionRead:   true,
		model.PermissionWrite:  true,
		model.PermissionDelete: true,
	}

	if len(permissions) == 0 {
		return apierr.ErrAccountPermissionRequired
	}

	for _, p := range permissions {
		if !validPermissions[p] {
			return apierr.AccountInvalidPermission(p)
		}
	}

	return nil
}

// generateAccountKey generates a unique account key (32 characters hex).
func generateAccountKey() (string, error) {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// generateSecretKey generates a secret key (64 characters hex).
func generateSecretKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}
