// Package model defines data structures.
package model

import "time"

// User represents a user entity.
type User struct {
	ID              uint       `json:"id" gorm:"primaryKey"`
	Email           string     `json:"email" gorm:"size:255;not null;uniqueIndex:idx_users_email"`
	EmailVerified   bool       `json:"email_verified" gorm:"default:false"`
	EmailVerifiedAt *time.Time `json:"email_verified_at,omitempty"`
	PasswordHash    string     `json:"-" gorm:"size:255"`

	// Profile
	FirstName   string  `json:"first_name" gorm:"size:100;not null"`
	LastName    string  `json:"last_name" gorm:"size:100;not null"`
	DisplayName *string `json:"display_name,omitempty" gorm:"size:255"`
	AvatarURL   *string `json:"avatar_url,omitempty" gorm:"size:500"`
	Bio         *string `json:"bio,omitempty" gorm:"type:text"`

	// Status
	Status UserStatus `json:"status" gorm:"size:20;default:'active'"`
	Role   UserRole   `json:"role" gorm:"size:20;default:'user'"`

	// Settings
	Preferences *string `json:"preferences,omitempty" gorm:"type:json"`
	Language    string  `json:"language" gorm:"size:10;default:'zh-CN'"`
	Timezone    string  `json:"timezone" gorm:"size:50;default:'Asia/Shanghai'"`

	// Tracking
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
	LastLoginIP *string    `json:"last_login_ip,omitempty" gorm:"size:45"`

	// Timestamps
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Associations
	OAuthIdentities []OAuthIdentity `json:"linked_accounts,omitempty" gorm:"foreignKey:UserID"`
	Sessions        []Session       `json:"-" gorm:"foreignKey:UserID"`
}

// UserStatus represents the status of a user.
type UserStatus string

const (
	UserStatusActive    UserStatus = "active"
	UserStatusSuspended UserStatus = "suspended"
	UserStatusDeleted   UserStatus = "deleted"
)

// UserRole represents the role of a user.
type UserRole string

const (
	UserRoleUser       UserRole = "user"
	UserRoleAdmin      UserRole = "admin"
	UserRoleSuperAdmin UserRole = "super_admin"
)

// TableName returns the table name for User.
func (User) TableName() string {
	return "users"
}

// EmailVerificationCode represents a verification code for email operations.
type EmailVerificationCode struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"size:255;not null"`
	Code      string    `json:"code" gorm:"size:10;not null"`
	Type      CodeType  `json:"type" gorm:"size:20;not null"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null"`
	Used      bool      `json:"used" gorm:"default:false"`
	Attempts  int       `json:"attempts" gorm:"default:0"`
	CreatedAt time.Time `json:"created_at"`
}

// CodeType represents the type of verification code.
type CodeType string

const (
	CodeTypeRegister      CodeType = "register"
	CodeTypeLogin         CodeType = "login"
	CodeTypeResetPassword CodeType = "reset_password"
	CodeTypeBindEmail     CodeType = "bind_email"
)

// TableName returns the table name for EmailVerificationCode.
func (EmailVerificationCode) TableName() string {
	return "email_verification_codes"
}

