// Package service provides business logic.
package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/mgcis-cn/ibookfs/internal/database"
	"github.com/mgcis-cn/ibookfs/internal/model"
	"github.com/mgcis-cn/ibookfs/internal/pkg/util"
)

// AuthService handles authentication operations.
type AuthService struct {
	jwtManager   *util.JWTManager
	emailService *EmailService
}

// NewAuthService creates a new auth service.
func NewAuthService(jwtManager *util.JWTManager, emailService *EmailService) *AuthService {
	return &AuthService{
		jwtManager:   jwtManager,
		emailService: emailService,
	}
}

// SendCodeRequest represents request to send verification code.
type SendCodeRequest struct {
	Email string `json:"email" binding:"required,email"`
	Type  string `json:"type" binding:"required"`
}

// LoginRequest represents email login request.
type LoginRequest struct {
	Email string `json:"email" binding:"required,email"`
	Code  string `json:"code" binding:"required,len=6"`
}

// RegisterRequest represents registration request.
type RegisterRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Code      string `json:"code" binding:"required,len=6"`
}

// AuthResponse represents authentication response.
type AuthResponse struct {
	User           *model.User           `json:"user"`
	Token          string                `json:"token"`
	RefreshToken   string                `json:"refresh_token,omitempty"`
	ExpiresIn      int64                 `json:"expires_in"`
	LinkedAccounts []model.OAuthIdentity `json:"linked_accounts,omitempty"`
	IsNewUser      bool                  `json:"is_new_user,omitempty"`
}

