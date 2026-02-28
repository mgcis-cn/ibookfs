// Package auth provides authentication business logic.
package auth

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/mgcis-cn/ibookfs/internal/apiserver/biz/email"
	model2 "github.com/mgcis-cn/ibookfs/internal/apiserver/model"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/store"
	"github.com/mgcis-cn/ibookfs/internal/pkg/util"
	"github.com/mgcis-cn/ibookfs/pkg/authn/jwt"
)

// AuthBiz defines the interface for authentication business logic.
type AuthBiz interface {
	SendCode(ctx context.Context, req SendCodeRequest) error
	Login(ctx context.Context, req LoginRequest, ipAddress, userAgent string) (*AuthResponse, error)
	Register(ctx context.Context, req RegisterRequest, ipAddress, userAgent string) (*AuthResponse, error)
	GenerateTokenPair(ctx context.Context, userID uint, email, role string) (*jwt.TokenPair, error)
	RefreshToken(ctx context.Context, refreshToken string) (*jwt.TokenPair, error)
	Logout(ctx context.Context, userID uint, token string) error
	GetCurrentUser(ctx context.Context, userID uint) (*model2.User, error)
	GetLinkedAccounts(ctx context.Context, userID uint) ([]model2.OAuthIdentity, error)
	BindEmail(ctx context.Context, userID uint, email, code string) error
	SetPrimaryAccount(ctx context.Context, userID uint, provider string, providerUserID *string) error
	UnlinkOAuth(ctx context.Context, userID uint, provider string, providerUserID *string) error
	UpdateProfile(ctx context.Context, userID uint, updates map[string]interface{}) error
	ValidateAccessToken(ctx context.Context, tokenString string) (uint, error)
	CreateSession(ctx context.Context, user *model2.User, token, ipAddress, userAgent string) error
}

// authBiz is the concrete implementation of AuthBiz.
type authBiz struct {
	jwtManager *jwt.Auth
	emailBiz   email.EmailBiz
	repo       store.IStore
}

