// Package handler contains HTTP request handlers.
package handler

import (
	"context"
	"crypto/rand"
	"encoding/hex"

	"github.com/mgcis-cn/ibookfs/internal/apiserver/biz/auth"
	apierr "github.com/mgcis-cn/ibookfs/internal/apiserver/errors"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/model"
	v1 "github.com/mgcis-cn/ibookfs/pkg/api/apiserver/v1"
	contextx "github.com/mgcis-cn/ibookfs/pkg/context"
)

type AuthRouter interface {
	SendCode(ctx context.Context, req *v1.SendCodeRequest) (*v1.SendCodeResponse, error)
	Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error)
	Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error)
	RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenResponse, error)
	Logout(ctx context.Context, req *v1.LogoutRequest) (*v1.LogoutResponse, error)
	GetCurrentUser(ctx context.Context, req *v1.GetCurrentUserRequest) (*v1.GetCurrentUserResponse, error)
	GetLinkedAccounts(ctx context.Context, req *v1.GetLinkedAccountsRequest) (*v1.GetLinkedAccountsResponse, error)
	LinkOAuth(ctx context.Context, req *v1.LinkOAuthRequest) (*v1.LinkOAuthResponse, error)
	UnlinkOAuth(ctx context.Context, req *v1.UnlinkOAuthRequest) (*v1.UnlinkOAuthResponse, error)
	BindEmail(ctx context.Context, req *v1.BindEmailRequest) (*v1.BindEmailResponse, error)
	SetPrimaryAccount(ctx context.Context, req *v1.SetPrimaryAccountRequest) (*v1.SetPrimaryAccountResponse, error)
	UpdateProfile(ctx context.Context, req *v1.UpdateProfileRequest) (*v1.UpdateProfileResponse, error)
}

// APIResponse represents a standard API response.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// SendCode handles POST /auth/send-code
func (h *handler) SendCode(ctx context.Context, req *v1.SendCodeRequest) (*v1.SendCodeResponse, error) {
	// Validate type
	validTypes := map[string]bool{
		"login":          true,
		"register":       true,
		"reset_password": true,
		"bind_email":     true,
	}
	if !validTypes[req.Type] {
		return &v1.SendCodeResponse{}, apierr.ErrInvalidCodeType
	}

	// Convert v1 request to biz request
	bizReq := auth.SendCodeRequest{
		Email: req.Email,
		Type:  req.Type,
	}

	if err := h.biz.Auth().SendCode(ctx, bizReq); err != nil {
		return &v1.SendCodeResponse{}, err
	}

	return &v1.SendCodeResponse{
		ExpiresIn: 300,
		Message:   "验证码已发送",
	}, nil
}

// Login handles POST /auth/login
func (h *handler) Login(ctx context.Context, req *v1.LoginRequest) (*v1.LoginResponse, error) {
	// Convert v1 request to biz request
	bizReq := auth.LoginRequest{
		Email: req.Email,
		Code:  req.Code,
	}

	c := contextx.Request(ctx)
	ipAddress := ""
	userAgent := c.UserAgent()

	bizResp, err := h.biz.Auth().Login(ctx, bizReq, ipAddress, userAgent)
	if err != nil {
		return nil, err
	}

	// Convert []model.OAuthIdentity to []any
	linkedAccounts := make([]any, len(bizResp.LinkedAccounts))
	for i, acc := range bizResp.LinkedAccounts {
		linkedAccounts[i] = acc
	}

	return &v1.LoginResponse{
		User:           bizResp.User,
		Token:          bizResp.Token,
		RefreshToken:   bizResp.RefreshToken,
		ExpiresIn:      bizResp.ExpiresIn,
		LinkedAccounts: linkedAccounts,
	}, nil
}

// Register handles POST /auth/register
func (h *handler) Register(ctx context.Context, req *v1.RegisterRequest) (*v1.RegisterResponse, error) {
	// Convert v1 request to biz request
	bizReq := auth.RegisterRequest{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Code:      req.Code,
	}

	c := contextx.Request(ctx)
	ipAddress := ""
	userAgent := c.UserAgent()

	bizResp, err := h.biz.Auth().Register(ctx, bizReq, ipAddress, userAgent)
	if err != nil {
		return nil, err
	}

	return &v1.RegisterResponse{
		User:      bizResp.User,
		Token:     bizResp.Token,
		ExpiresIn: bizResp.ExpiresIn,
		IsNewUser: bizResp.IsNewUser,
	}, nil
}