// SendCode sends a verification code to email.
func (s *AuthService) SendCode(req SendCodeRequest) error {
	db := database.Default()

	// Check rate limiting
	var recentCode int64
	db.Model(&model.EmailVerificationCode{}).
		Where("email = ? AND type = ? AND created_at > ?", req.Email, req.Type, time.Now().Add(-60*time.Second)).
		Count(&recentCode)

	if recentCode > 0 {
		return errors.New("请等待60秒后重新发送")
	}

	// Generate verification code based on email configuration
	var code string
	isEmailConfigured := s.emailService.cfg.IsEmailEnabled()

	if !isEmailConfigured {
		// Email not properly configured - use demo code
		code = "123456"
		fmt.Printf("[DEMO MODE] Email not configured, using fixed verification code: %s\n", code)
	} else {
		// Email is configured - generate random code
		code = util.GenerateVerificationCode()
	}

	// Store code
	verificationCode := &model.EmailVerificationCode{
		Email:     req.Email,
		Code:      code,
		Type:      model.CodeType(req.Type),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := db.Create(verificationCode).Error; err != nil {
		return fmt.Errorf("failed to store verification code: %w", err)
	}

	// Send email with verification code
	if isEmailConfigured {
		// Email is configured - send real email and fail on error
		if err := s.emailService.SendVerificationCode(req.Email, code, req.Type); err != nil {
			// Delete the verification code since email failed
			db.Delete(verificationCode)
			return fmt.Errorf("failed to send email: %w", err)
		}
	} else {
		// Email not configured - log only
		fmt.Printf("[DEMO MODE] Verification code for %s: %s (email sending skipped - not configured)\n", req.Email, code)
	}

	return nil
}

// Login authenticates user with email and verification code.
func (s *AuthService) Login(req LoginRequest, ipAddress, userAgent string) (*AuthResponse, error) {
	db := database.Default()

	// Find valid verification code
	var code model.EmailVerificationCode
	err := db.Where("email = ? AND type = ? AND used = false AND expires_at > ?",
		req.Email, model.CodeTypeLogin, time.Now()).
		Order("created_at DESC").
		First(&code).Error

	if err != nil {
		return nil, errors.New("验证码无效或已过期")
	}

	// Verify code
	if code.Code != req.Code {
		code.Attempts++
		db.Save(&code)
		return nil, errors.New("验证码错误")
	}

	// Mark code as used
	code.Used = true
	db.Save(&code)

	// Find user
	var user model.User
	err = db.Where("email = ?", req.Email).First(&user).Error
	if err != nil {
		return nil, errors.New("用户不存在，请先注册")
	}

	// Check user status
	if user.Status != model.UserStatusActive {
		return nil, errors.New("账号已被禁用")
	}

	// Generate tokens
	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	user.LastLoginIP = &ipAddress
	db.Save(&user)

	// Log login history
	s.logLoginHistory(&user, model.LoginMethodEmail, true, "", ipAddress, userAgent)

	// Load OAuth identities
	var oauthIdentities []model.OAuthIdentity
	db.Where("user_id = ?", user.ID).Find(&oauthIdentities)

	// Create session
	s.createSession(&user, tokenPair.AccessToken, ipAddress, userAgent)

	return &AuthResponse{
		User:           &user,
		Token:          tokenPair.AccessToken,
		RefreshToken:   tokenPair.RefreshToken,
		ExpiresIn:      int64(time.Until(tokenPair.ExpiresAt).Seconds()),
		LinkedAccounts: oauthIdentities,
	}, nil
}

// Register creates a new user account.
func (s *AuthService) Register(req RegisterRequest, ipAddress, userAgent string) (*AuthResponse, error) {
	db := database.Default()

	// Check if email already exists
	var existingUser model.User
	err := db.Where("email = ?", req.Email).First(&existingUser).Error
	if err == nil {
		return nil, errors.New("该邮箱已被注册")
	}

	// Verify code
	var code model.EmailVerificationCode
	err = db.Where("email = ? AND type = ? AND used = false AND expires_at > ?",
		req.Email, model.CodeTypeRegister, time.Now()).
		Order("created_at DESC").
		First(&code).Error

	if err != nil {
		return nil, errors.New("验证码无效或已过期")
	}

	if code.Code != req.Code {
		code.Attempts++
		db.Save(&code)
		return nil, errors.New("验证码错误")
	}

	// Mark code as used
	code.Used = true
	db.Save(&code)

	// Create user
	displayName := req.FirstName + " " + req.LastName
	user := &model.User{
		Email:           req.Email,
		EmailVerified:   true,
		EmailVerifiedAt: util.Ptr(time.Now()),
		FirstName:       req.FirstName,
		LastName:        req.LastName,
		DisplayName:     &displayName,
		Status:          model.UserStatusActive,
		Role:            model.UserRoleUser,
		Language:        "zh-CN",
		Timezone:        "Asia/Shanghai",
	}

	if err := db.Create(user).Error; err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Log login history
	s.logLoginHistory(user, model.LoginMethodEmail, true, "", ipAddress, userAgent)

	// Create session
	s.createSession(user, tokenPair.AccessToken, ipAddress, userAgent)

	return &AuthResponse{
		User:         user,
		Token:        tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresIn:    int64(time.Until(tokenPair.ExpiresAt).Seconds()),
		IsNewUser:    true,
	}, nil
}

// RefreshToken refreshes an access token using refresh token.
func (s *AuthService) RefreshToken(refreshToken string) (*util.TokenPair, error) {
	db := database.Default()

	claims, err := s.jwtManager.ValidateToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Verify user still exists and is active
	var user model.User
	err = db.Where("id = ? AND status = ?", claims.UserID, model.UserStatusActive).First(&user).Error
	if err != nil {
		return nil, errors.New("user not found or inactive")
	}

	// Generate new tokens
	tokenPair, err := s.jwtManager.GenerateTokenPair(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return tokenPair, nil
}

// Logout handles user logout.
func (s *AuthService) Logout(userID uint, token string) error {
	db := database.Default()
	tokenHash := util.HashToken(token)
	return db.Model(&model.Session{}).
		Where("token_hash = ? AND user_id = ?", tokenHash, userID).
		Update("revoked", true).
		Error
}

// GetCurrentUser retrieves current user by ID.
func (s *AuthService) GetCurrentUser(userID uint) (*model.User, error) {
	db := database.Default()
	var user model.User
	err := db.Preload("OAuthIdentities").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetLinkedAccounts retrieves all linked accounts for a user.
func (s *AuthService) GetLinkedAccounts(userID uint) ([]model.OAuthIdentity, error) {
	db := database.Default()
	var identities []model.OAuthIdentity
	err := db.Where("user_id = ?", userID).Find(&identities).Error
	if err != nil {
		return nil, err
	}
	return identities, nil
}

// BindEmail binds an email to OAuth-only user.
func (s *AuthService) BindEmail(userID uint, email, code string) error {
	db := database.Default()

	// Check if email is already used by another user
	var existingUser model.User
	err := db.Where("email = ? AND id != ?", email, userID).First(&existingUser).Error
	if err == nil {
		return errors.New("该邮箱已被其他账号使用")
	}

	// Verify code
	var verificationCode model.EmailVerificationCode
	err = db.Where("email = ? AND type = ? AND used = false AND expires_at > ?",
		email, model.CodeTypeBindEmail, time.Now()).
		Order("created_at DESC").
		First(&verificationCode).Error

	if err != nil {
		return errors.New("验证码无效或已过期")
	}

	if verificationCode.Code != code {
		verificationCode.Attempts++
		db.Save(&verificationCode)
		return errors.New("验证码错误")
	}

	// Mark code as used
	verificationCode.Used = true
	db.Save(&verificationCode)

	// Update user
	now := time.Now()
	return db.Model(&model.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"email":             email,
			"email_verified":    true,
			"email_verified_at": now,
		}).Error
}

// SetPrimaryAccount sets a specific account as primary.
func (s *AuthService) SetPrimaryAccount(userID uint, provider string, providerUserID *string) error {
	db := database.Default()
	tx := db.Begin()

	// Unset all primary accounts for this user
	if err := tx.Model(&model.OAuthIdentity{}).
		Where("user_id = ?", userID).
		Update("is_primary", false).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	// If setting email as primary, we need to handle differently
	if provider == "email" {
		tx.Rollback()
		return nil
	}

	// Set new primary OAuth account
	if err := tx.Model(&model.OAuthIdentity{}).
		Where("user_id = ? AND provider = ?", userID, provider).
		Update("is_primary", true).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

// UnlinkOAuth unlinks an OAuth account from user.
func (s *AuthService) UnlinkOAuth(userID uint, provider string, providerUserID *string) error {
	db := database.Default()

	// Check if this is the only login method
	var emailLoginCount int64
	db.Model(&model.User{}).
		Where("id = ? AND email != ''", userID).
		Count(&emailLoginCount)

	var oauthCount int64
	db.Model(&model.OAuthIdentity{}).
		Where("user_id = ?", userID).
		Count(&oauthCount)

	if emailLoginCount == 0 && oauthCount <= 1 {
		return errors.New("无法解绑最后一个登录方式")
	}

	// Unlink the OAuth account
	return db.Where("user_id = ? AND provider = ?", userID, provider).
		Delete(&model.OAuthIdentity{}).Error
}

// UpdateProfile updates user profile.
func (s *AuthService) UpdateProfile(userID uint, updates map[string]interface{}) error {
	db := database.Default()
	return db.Model(&model.User{}).
		Where("id = ?", userID).
		Updates(updates).
		Error
}

// Helper functions

func (s *AuthService) logLoginHistory(user *model.User, method model.LoginMethod, success bool, failureReason, ipAddress, userAgent string) {
	db := database.Default()
	history := &model.LoginHistory{
		UserID:        &user.ID,
		LoginMethod:   method,
		Success:       success,
		FailureReason: &failureReason,
		IPAddress:     &ipAddress,
		UserAgent:     &userAgent,
	}
	db.Create(history)
}

func (s *AuthService) createSession(user *model.User, token, ipAddress, userAgent string) error {
	db := database.Default()
	tokenHash := util.HashToken(token)
	expiresAt := time.Now().Add(24 * time.Hour)

	session := &model.Session{
		UserID:         user.ID,
		TokenHash:      tokenHash,
		IPAddress:      &ipAddress,
		UserAgent:      &userAgent,
		ExpiresAt:      expiresAt,
		LastActivityAt: &[]time.Time{time.Now()}[0],
	}

	return db.Create(session).Error
}

// ValidateAccessToken validates an access token and returns user ID if valid.
func (s *AuthService) ValidateAccessToken(tokenString string) (uint, error) {
	db := database.Default()

	claims, err := s.jwtManager.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}

	// Check if session exists and is not revoked
	tokenHash := util.HashToken(tokenString)
	var session model.Session
	err = db.Where("token_hash = ? AND revoked = false AND expires_at > ?",
		tokenHash, time.Now()).First(&session).Error
	if err != nil {
		return 0, errors.New("session not found or expired")
	}

	return claims.UserID, nil
}