// NewAuthBiz creates a new auth business logic with injected dependencies.
func NewAuthBiz(jwt *jwt.Auth, emailBiz email.EmailBiz, repo store.IStore) AuthBiz {
	return &authBiz{
		jwtManager: jwt,
		emailBiz:   emailBiz,
		repo:       repo,
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
	User           *model2.User           `json:"user"`
	Token          string                 `json:"token"`
	RefreshToken   string                 `json:"refresh_token,omitempty"`
	ExpiresIn      int64                  `json:"expires_in"`
	LinkedAccounts []model2.OAuthIdentity `json:"linked_accounts,omitempty"`
	IsNewUser      bool                   `json:"is_new_user,omitempty"`
}

// SendCode sends a verification code to email.
func (s *authBiz) SendCode(ctx context.Context, req SendCodeRequest) error {
	userStore := s.repo.User()

	// Check rate limiting
	recentCount, err := userStore.CountRecentCodes(ctx, req.Email, time.Now().Add(-60*time.Second))
	if err != nil {
		return fmt.Errorf("failed to check rate limit: %w", err)
	}

	if recentCount > 0 {
		return errors.New("请等待60秒后重新发送")
	}

	// Generate verification code based on email configuration
	var code string
	emailEnabled := s.emailBiz.IsEmailEnabled()

	if !emailEnabled {
		// Email not properly configured - use demo code
		code = "123456"
		fmt.Printf("[DEMO MODE] Email not configured, using fixed verification code: %s\n", code)
	} else {
		// Email is configured - generate random code
		code = util.GenerateVerificationCode()
	}

	// Store code
	verificationCode := &model2.EmailVerificationCode{
		Email:     req.Email,
		Code:      code,
		Type:      model2.CodeType(req.Type),
		ExpiresAt: time.Now().Add(5 * time.Minute),
	}

	if err := userStore.CreateVerificationCode(ctx, verificationCode); err != nil {
		return fmt.Errorf("failed to store verification code: %w", err)
	}

	// Send email with verification code
	if emailEnabled {
		// Email is configured - send real email and fail on error
		if err := s.emailBiz.SendVerificationCode(req.Email, code, req.Type); err != nil {
			// Delete the verification code since email failed
			_ = userStore.DeleteVerificationCode(ctx, verificationCode)
			return fmt.Errorf("failed to send email: %w", err)
		}
	} else {
		// Email not configured - log only
		fmt.Printf("[DEMO MODE] Verification code for %s: %s (email sending skipped - not configured)\n", req.Email, code)
	}

	return nil
}

// Login authenticates user with email and verification code.
func (s *authBiz) Login(ctx context.Context, req LoginRequest, ipAddress, userAgent string) (*AuthResponse, error) {
	userStore := s.repo.User()

	// Find valid verification code
	code, err := userStore.GetLatestVerificationCode(ctx, req.Email, model2.CodeTypeLogin)
	if err != nil {
		return nil, errors.New("验证码无效或已过期")
	}

	// Verify code
	if code.Code != req.Code {
		code.Attempts++
		_ = userStore.UpdateVerificationCode(ctx, code)
		return nil, errors.New("验证码错误")
	}

	// Mark code as used
	code.Used = true
	_ = userStore.UpdateVerificationCode(ctx, code)

	// Find user
	user, err := userStore.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("用户不存在，请先注册")
	}

	// Check user status
	if user.Status != model2.UserStatusActive {
		return nil, errors.New("账号已被禁用")
	}

	// Generate tokens
	tokenPair, err := s.jwtManager.Sign(fmt.Sprintf("%d", user.ID), user.Email, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Update last login
	now := time.Now()
	user.LastLoginAt = &now
	user.LastLoginIP = &ipAddress
	_ = userStore.UpdateUserModel(ctx, user)

	// Log login history
	s.logLoginHistory(ctx, user, model2.LoginMethodEmail, true, "", ipAddress, userAgent)

	// Load OAuth identities
	oauthIdentities, _ := userStore.ListOAuthIdentities(ctx, user.ID)

	// Create session
	_ = s.CreateSession(ctx, user, tokenPair.Token, ipAddress, userAgent)

	return &AuthResponse{
		User:           user,
		Token:          tokenPair.Token,
		RefreshToken:   tokenPair.Token,
		ExpiresIn:      tokenPair.ExpiresAt,
		LinkedAccounts: oauthIdentities,
	}, nil
}

// Register creates a new user account.
func (s *authBiz) Register(ctx context.Context, req RegisterRequest, ipAddress, userAgent string) (*AuthResponse, error) {
	userStore := s.repo.User()

	// Check if email already exists
	_, err := userStore.GetUserByEmail(ctx, req.Email)
	if err == nil {
		return nil, errors.New("该邮箱已被注册")
	}

	// Verify code
	code, err := userStore.GetLatestVerificationCode(ctx, req.Email, model2.CodeTypeRegister)
	if err != nil {
		return nil, errors.New("验证码无效或已过期")
	}

	if code.Code != req.Code {
		code.Attempts++
		_ = userStore.UpdateVerificationCode(ctx, code)
		return nil, errors.New("验证码错误")
	}

	// Mark code as used
	code.Used = true
	_ = userStore.UpdateVerificationCode(ctx, code)

	// Create user
	displayName := req.FirstName + " " + req.LastName
	user := &model2.User{
		Email:       req.Email,
		FirstName:   req.FirstName,
		LastName:    req.LastName,
		DisplayName: &displayName,
		Status:      model2.UserStatusActive,
		Role:        model2.UserRoleUser,
		Language:    "zh-CN",
		Timezone:    "Asia/Shanghai",
	}

	if err := userStore.CreateUser(ctx, user); err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	// Generate tokens
	tokenPair, err := s.jwtManager.Sign(fmt.Sprintf("%d", user.ID), user.Email, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Log login history
	s.logLoginHistory(ctx, user, model2.LoginMethodEmail, true, "", ipAddress, userAgent)

	// Create session
	_ = s.CreateSession(ctx, user, tokenPair.Token, ipAddress, userAgent)

	return &AuthResponse{
		User:         user,
		Token:        tokenPair.Token,
		RefreshToken: tokenPair.Token,
		ExpiresIn:    tokenPair.ExpiresAt,
		IsNewUser:    true,
	}, nil
}

// GenerateTokenPair generates access and refresh tokens for a user.
func (s *authBiz) GenerateTokenPair(ctx context.Context, userID uint, email, role string) (*jwt.TokenPair, error) {
	return s.jwtManager.Sign(fmt.Sprintf("%d", userID), email, role)
}

// RefreshToken refreshes an access token using refresh token.
func (s *authBiz) RefreshToken(ctx context.Context, refreshToken string) (*jwt.TokenPair, error) {
	userStore := s.repo.User()

	claims, err := s.jwtManager.ParseClaims(ctx, refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Verify user still exists and is active
	userID, err := strconv.ParseUint(claims.UserID, 10, 64)
	if err != nil {
		return nil, errors.New("invalid user ID in token")
	}

	user, err := userStore.GetUserByID(ctx, uint(userID))
	if err != nil {
		return nil, errors.New("user not found")
	}

	if user.Status != model2.UserStatusActive {
		return nil, errors.New("user not found or inactive")
	}

	// Generate new tokens
	tokenPair, err := s.jwtManager.Sign(fmt.Sprintf("%d", user.ID), user.Email, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	return tokenPair, nil
}

// Logout handles user logout.
func (s *authBiz) Logout(ctx context.Context, userID uint, token string) error {
	tokenHash := util.HashToken(token)
	userStore := s.repo.User()
	return userStore.RevokeSession(ctx, tokenHash, userID)
}

// GetCurrentUser retrieves current user by ID.
func (s *authBiz) GetCurrentUser(ctx context.Context, userID uint) (*model2.User, error) {
	return s.repo.User().GetUserWithPreloads(ctx, userID)
}

// GetLinkedAccounts retrieves all linked accounts for a user.
func (s *authBiz) GetLinkedAccounts(ctx context.Context, userID uint) ([]model2.OAuthIdentity, error) {
	return s.repo.User().ListOAuthIdentities(ctx, userID)
}

// BindEmail binds an email to OAuth-only user.
func (s *authBiz) BindEmail(ctx context.Context, userID uint, email, code string) error {
	userStore := s.repo.User()

	// Check if email is already used by another user
	count, err := userStore.CountUsersWithEmail(ctx, email, &userID)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("该邮箱已被其他账号使用")
	}

	// Verify code
	verificationCode, err := userStore.GetLatestVerificationCode(ctx, email, model2.CodeTypeBindEmail)
	if err != nil {
		return errors.New("验证码无效或已过期")
	}

	if verificationCode.Code != code {
		verificationCode.Attempts++
		_ = userStore.UpdateVerificationCode(ctx, verificationCode)
		return errors.New("验证码错误")
	}

	// Mark code as used
	verificationCode.Used = true
	_ = userStore.UpdateVerificationCode(ctx, verificationCode)

	// Update user
	return userStore.UpdateUser(ctx, userID, map[string]interface{}{
		"email": email,
	})
}

// SetPrimaryAccount sets a specific account as primary.
func (s *authBiz) SetPrimaryAccount(ctx context.Context, userID uint, provider string, providerUserID *string) error {
	userStore := s.repo.User()

	// Unset all primary accounts for this user
	if err := userStore.UpdateAllOAuthIdentities(ctx, userID, map[string]interface{}{
		"is_primary": false,
	}); err != nil {
		return err
	}

	// If setting email as primary, we need to handle differently
	if provider == "email" {
		return nil
	}

	// Set new primary OAuth account
	return userStore.UpdateOAuthIdentity(ctx, userID, provider, map[string]interface{}{
		"is_primary": true,
	})
}

// UnlinkOAuth unlinks an OAuth account from user.
func (s *authBiz) UnlinkOAuth(ctx context.Context, userID uint, provider string, providerUserID *string) error {
	userStore := s.repo.User()

	// Check if this is the only login method
	user, _ := userStore.GetUserByID(ctx, userID)
	emailLoginCount := 0
	if user.Email != "" {
		emailLoginCount = 1
	}

	oauthCount, _ := userStore.CountOAuthIdentities(ctx, userID)

	if emailLoginCount == 0 && oauthCount <= 1 {
		return errors.New("无法解绑最后一个登录方式")
	}

	// Unlink the OAuth account
	return userStore.DeleteOAuthIdentity(ctx, userID, provider)
}

// UpdateProfile updates user profile.
func (s *authBiz) UpdateProfile(ctx context.Context, userID uint, updates map[string]interface{}) error {
	return s.repo.User().UpdateUser(ctx, userID, updates)
}

// ValidateAccessToken validates an access token and returns user ID if valid.
func (s *authBiz) ValidateAccessToken(ctx context.Context, tokenString string) (uint, error) {
	userStore := s.repo.User()

	claims, err := s.jwtManager.ParseClaims(ctx, tokenString)
	if err != nil {
		return 0, err
	}

	// Check if session exists and is not revoked
	tokenHash := util.HashToken(tokenString)
	_, err = userStore.GetValidSession(ctx, tokenHash)
	if err != nil {
		return 0, errors.New("session not found or expired")
	}

	userId, err := strconv.ParseUint(claims.UserID, 10, 64)
	return uint(userId), err
}

// Helper functions

func (s *authBiz) logLoginHistory(ctx context.Context, user *model2.User, method model2.LoginMethod, success bool, failureReason, ipAddress, userAgent string) {
	userStore := s.repo.User()
	history := &model2.LoginHistory{
		UserID:        &user.ID,
		LoginMethod:   method,
		Success:       success,
		FailureReason: &failureReason,
		IPAddress:     &ipAddress,
		UserAgent:     &userAgent,
	}
	_ = userStore.CreateLoginHistory(ctx, history)
}

func (s *authBiz) CreateSession(ctx context.Context, user *model2.User, token, ipAddress, userAgent string) error {
	userStore := s.repo.User()
	tokenHash := util.HashToken(token)
	expiresAt := time.Now().Add(24 * time.Hour)
	now := time.Now()

	session := &model2.Session{
		UserID:         user.ID,
		TokenHash:      tokenHash,
		IPAddress:      &ipAddress,
		UserAgent:      &userAgent,
		ExpiresAt:      expiresAt,
		LastActivityAt: &now,
	}

	return userStore.CreateSession(ctx, session)
}