// OAuthCallback handles POST /auth/oauth/callback
func (h *handler) OAuthCallback(ctx context.Context, req *v1.OAuthCallbackRequest) (*v1.OAuthCallbackResponse, error) {

	c := contextx.Request(ctx)

	ipAddress := ""
	userAgent := c.UserAgent()
	redirectURI := c.Header.Get("Origin") + "/auth/callback"

	resp, err := h.biz.OAuth().HandleCallback(
		ctx,
		model.OAuthProvider(req.Provider),
		req.Code,
		req.State,
		redirectURI,
		ipAddress,
		userAgent,
	)

	if err != nil {
		return &v1.OAuthCallbackResponse{}, err
	}

	return &v1.OAuthCallbackResponse{Data: resp}, nil
}

// RefreshToken handles POST /auth/refresh
func (h *handler) RefreshToken(ctx context.Context, req *v1.RefreshTokenRequest) (*v1.RefreshTokenResponse, error) {
	tokenPair, err := h.biz.Auth().RefreshToken(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	return &v1.RefreshTokenResponse{
		Token:        tokenPair.Token,
		RefreshToken: tokenPair.Token,
		ExpiresIn:    tokenPair.ExpiresAt,
	}, nil
}

// Logout handles POST /auth/logout
func (h *handler) Logout(ctx context.Context, req *v1.LogoutRequest) (*v1.LogoutResponse, error) {
	err := h.biz.Auth().Logout(ctx, uint(contextx.UserId(ctx)), req.Token)
	if err != nil {
		return &v1.LogoutResponse{}, err
	}
	return &v1.LogoutResponse{}, nil
}

// GetCurrentUser handles GET /auth/me
func (h *handler) GetCurrentUser(ctx context.Context, req *v1.GetCurrentUserRequest) (*v1.GetCurrentUserResponse, error) {
	userID := uint(contextx.UserId(ctx))
	user, err := h.biz.Auth().GetCurrentUser(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Load linked accounts
	linkedAccounts, _ := h.biz.Auth().GetLinkedAccounts(ctx, userID)

	// Convert []model.OAuthIdentity to []any
	accounts := make([]any, len(linkedAccounts))
	for i, acc := range linkedAccounts {
		accounts[i] = acc
	}

	return &v1.GetCurrentUserResponse{
		User:           user,
		LinkedAccounts: accounts,
	}, nil
}

// GetLinkedAccounts handles GET /auth/accounts
func (h *handler) GetLinkedAccounts(ctx context.Context, req *v1.GetLinkedAccountsRequest) (*v1.GetLinkedAccountsResponse, error) {
	userID := uint(contextx.UserId(ctx))
	accounts, err := h.biz.Auth().GetLinkedAccounts(ctx, userID)
	if err != nil {
		return nil, err
	}

	// Add email account if user has email
	user, _ := h.biz.Auth().GetCurrentUser(ctx, userID)
	var result []interface{}

	if user.Email != "" {
		result = append(result, map[string]interface{}{
			"provider":       "email",
			"provider_email": user.Email,
			"is_primary":     true,
			"linked_at":      user.CreatedAt,
		})
	}

	for _, acc := range accounts {
		result = append(result, map[string]interface{}{
			"provider":          acc.Provider,
			"provider_user_id":  &acc.ProviderUserID,
			"provider_username": acc.ProviderUsername,
			"provider_email":    acc.ProviderEmail,
			"is_primary":        acc.IsPrimary,
			"linked_at":         acc.LinkedAt,
			"last_used_at":      acc.LastUsedAt,
		})
	}

	return &v1.GetLinkedAccountsResponse{
		Accounts: result,
	}, nil
}

// LinkOAuth handles POST /auth/oauth/link
func (h *handler) LinkOAuth(ctx context.Context, req *v1.LinkOAuthRequest) (*v1.LinkOAuthResponse, error) {
	userID := uint(contextx.UserId(ctx))
	c := contextx.Request(ctx)
	redirectURI := c.Header.Get("Origin") + "/settings/accounts/link"

	err := h.biz.OAuth().LinkOAuth(ctx, userID, model.OAuthProvider(req.Provider), req.Code, redirectURI)
	if err != nil {
		return &v1.LinkOAuthResponse{}, err
	}

	return &v1.LinkOAuthResponse{
		Message: "OAuth account linked successfully",
	}, nil
}

// UnlinkOAuth handles DELETE /auth/oauth/unlink
func (h *handler) UnlinkOAuth(ctx context.Context, req *v1.UnlinkOAuthRequest) (*v1.UnlinkOAuthResponse, error) {
	userID := uint(contextx.UserId(ctx))

	err := h.biz.Auth().UnlinkOAuth(ctx, userID, req.Provider, req.ProviderUserID)
	if err != nil {
		return &v1.UnlinkOAuthResponse{}, err
	}

	return &v1.UnlinkOAuthResponse{
		Message: "OAuth account unlinked successfully",
	}, nil
}

// BindEmail handles POST /auth/bind-email
func (h *handler) BindEmail(ctx context.Context, req *v1.BindEmailRequest) (*v1.BindEmailResponse, error) {
	userID := uint(contextx.UserId(ctx))

	err := h.biz.Auth().BindEmail(ctx, userID, req.Email, req.Code)
	if err != nil {
		return &v1.BindEmailResponse{}, err
	}

	user, _ := h.biz.Auth().GetCurrentUser(ctx, userID)

	return &v1.BindEmailResponse{
		User: user,
	}, nil
}

// SetPrimaryAccount handles PUT /auth/accounts/primary
func (h *handler) SetPrimaryAccount(ctx context.Context, req *v1.SetPrimaryAccountRequest) (*v1.SetPrimaryAccountResponse, error) {
	userID := uint(contextx.UserId(ctx))

	err := h.biz.Auth().SetPrimaryAccount(ctx, userID, req.Provider, req.ProviderUserID)
	if err != nil {
		return &v1.SetPrimaryAccountResponse{}, err
	}

	return &v1.SetPrimaryAccountResponse{
		Message: "Primary account updated",
	}, nil
}

// UpdateProfile handles PUT /auth/me
func (h *handler) UpdateProfile(ctx context.Context, req *v1.UpdateProfileRequest) (*v1.UpdateProfileResponse, error) {
	userID := uint(contextx.UserId(ctx))

	updates := make(map[string]interface{})
	if req.FirstName != nil {
		updates["first_name"] = *req.FirstName
	}
	if req.LastName != nil {
		updates["last_name"] = *req.LastName
	}
	if req.DisplayName != nil {
		updates["display_name"] = *req.DisplayName
	}
	if req.Language != nil {
		updates["language"] = *req.Language
	}
	if req.Timezone != nil {
		updates["timezone"] = *req.Timezone
	}

	err := h.biz.Auth().UpdateProfile(ctx, userID, updates)
	if err != nil {
		return &v1.UpdateProfileResponse{}, err
	}

	user, _ := h.biz.Auth().GetCurrentUser(ctx, userID)

	return &v1.UpdateProfileResponse{
		User: user,
	}, nil
}

// GetConfig handles GET /auth/config
// Returns the authentication configuration including supported OAuth providers and email status
func (h *handler) GetConfig(ctx context.Context, req *v1.GetConfigRequest) (*v1.GetConfigResponse, error) {
	// Build list of supported OAuth providers from OAuth Factory
	supportedProviders := h.oauthFactory.GetProviderKinds()

	// Check if email is enabled
	emailEnabled := h.biz.Email().IsEmailEnabled()

	// Helper function to safely check if provider is configured
	isConfigured := func(name string) bool {
		_, ok := h.oauthFactory.Get(name)
		return ok
	}

	// Check if any OAuth provider is enabled
	oauthEnabled := len(supportedProviders) > 0

	return &v1.GetConfigResponse{
		OauthProviders: supportedProviders,
		EmailEnabled:   emailEnabled,
		OAuthEnabled:   oauthEnabled,
		OAuth: map[string]interface{}{
			"github": map[string]any{
				"configured": isConfigured("github"),
			},
			"gitee": map[string]any{
				"configured": isConfigured("gitee"),
			},
			"icloud": map[string]any{
				"configured": isConfigured("icloud"),
			},
			"google": map[string]any{
				"configured": isConfigured("google"),
			},
			"wechat": map[string]any{
				"configured": isConfigured("wechat"),
			},
		},
	}, nil
}

// OAuthAuthorize handles GET /auth/oauth/authorize
func (h *handler) OAuthAuthorize(ctx context.Context, req *v1.OAuthAuthorizeRequest) (*v1.OAuthAuthorizeResponse, error) {

	if req.Provider == "" {
		return &v1.OAuthAuthorizeResponse{}, apierr.ErrProviderRequired
	}
	// redirect_uri is controlled by config file, no code concatenation
	// State is auto-generated if empty
	if req.State == "" {
		stateBytes := make([]byte, 32)
		if _, err := rand.Read(stateBytes); err != nil {
			return &v1.OAuthAuthorizeResponse{}, apierr.ErrStateGenFailed
		}
		req.State = hex.EncodeToString(stateBytes)
	}

	authorizeURL, err := h.biz.OAuth().GetAuthorizeURL(ctx, model.OAuthProvider(req.Provider), req.RedirectURI, req.State)
	if err != nil {
		return &v1.OAuthAuthorizeResponse{}, err
	}

	return &v1.OAuthAuthorizeResponse{
		AuthorizeURL: authorizeURL,
		State:        req.State,
	}, nil
}
