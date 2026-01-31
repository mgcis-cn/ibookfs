// Package model defines data structures.
package model

import "time"

// OAuthIdentity represents an OAuth identity linked to a user.
type OAuthIdentity struct {
	ID             uint   `json:"id" gorm:"primaryKey"`
	UserID         uint   `json:"user_id" gorm:"not null;index:idx_oauth_user_id"`
	Provider       string `json:"provider" gorm:"size:20;not null;index:idx_oauth_provider"`
	ProviderUserID string `json:"provider_user_id" gorm:"size:255;not null"`

	// Provider info
	ProviderUsername *string `json:"provider_username,omitempty" gorm:"size:255"`
	ProviderEmail    *string `json:"provider_email,omitempty" gorm:"size:255"`

	// OAuth tokens (encrypted)
	AccessToken    *string    `json:"-" gorm:"type:text"`
	RefreshToken   *string    `json:"-" gorm:"type:text"`
	TokenExpiresAt *time.Time `json:"token_expires_at,omitempty"`
	Scope          *string    `json:"scope,omitempty" gorm:"size:500"`

	// Profile data from provider
	ProfileData *string `json:"-" gorm:"type:json"`

	// Link status
	IsPrimary  bool       `json:"is_primary" gorm:"default:false;index:idx_oauth_is_primary"`
	LinkedAt   *time.Time `json:"linked_at,omitempty"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Associations
	User *User `json:"-" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for OAuthIdentity.
func (OAuthIdentity) TableName() string {
	return "oauth_identities"
}

// OAuthProvider represents supported OAuth providers.
type OAuthProvider string

const (
	OAuthProviderGitHub OAuthProvider = "github"
	OAuthProviderGitee  OAuthProvider = "gitee"
	OAuthProviderICloud OAuthProvider = "icloud"
	OAuthProviderGoogle OAuthProvider = "google"
	OAuthProviderWechat OAuthProvider = "wechat"
)

// IsValid checks if the provider is valid.
func (p OAuthProvider) IsValid() bool {
	switch p {
	case OAuthProviderGitHub, OAuthProviderGitee, OAuthProviderICloud,
		OAuthProviderGoogle, OAuthProviderWechat:
		return true
	default:
		return false
	}
}

// String returns the string representation of the provider.
func (p OAuthProvider) String() string {
	return string(p)
}

// ToLoginMethod converts OAuthProvider to LoginMethod.
func (p OAuthProvider) ToLoginMethod() LoginMethod {
	return LoginMethod(p)
}

// OAuthConfig represents OAuth configuration for a provider.
type OAuthConfig struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	RedirectURL  string   `json:"redirect_url"`
	AuthURL      string   `json:"auth_url"`
	TokenURL     string   `json:"token_url"`
	UserInfoURL  string   `json:"user_info_url"`
	Scopes       []string `json:"scopes"`
}

// OAuthUserInfo represents user info from OAuth provider.
type OAuthUserInfo struct {
	ID        string        `json:"id"`
	Username  string        `json:"login"`
	Email     string        `json:"email"`
	Name      string        `json:"name"`
	AvatarURL string        `json:"avatar_url"`
	Provider  OAuthProvider `json:"provider"`
}

// OAuthTokenResponse represents OAuth token response.
type OAuthTokenResponse struct {
	AccessToken  string    `json:"access_token"`
	TokenType    string    `json:"token_type"`
	ExpiresIn    int       `json:"expires_in"`
	ExpiresAt    time.Time `json:"-"` // Calculated from ExpiresIn
	RefreshToken string    `json:"refresh_token"`
	Scope        string    `json:"scope"`
}