// Session represents a user session.
type Session struct {
	ID               uint    `json:"id" gorm:"primaryKey"`
	UserID           uint    `json:"user_id" gorm:"not null;index:idx_sessions_user_id"`
	TokenHash        string  `json:"-" gorm:"size:64;not null;uniqueIndex:idx_sessions_token_hash"`
	RefreshTokenHash *string `json:"-" gorm:"size:64"`

	// Device info
	UserAgent  *string `json:"user_agent,omitempty" gorm:"size:500"`
	IPAddress  *string `json:"ip_address,omitempty" gorm:"size:45"`
	DeviceType *string `json:"device_type,omitempty" gorm:"size:50"`
	DeviceName *string `json:"device_name,omitempty" gorm:"size:100"`

	// Status
	ExpiresAt      time.Time  `json:"expires_at" gorm:"not null;index:idx_sessions_expires"`
	LastActivityAt *time.Time `json:"last_activity_at,omitempty"`
	Revoked        bool       `json:"revoked" gorm:"default:false;index:idx_sessions_revoked"`
	RevokedAt      *time.Time `json:"revoked_at,omitempty"`
	RevokedReason  *string    `json:"revoked_reason,omitempty" gorm:"size:255"`

	CreatedAt time.Time `json:"created_at"`

	// Associations
	User *User `json:"-" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for Session.
func (Session) TableName() string {
	return "sessions"
}

// LoginHistory represents a login attempt record.
type LoginHistory struct {
	ID            uint        `json:"id" gorm:"primaryKey"`
	UserID        *uint       `json:"user_id,omitempty" gorm:"index:idx_login_user_id"`
	LoginMethod   LoginMethod `json:"login_method" gorm:"size:20;not null"`
	Success       bool        `json:"success" gorm:"index:idx_login_success"`
	FailureReason *string     `json:"failure_reason,omitempty" gorm:"size:255"`

	// Request info
	IPAddress         *string `json:"ip_address,omitempty" gorm:"size:45"`
	UserAgent         *string `json:"user_agent,omitempty" gorm:"size:500"`
	DeviceFingerprint *string `json:"device_fingerprint,omitempty" gorm:"size:255"`

	// Geo location
	Country *string `json:"country,omitempty" gorm:"size:100"`
	Region  *string `json:"region,omitempty" gorm:"size:100"`
	City    *string `json:"city,omitempty" gorm:"size:100"`

	CreatedAt time.Time `json:"created_at" gorm:"index:idx_login_created_at"`

	// Associations
	User *User `json:"-" gorm:"foreignKey:UserID"`
}

// LoginMethod represents the login method used.
type LoginMethod string

const (
	LoginMethodEmail  LoginMethod = "email"
	LoginMethodGitHub LoginMethod = "github"
	LoginMethodGitee  LoginMethod = "gitee"
	LoginMethodICloud LoginMethod = "icloud"
	LoginMethodGoogle LoginMethod = "google"
	LoginMethodWechat LoginMethod = "wechat"
)

// TableName returns the table name for LoginHistory.
func (LoginHistory) TableName() string {
	return "login_history"
}

// AccountSecurityLog represents a security event log.
type AccountSecurityLog struct {
	ID          uint              `json:"id" gorm:"primaryKey"`
	UserID      uint              `json:"user_id" gorm:"not null;index:idx_security_user_id"`
	EventType   SecurityEventType `json:"event_type" gorm:"size:50;not null;index:idx_security_event_type"`
	Description *string           `json:"description,omitempty" gorm:"type:text"`
	IPAddress   *string           `json:"ip_address,omitempty" gorm:"size:45"`
	UserAgent   *string           `json:"user_agent,omitempty" gorm:"size:500"`
	Metadata    *string           `json:"metadata,omitempty" gorm:"type:json"`
	CreatedAt   time.Time         `json:"created_at" gorm:"index:idx_security_created_at"`

	// Associations
	User *User `json:"-" gorm:"foreignKey:UserID"`
}

// SecurityEventType represents the type of security event.
type SecurityEventType string

const (
	SecurityEventPasswordChanged        SecurityEventType = "password_changed"
	SecurityEventEmailBound             SecurityEventType = "email_bound"
	SecurityEventEmailUnbound           SecurityEventType = "email_unbound"
	SecurityEventOAuthLinked            SecurityEventType = "oauth_linked"
	SecurityEventOAuthUnlinked          SecurityEventType = "oauth_unlinked"
	SecurityEventAccountLocked          SecurityEventType = "account_locked"
	SecurityEventAccountUnlocked        SecurityEventType = "account_unlocked"
	SecurityEventPasswordResetRequested SecurityEventType = "password_reset_requested"
	SecurityEventSuspiciousActivity     SecurityEventType = "suspicious_activity"
	SecurityEventSessionRevoked         SecurityEventType = "session_revoked"
)

// TableName returns the table name for AccountSecurityLog.
func (AccountSecurityLog) TableName() string {
	return "account_security_log"
}
