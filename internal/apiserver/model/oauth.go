// Package model defines data structures.
package model

import "time"

// OAuthIdentity represents an OAuth identity linked to a user.
type OAuthIdentity struct {
	ID             uint   `json:"id" gorm:"primaryKey;comment:主键"`
	UserID         uint   `json:"user_id" gorm:"not null;index:idx_oauth_user_id;comment:用户ID"`
	Provider       string `json:"provider" gorm:"size:20;not null;index:idx_oauth_provider;comment:OAuth提供商"`
	ProviderUserID string `json:"provider_user_id" gorm:"size:255;not null;comment:第三方平台用户ID"`

	// Provider info
	ProviderUsername *string `json:"provider_username,omitempty" gorm:"size:255;comment:第三方平台用户名"`
	ProviderEmail    *string `json:"provider_email,omitempty" gorm:"size:255;comment:第三方平台邮箱"`

	// OAuth tokens (encrypted)
	AccessToken    *string    `json:"-" gorm:"type:text;comment:OAuth访问令牌"`
	RefreshToken   *string    `json:"-" gorm:"type:text;comment:OAuth刷新令牌"`
	TokenExpiresAt *time.Time `json:"token_expires_at,omitempty" gorm:"comment:令牌过期时间"`
	Scope          *string    `json:"scope,omitempty" gorm:"size:500;comment:OAuth授权范围"`

	// Profile data from provider
	ProfileData *string `json:"-" gorm:"type:longtext;comment:第三方平台用户资料（JSON）"`

	// Link status
	IsPrimary  bool       `json:"is_primary" gorm:"default:false;index:idx_oauth_is_primary;comment:是否为主账号"`
	LinkedAt   *time.Time `json:"linked_at,omitempty" gorm:"comment:绑定时间"`
	LastUsedAt *time.Time `json:"last_used_at,omitempty" gorm:"comment:最后使用时间"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:更新时间"`

	// Associations
	User *User `json:"-" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for OAuthIdentity.
func (*OAuthIdentity) TableName() string {
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

// OAuthState represents an OAuth state token for preventing CSRF and duplicate processing.
type OAuthState struct {
	ID        uint      `json:"id" gorm:"primaryKey;comment:主键"`
	State     string    `json:"state" gorm:"size:255;not null;uniqueIndex:idx_oauth_state;comment:OAuth状态令牌"`
	Provider  string    `json:"provider" gorm:"size:20;not null;comment:OAuth提供商"`
	Processed bool      `json:"processed" gorm:"default:false;comment:是否已处理"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index:idx_oauth_state_expires;comment:过期时间"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
}

// TableName returns the table name for OAuthState.
func (*OAuthState) TableName() string {
	return "oauth_states"
}
