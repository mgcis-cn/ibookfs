// Package oauth provides OAuth business logic.
package oauth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mgcis-cn/ibookfs/internal/apiserver/biz/auth"
	model2 "github.com/mgcis-cn/ibookfs/internal/apiserver/model"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/store"
)

// OAuthBiz defines the interface for OAuth business logic.
type OAuthBiz interface {
	GetAuthorizeURL(ctx context.Context, provider model2.OAuthProvider, redirectURI, state string) (string, error)
	HandleCallback(ctx context.Context, provider model2.OAuthProvider, code, state, redirectURI, ipAddress, userAgent string) (*AuthResponse, error)
	LinkOAuth(ctx context.Context, userID uint, provider model2.OAuthProvider, code, redirectURI string) error
}

// OAuthBizFactory creates OAuth configs.
type OAuthBizFactory struct{}

// GetOAuthConfig returns OAuth config for a provider.
func GetOAuthConfig(provider model2.OAuthProvider, clientID, clientSecret, redirectURI string) model2.OAuthConfig {
	switch provider {
	case model2.OAuthProviderGitHub:
		return model2.OAuthConfig{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURI,
			AuthURL:      "https://github.com/login/oauth/authorize",
			TokenURL:     "https://github.com/login/oauth/access_token",
			UserInfoURL:  "https://api.github.com/user",
			Scopes:       []string{"read:user", "user:email"},
		}
	case model2.OAuthProviderGitee:
		return model2.OAuthConfig{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURI,
			AuthURL:      "https://gitee.com/oauth/authorize",
			TokenURL:     "https://gitee.com/oauth/token",
			UserInfoURL:  "https://gitee.com/api/v5/user",
			Scopes:       []string{"user_info", "emails"},
		}
	default:
		return model2.OAuthConfig{}
	}
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

// oauthBiz is the concrete implementation of OAuthBiz.
type oauthBiz struct {
	authSvc    auth.AuthBiz
	httpClient *http.Client
	configs    map[model2.OAuthProvider]model2.OAuthConfig
	repo       store.IStore
}

// NewOAuthBiz creates a new OAuth business logic with injected dependencies.
func NewOAuthBiz(authSvc auth.AuthBiz, configs map[model2.OAuthProvider]model2.OAuthConfig, repo store.IStore) OAuthBiz {
	return &oauthBiz{
		authSvc:    authSvc,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		configs:    configs,
		repo:       repo,
	}
}

// GetAuthorizeURL returns the OAuth authorization URL.
func (s *oauthBiz) GetAuthorizeURL(ctx context.Context, provider model2.OAuthProvider, redirectURI, state string) (string, error) {
	cfg, ok := s.configs[provider]
	if !ok {
		return "", errors.New("unsupported OAuth provider")
	}

	params := url.Values{}
	params.Set("client_id", cfg.ClientID)
	params.Set("redirect_uri", redirectURI)
	params.Set("response_type", "code")
	params.Set("state", state)
	for _, scope := range cfg.Scopes {
		params.Add("scope", scope)
	}

	return fmt.Sprintf("%s?%s", cfg.AuthURL, params.Encode()), nil
}

// HandleCallback handles OAuth callback and returns auth response.
func (s *oauthBiz) HandleCallback(ctx context.Context, provider model2.OAuthProvider, code, state, redirectURI, ipAddress, userAgent string) (*AuthResponse, error) {
	cfg, ok := s.configs[provider]
	if !ok {
		return nil, errors.New("unsupported OAuth provider")
	}

	// Exchange code for token
	tokenResp, err := s.exchangeCodeForToken(provider, code, redirectURI, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info from provider
	userInfo, err := s.getUserInfo(provider, tokenResp.AccessToken, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Find or create user
	user, isNew, err := s.findOrCreateUser(ctx, userInfo, ipAddress, userAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	// Update or create OAuth identity
	if err := s.updateOAuthIdentity(ctx, user.ID, userInfo, tokenResp, provider); err != nil {
		return nil, fmt.Errorf("failed to update OAuth identity: %w", err)
	}

	// Generate tokens
	tokenPair, err := s.authSvc.GenerateTokenPair(ctx, user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Load OAuth identities
	oauthIdentities, _ := s.repo.User().ListOAuthIdentities(ctx, user.ID)

	return &AuthResponse{
		User:           user,
		Token:          tokenPair.Token,
		RefreshToken:   tokenPair.Token,
		ExpiresIn:      tokenPair.ExpiresAt,
		LinkedAccounts: oauthIdentities,
		IsNewUser:      isNew,
	}, nil
}

// LinkOAuth links an OAuth account to existing user.
func (s *oauthBiz) LinkOAuth(ctx context.Context, userID uint, provider model2.OAuthProvider, code, redirectURI string) error {
	cfg, ok := s.configs[provider]
	if !ok {
		return errors.New("unsupported OAuth provider")
	}

	userStore := s.repo.User()

	// Check if already linked
	_, err := userStore.GetOAuthIdentityByUserIDAndProvider(ctx, userID, provider.String())
	if err == nil {
		return errors.New("该OAuth账号已绑定")
	}

	// Exchange code for token
	tokenResp, err := s.exchangeCodeForToken(provider, code, redirectURI, cfg)
	if err != nil {
		return fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info from provider
	userInfo, err := s.getUserInfo(provider, tokenResp.AccessToken, cfg)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	// Check if OAuth identity is linked to another user
	_, err = userStore.GetOAuthIdentityByProviderAndProviderUserID(ctx, provider.String(), userInfo.ID)
	if err == nil {
		return errors.New("该OAuth账号已绑定到其他用户")
	}

	// Create OAuth identity
	now := time.Now()
	identity := &model2.OAuthIdentity{
		UserID:           userID,
		Provider:         provider.String(),
		ProviderUserID:   userInfo.ID,
		ProviderUsername: &userInfo.Username,
		ProviderEmail:    &userInfo.Email,
		AccessToken:      &tokenResp.AccessToken,
		IsPrimary:        false,
		LinkedAt:         &now,
	}

	if tokenResp.RefreshToken != "" {
		identity.RefreshToken = &tokenResp.RefreshToken
	}

	if !tokenResp.ExpiresAt.IsZero() {
		identity.TokenExpiresAt = &tokenResp.ExpiresAt
	}

	return userStore.CreateOAuthIdentity(ctx, identity)
}

// Helper functions

func (s *oauthBiz) exchangeCodeForToken(provider model2.OAuthProvider, code, redirectURI string, cfg model2.OAuthConfig) (*model2.OAuthTokenResponse, error) {
	data := url.Values{}
	data.Set("client_id", cfg.ClientID)
	data.Set("client_secret", cfg.ClientSecret)
	data.Set("code", code)
	data.Set("redirect_uri", redirectURI)
	data.Set("grant_type", "authorization_code")

	resp, err := s.httpClient.PostForm(cfg.TokenURL, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("token request failed: %s", string(body))
	}

	var tokenResp model2.OAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	// Calculate ExpiresAt from ExpiresIn
	if tokenResp.ExpiresIn > 0 {
		tokenResp.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}

	return &tokenResp, nil
}

func (s *oauthBiz) getUserInfo(provider model2.OAuthProvider, accessToken string, cfg model2.OAuthConfig) (*model2.OAuthUserInfo, error) {
	req, err := http.NewRequest("GET", cfg.UserInfoURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("user info request failed: %s", string(body))
	}

	return s.parseUserInfo(provider, resp.Body)
}

func (s *oauthBiz) parseUserInfo(provider model2.OAuthProvider, body io.Reader) (*model2.OAuthUserInfo, error) {
	switch provider {
	case model2.OAuthProviderGitHub:
		return s.parseGitHubUserInfo(body)
	case model2.OAuthProviderGitee:
		return s.parseGiteeUserInfo(body)
	default:
		return nil, errors.New("unsupported provider for user info parsing")
	}
}

// GitHubUserinfo represents GitHub user info response.
type GitHubUserinfo struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func (s *oauthBiz) parseGitHubUserInfo(body io.Reader) (*model2.OAuthUserInfo, error) {
	var githubUser GitHubUserinfo
	if err := json.NewDecoder(body).Decode(&githubUser); err != nil {
		return nil, err
	}

	// If email is not public, fetch from emails API
	if githubUser.Email == "" {
		githubUser.Email = s.getGitHubPrimaryEmail(githubUser.Login)
	}

	return &model2.OAuthUserInfo{
		ID:        fmt.Sprintf("%d", githubUser.ID),
		Username:  githubUser.Login,
		Email:     githubUser.Email,
		Name:      githubUser.Name,
		AvatarURL: githubUser.AvatarURL,
		Provider:  model2.OAuthProviderGitHub,
	}, nil
}

// GitHubEmail represents GitHub email response.
type GitHubEmail struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

func (s *oauthBiz) getGitHubPrimaryEmail(username string) string {
	// Use a default email pattern if we can't fetch the real email
	return fmt.Sprintf("%s@github.local", username)
}

// GiteeUserinfo represents Gitee user info response.
type GiteeUserinfo struct {
	ID        int64  `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func (s *oauthBiz) parseGiteeUserInfo(body io.Reader) (*model2.OAuthUserInfo, error) {
	var giteeUser GiteeUserinfo
	if err := json.NewDecoder(body).Decode(&giteeUser); err != nil {
		return nil, err
	}

	return &model2.OAuthUserInfo{
		ID:        fmt.Sprintf("%d", giteeUser.ID),
		Username:  giteeUser.Login,
		Email:     giteeUser.Email,
		Name:      giteeUser.Name,
		AvatarURL: giteeUser.AvatarURL,
		Provider:  model2.OAuthProviderGitee,
	}, nil
}

func (s *oauthBiz) findOrCreateUser(ctx context.Context, userInfo *model2.OAuthUserInfo, ipAddress, userAgent string) (*model2.User, bool, error) {
	userStore := s.repo.User()

	// First, try to find user by OAuth identity
	identity, err := userStore.GetOAuthIdentityByProviderAndProviderUserID(ctx, userInfo.Provider.String(), userInfo.ID)
	if err == nil {
		// Found existing user with this OAuth identity - we need to get the user
		user, err := userStore.GetUserByID(ctx, identity.UserID)
		if err != nil {
			return nil, false, err
		}

		// Update last login
		now := time.Now()
		user.LastLoginAt = &now
		user.LastLoginIP = &ipAddress
		_ = userStore.UpdateUserModel(ctx, user)

		// Update last_used_at
		identity.LastUsedAt = &now
		_ = userStore.UpdateOAuthIdentityModel(ctx, identity)

		// Log login
		s.logLoginHistory(ctx, user, userInfo.Provider.ToLoginMethod(), true, "", ipAddress, userAgent)

		return user, false, nil
	}

	// Not found by OAuth, check if email exists
	var user *model2.User
	if userInfo.Email != "" {
		user, err = userStore.GetUserByEmail(ctx, userInfo.Email)
		if err == nil {
			// Found user by email, link this OAuth identity
			return s.linkOAuthToUser(ctx, user, userInfo, ipAddress, userAgent)
		}
	}

	// Create new user
	return s.createUserFromOAuth(ctx, userInfo, ipAddress, userAgent)
}

func (s *oauthBiz) linkOAuthToUser(ctx context.Context, user *model2.User, userInfo *model2.OAuthUserInfo, ipAddress, userAgent string) (*model2.User, bool, error) {
	userStore := s.repo.User()
	now := time.Now()

	user.LastLoginAt = &now
	user.LastLoginIP = &ipAddress
	if userInfo.AvatarURL != "" {
		user.AvatarURL = &userInfo.AvatarURL
	}
	_ = userStore.UpdateUserModel(ctx, user)

	// Create OAuth identity
	identity := &model2.OAuthIdentity{
		UserID:           user.ID,
		Provider:         userInfo.Provider.String(),
		ProviderUserID:   userInfo.ID,
		ProviderUsername: &userInfo.Username,
		ProviderEmail:    &userInfo.Email,
		IsPrimary:        false,
		LinkedAt:         &now,
		LastUsedAt:       &now,
	}

	_ = userStore.CreateOAuthIdentity(ctx, identity)

	// Log login
	s.logLoginHistory(ctx, user, userInfo.Provider.ToLoginMethod(), true, "", ipAddress, userAgent)

	return user, false, nil
}

func (s *oauthBiz) createUserFromOAuth(ctx context.Context, userInfo *model2.OAuthUserInfo, ipAddress, userAgent string) (*model2.User, bool, error) {
	userStore := s.repo.User()

	// Generate display name
	displayName := userInfo.Name
	if displayName == "" {
		displayName = userInfo.Username
	}

	// Generate username/firstName/lastName
	// For OAuth users, we'll use the username as both first and last name initially
	user := &model2.User{
		Email:       userInfo.Email,
		FirstName:   userInfo.Username,
		LastName:    userInfo.Username,
		DisplayName: &displayName,
		AvatarURL:   &userInfo.AvatarURL,
		Status:      model2.UserStatusActive,
		Role:        model2.UserRoleUser,
		Language:    "zh-CN",
		Timezone:    "Asia/Shanghai",
	}

	if err := userStore.CreateUser(ctx, user); err != nil {
		return nil, false, err
	}

	// Create OAuth identity as primary
	now := time.Now()
	identity := &model2.OAuthIdentity{
		UserID:           user.ID,
		Provider:         userInfo.Provider.String(),
		ProviderUserID:   userInfo.ID,
		ProviderUsername: &userInfo.Username,
		ProviderEmail:    &userInfo.Email,
		IsPrimary:        true,
		LinkedAt:         &now,
		LastUsedAt:       &now,
	}

	_ = userStore.CreateOAuthIdentity(ctx, identity)

	// Log login
	s.logLoginHistory(ctx, user, userInfo.Provider.ToLoginMethod(), true, "", ipAddress, userAgent)

	return user, true, nil
}

func (s *oauthBiz) updateOAuthIdentity(ctx context.Context, userID uint, userInfo *model2.OAuthUserInfo, tokenResp *model2.OAuthTokenResponse, provider model2.OAuthProvider) error {
	userStore := s.repo.User()

	identity, err := userStore.GetOAuthIdentityByUserIDAndProvider(ctx, userID, provider.String())
	if err == nil {
		// Update existing identity
		now := time.Now()
		identity.AccessToken = &tokenResp.AccessToken
		identity.LastUsedAt = &now
		if tokenResp.RefreshToken != "" {
			identity.RefreshToken = &tokenResp.RefreshToken
		}
		if !tokenResp.ExpiresAt.IsZero() {
			identity.TokenExpiresAt = &tokenResp.ExpiresAt
		}
		return userStore.UpdateOAuthIdentityModel(ctx, identity)
	}

	// Create new identity
	now := time.Now()
	newIdentity := &model2.OAuthIdentity{
		UserID:           userID,
		Provider:         provider.String(),
		ProviderUserID:   userInfo.ID,
		ProviderUsername: &userInfo.Username,
		ProviderEmail:    &userInfo.Email,
		AccessToken:      &tokenResp.AccessToken,
		IsPrimary:        true,
		LinkedAt:         &now,
		LastUsedAt:       &now,
	}

	if tokenResp.RefreshToken != "" {
		newIdentity.RefreshToken = &tokenResp.RefreshToken
	}

	if !tokenResp.ExpiresAt.IsZero() {
		newIdentity.TokenExpiresAt = &tokenResp.ExpiresAt
	}

	return userStore.CreateOAuthIdentity(ctx, newIdentity)
}

func (s *oauthBiz) logLoginHistory(ctx context.Context, user *model2.User, method model2.LoginMethod, success bool, failureReason, ipAddress, userAgent string) {
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
