// Package model defines account secret data structures for API authentication.
package model

import "time"

// AccountSecret represents a secret key for API authentication.
// Similar to cloud providers' AccessKey system (e.g., Alibaba Cloud AK/SK).
type AccountSecret struct {
	ID         uint   `json:"id" gorm:"primaryKey;comment:主键"`
	OwnerID    uint   `json:"owner_id" gorm:"not null;index:idx_account_secrets_owner_id;comment:所有者ID"`
	Name       string `json:"name" gorm:"size:100;not null;comment:密钥名称"`
	AccountKey string `json:"account_key" gorm:"size:32;uniqueIndex;not null;comment:账户标识(公开)"`
	SecretKey  string `json:"-" gorm:"size:64;not null;comment:私密密钥(保密,哈希存储)"`

	// Permission scope
	Scope        AccountSecretScope `json:"scope" gorm:"size:20;not null;index:idx_account_secrets_scope;comment:权限范围"`
	ResourceID   *uint              `json:"resource_id,omitempty" gorm:"index:idx_account_secrets_resource_id;comment:资源ID(resource级别时)"`
	ResourceType *string            `json:"resource_type,omitempty" gorm:"size:50;comment:资源类型"`

	// Permissions (e.g., read, write, delete)
	Permissions []string `json:"permissions" gorm:"type:json;comment:权限列表"`

	// Status
	IsEnabled bool       `json:"is_enabled" gorm:"default:true;index:idx_account_secrets_enabled;comment:是否启用"`
	ExpiresAt *time.Time `json:"expires_at,omitempty" gorm:"index:idx_account_secrets_expires;comment:过期时间"`

	// Timestamps
	CreatedAt  time.Time  `json:"created_at" gorm:"comment:创建时间"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty" gorm:"comment:最后使用时间"`
}

// TableName returns the table name for AccountSecret.
func (*AccountSecret) TableName() string {
	return "account_secrets"
}

// AccountSecretScope represents the permission scope of an account secret.
type AccountSecretScope string

const (
	// AccountSecretScopeUser indicates the secret can access all resources of the user.
	AccountSecretScopeUser AccountSecretScope = "user"

	// AccountSecretScopeResource indicates the secret can only access a specific resource.
	AccountSecretScopeResource AccountSecretScope = "resource"
)

// AccountSecretPermission represents available permission types.
const (
	PermissionRead   = "read"
	PermissionWrite  = "write"
	PermissionDelete = "delete"
)

// IsValid checks if the secret key is valid (enabled and not expired).
func (a *AccountSecret) IsValid() bool {
	if !a.IsEnabled {
		return false
	}
	if a.ExpiresAt != nil && a.ExpiresAt.Before(time.Now()) {
		return false
	}
	return true
}

// HasPermission checks if the secret has the specified permission.
func (a *AccountSecret) HasPermission(permission string) bool {
	for _, p := range a.Permissions {
		if p == permission {
			return true
		}
	}
	return false
}

// IsUserLevel returns true if this is a user-level secret.
func (a *AccountSecret) IsUserLevel() bool {
	return a.Scope == AccountSecretScopeUser
}

// IsResourceLevel returns true if this is a resource-level secret.
func (a *AccountSecret) IsResourceLevel() bool {
	return a.Scope == AccountSecretScopeResource
}
