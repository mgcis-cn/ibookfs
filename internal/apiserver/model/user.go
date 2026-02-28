// Package model defines data structures.
package model

import "time"

// User represents a user entity.
type User struct {
	ID    uint   `json:"id" gorm:"primaryKey;comment:主键"`
	Email string `json:"email" gorm:"size:255;not null;uniqueIndex:idx_users_email;comment:邮箱"`

	// Profile
	FirstName   string  `json:"first_name" gorm:"size:100;not null;comment:名"`
	LastName    string  `json:"last_name" gorm:"size:100;not null;comment:姓"`
	DisplayName *string `json:"display_name,omitempty" gorm:"size:255;comment:显示名称"`
	AvatarURL   *string `json:"avatar_url,omitempty" gorm:"size:500;comment:头像URL"`

	// Status
	Status UserStatus `json:"status" gorm:"size:20;default:'active';comment:状态"`
	Role   UserRole   `json:"role" gorm:"size:20;default:'user';comment:角色"`

	// Settings
	Preferences *string `json:"preferences,omitempty" gorm:"type:longtext;comment:用户偏好设置（JSON）"`
	Language    string  `json:"language" gorm:"size:10;default:'zh-CN';comment:语言设置"`
	Timezone    string  `json:"timezone" gorm:"size:50;default:'Asia/Shanghai';comment:时区"`

	// Tracking
	LastLoginAt *time.Time `json:"last_login_at,omitempty" gorm:"comment:最后登录时间"`
	LastLoginIP *string    `json:"last_login_ip,omitempty" gorm:"size:45;comment:最后登录IP"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt time.Time `json:"updated_at" gorm:"comment:更新时间"`

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
func (*User) TableName() string {
	return "users"
}

// EmailVerificationCode represents a verification code for email operations.
type EmailVerificationCode struct {
	ID        uint      `json:"id" gorm:"primaryKey;comment:主键"`
	Email     string    `json:"email" gorm:"size:255;not null;comment:邮箱地址"`
	Code      string    `json:"code" gorm:"size:10;not null;comment:验证码（6位数字）"`
	Type      CodeType  `json:"type" gorm:"size:20;not null;comment:类型"`
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;comment:过期时间"`
	Used      bool      `json:"used" gorm:"default:false;comment:是否已使用"`
	Attempts  int       `json:"attempts" gorm:"default:0;comment:验证尝试次数"`
	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`
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
	ID               uint    `json:"id" gorm:"primaryKey;comment:主键"`
	UserID           uint    `json:"user_id" gorm:"not null;index:idx_sessions_user_id;comment:用户ID"`
	TokenHash        string  `json:"-" gorm:"size:64;not null;uniqueIndex:idx_sessions_token_hash;comment:访问令牌哈希"`
	RefreshTokenHash *string `json:"-" gorm:"size:64;comment:刷新令牌哈希"`

	// Device info
	UserAgent  *string `json:"user_agent,omitempty" gorm:"size:500;comment:用户代理（浏览器信息）"`
	IPAddress  *string `json:"ip_address,omitempty" gorm:"size:45;comment:IP地址"`
	DeviceType *string `json:"device_type,omitempty" gorm:"size:50;comment:设备类型"`
	DeviceName *string `json:"device_name,omitempty" gorm:"size:100;comment:设备名称"`

	// Status
	ExpiresAt      time.Time  `json:"expires_at" gorm:"not null;index:idx_sessions_expires;comment:过期时间"`
	LastActivityAt *time.Time `json:"last_activity_at,omitempty" gorm:"comment:最后活跃时间"`
	Revoked        bool       `json:"revoked" gorm:"default:false;index:idx_sessions_revoked;comment:是否已撤销"`
	RevokedAt      *time.Time `json:"revoked_at,omitempty" gorm:"comment:撤销时间"`
	RevokedReason  *string    `json:"revoked_reason,omitempty" gorm:"size:255;comment:撤销原因"`

	CreatedAt time.Time `json:"created_at" gorm:"comment:创建时间"`

	// Associations
	User *User `json:"-" gorm:"foreignKey:UserID"`
}

// TableName returns the table name for Session.
func (*Session) TableName() string {
	return "sessions"
}

// LoginHistory represents a login attempt record.
type LoginHistory struct {
	ID            uint        `json:"id" gorm:"primaryKey;comment:主键"`
	UserID        *uint       `json:"user_id,omitempty" gorm:"index:idx_login_user_id;comment:用户ID"`
	LoginMethod   LoginMethod `json:"login_method" gorm:"size:20;not null;comment:登录方式"`
	Success       bool        `json:"success" gorm:"index:idx_login_success;comment:是否成功"`
	FailureReason *string     `json:"failure_reason,omitempty" gorm:"size:255;comment:失败原因"`

	// Request info
	IPAddress         *string `json:"ip_address,omitempty" gorm:"size:45;comment:IP地址"`
	UserAgent         *string `json:"user_agent,omitempty" gorm:"size:500;comment:用户代理"`
	DeviceFingerprint *string `json:"device_fingerprint,omitempty" gorm:"size:255;comment:设备指纹"`

	// Geo location
	Country *string `json:"country,omitempty" gorm:"size:100;comment:国家"`
	Region  *string `json:"region,omitempty" gorm:"size:100;comment:地区/省份"`
	City    *string `json:"city,omitempty" gorm:"size:100;comment:城市"`

	CreatedAt time.Time `json:"created_at" gorm:"index:idx_login_created_at;comment:创建时间"`

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
func (*LoginHistory) TableName() string {
	return "login_history"
}
