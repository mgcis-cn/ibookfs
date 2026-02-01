// Package service provides business logic.
package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/mgcis-cn/ibookfs/internal/database"
	"github.com/mgcis-cn/ibookfs/internal/model"
	"gorm.io/gorm"
)

// OAuthService handles OAuth operations.
type OAuthService struct {
	authSvc    *AuthService
	httpClient *http.Client
	configs    map[model.OAuthProvider]model.OAuthConfig
}

// NewOAuthService creates a new OAuth service.
func NewOAuthService(authSvc *AuthService, configs map[model.OAuthProvider]model.OAuthConfig) *OAuthService {
	return &OAuthService{
		authSvc:    authSvc,
		httpClient: &http.Client{Timeout: 30 * time.Second},
		configs:    configs,
	}
}

// GetAuthorizeURL returns the OAuth authorization URL.
func (s *OAuthService) GetAuthorizeURL(provider model.OAuthProvider, redirectURI, state string) (string, error) {
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
func (s *OAuthService) HandleCallback(provider model.OAuthProvider, code, state, redirectURI, ipAddress, userAgent string) (*AuthResponse, error) {
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
	user, isNew, err := s.findOrCreateUser(userInfo, ipAddress, userAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	// Update or create OAuth identity
	if err := s.updateOAuthIdentity(user.ID, userInfo, tokenResp, provider); err != nil {
		return nil, fmt.Errorf("failed to update OAuth identity: %w", err)
	}

	// Generate tokens
	tokenPair, err := s.authSvc.jwtManager.GenerateTokenPair(user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Load OAuth identities
	var oauthIdentities []model.OAuthIdentity
	database.Default().Where("user_id = ?", user.ID).Find(&oauthIdentities)

	return &AuthResponse{
		User:           user,
		Token:          tokenPair.AccessToken,
		RefreshToken:   tokenPair.RefreshToken,
		ExpiresIn:      int64(time.Until(tokenPair.ExpiresAt).Seconds()),
		LinkedAccounts: oauthIdentities,
		IsNewUser:      isNew,
	}, nil
}

// LinkOAuth links an OAuth account to existing user.
func (s *OAuthService) LinkOAuth(userID uint, provider model.OAuthProvider, code, redirectURI string) error {
	cfg, ok := s.configs[provider]
	if !ok {
		return errors.New("unsupported OAuth provider")
	}

	// Check if already linked
	var existingIdentity model.OAuthIdentity
	err := database.Default().Where("user_id = ? AND provider = ?", userID, provider).First(&existingIdentity).Error
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
	var otherUserIdentity model.OAuthIdentity
	err = database.Default().Where("provider = ? AND provider_user_id = ?", provider, userInfo.ID).First(&otherUserIdentity).Error
	if err == nil {
		return errors.New("该OAuth账号已绑定到其他用户")
	}

	// Create OAuth identity
	identity := &model.OAuthIdentity{
		UserID:           userID,
		Provider:         provider.String(),
		ProviderUserID:   userInfo.ID,
		ProviderUsername: &userInfo.Username,
		ProviderEmail:    &userInfo.Email,
		AccessToken:      &tokenResp.AccessToken,
		IsPrimary:        false,
	}

	if tokenResp.RefreshToken != "" {
		identity.RefreshToken = &tokenResp.RefreshToken
	}

	if !tokenResp.ExpiresAt.IsZero() {
		identity.TokenExpiresAt = &tokenResp.ExpiresAt
	}

	now := time.Now()
	identity.LinkedAt = &now

	return database.Default().Create(identity).Error
}

// Helper functions

func (s *OAuthService) exchangeCodeForToken(provider model.OAuthProvider, code, redirectURI string, cfg model.OAuthConfig) (*model.OAuthTokenResponse, error) {
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

	var tokenResp model.OAuthTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return nil, err
	}

	// Calculate ExpiresAt from ExpiresIn
	if tokenResp.ExpiresIn > 0 {
		tokenResp.ExpiresAt = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second)
	}

	return &tokenResp, nil
}

func (s *OAuthService) getUserInfo(provider model.OAuthProvider, accessToken string, cfg model.OAuthConfig) (*model.OAuthUserInfo, error) {
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

func (s *OAuthService) parseUserInfo(provider model.OAuthProvider, body io.Reader) (*model.OAuthUserInfo, error) {
	switch provider {
	case model.OAuthProviderGitHub:
		return s.parseGitHubUserInfo(body)
	case model.OAuthProviderGitee:
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

func (s *OAuthService) parseGitHubUserInfo(body io.Reader) (*model.OAuthUserInfo, error) {
	var githubUser GitHubUserinfo
	if err := json.NewDecoder(body).Decode(&githubUser); err != nil {
		return nil, err
	}

	// If email is not public, fetch from emails API
	if githubUser.Email == "" {
		githubUser.Email = s.getGitHubPrimaryEmail(githubUser.Login)
	}

	return &model.OAuthUserInfo{
		ID:        fmt.Sprintf("%d", githubUser.ID),
		Username:  githubUser.Login,
		Email:     githubUser.Email,
		Name:      githubUser.Name,
		AvatarURL: githubUser.AvatarURL,
		Provider:  model.OAuthProviderGitHub,
	}, nil
}

// GitHubEmail represents GitHub email response.
type GitHubEmail struct {
	Email   string `json:"email"`
	Primary bool   `json:"primary"`
}

func (s *OAuthService) getGitHubPrimaryEmail(username string) string {
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

func (s *OAuthService) parseGiteeUserInfo(body io.Reader) (*model.OAuthUserInfo, error) {
	var giteeUser GiteeUserinfo
	if err := json.NewDecoder(body).Decode(&giteeUser); err != nil {
		return nil, err
	}

	return &model.OAuthUserInfo{
		ID:        fmt.Sprintf("%d", giteeUser.ID),
		Username:  giteeUser.Login,
		Email:     giteeUser.Email,
		Name:      giteeUser.Name,
		AvatarURL: giteeUser.AvatarURL,
		Provider:  model.OAuthProviderGitee,
	}, nil
}

func (s *OAuthService) findOrCreateUser(userInfo *model.OAuthUserInfo, ipAddress, userAgent string) (*model.User, bool, error) {
	// First, try to find user by OAuth identity
	var identity model.OAuthIdentity
	err := database.Default().Where("provider = ? AND provider_user_id = ?", userInfo.Provider.String(), userInfo.ID).
		Preload("User").
		First(&identity).Error

	if err == nil {
		// Found existing user with this OAuth identity
		user := identity.User
		now := time.Now()
		user.LastLoginAt = &now
		user.LastLoginIP = &ipAddress
		database.Default().Save(user)

		// Update last_used_at
		identity.LastUsedAt = &now
		database.Default().Save(&identity)

		// Log login
		s.authSvc.logLoginHistory(user, userInfo.Provider.ToLoginMethod(), true, "", ipAddress, userAgent)

		return user, false, nil
	}

	// Not found by OAuth, check if email exists
	var user model.User
	if userInfo.Email != "" {
		err = database.Default().Where("email = ?", userInfo.Email).First(&user).Error
		if err == nil {
			// Found user by email, link this OAuth identity
			return s.linkOAuthToUser(&user, userInfo, ipAddress, userAgent)
		}
	}

	// Create new user
	return s.createUserFromOAuth(userInfo, ipAddress, userAgent)
}

func (s *OAuthService) linkOAuthToUser(user *model.User, userInfo *model.OAuthUserInfo, ipAddress, userAgent string) (*model.User, bool, error) {
	now := time.Now()
	user.LastLoginAt = &now
	user.LastLoginIP = &ipAddress
	if userInfo.AvatarURL != "" {
		user.AvatarURL = &userInfo.AvatarURL
	}
	database.Default().Save(user)

	// Create OAuth identity
	identity := &model.OAuthIdentity{
		UserID:           user.ID,
		Provider:         userInfo.Provider.String(),
		ProviderUserID:   userInfo.ID,
		ProviderUsername: &userInfo.Username,
		ProviderEmail:    &userInfo.Email,
		IsPrimary:        false,
		LinkedAt:         &now,
		LastUsedAt:       &now,
	}

	database.Default().Create(identity)

	// Log login
	s.authSvc.logLoginHistory(user, userInfo.Provider.ToLoginMethod(), true, "", ipAddress, userAgent)

	return user, false, nil
}

func (s *OAuthService) createUserFromOAuth(userInfo *model.OAuthUserInfo, ipAddress, userAgent string) (*model.User, bool, error) {
	// Generate display name
	displayName := userInfo.Name
	if displayName == "" {
		displayName = userInfo.Username
	}

	// Generate username/firstName/lastName
	// For OAuth users, we'll use the username as both first and last name initially
	user := &model.User{
		Email:       userInfo.Email,
		FirstName:   userInfo.Username,
		LastName:    userInfo.Username,
		DisplayName: &displayName,
		AvatarURL:   &userInfo.AvatarURL,
		Status:      model.UserStatusActive,
		Role:        model.UserRoleUser,
		Language:    "zh-CN",
		Timezone:    "Asia/Shanghai",
	}

	if err := database.Default().Create(user).Error; err != nil {
		return nil, false, err
	}

	// Create OAuth identity as primary
	now := time.Now()
	identity := &model.OAuthIdentity{
		UserID:           user.ID,
		Provider:         userInfo.Provider.String(),
		ProviderUserID:   userInfo.ID,
		ProviderUsername: &userInfo.Username,
		ProviderEmail:    &userInfo.Email,
		IsPrimary:        true,
		LinkedAt:         &now,
		LastUsedAt:       &now,
	}

	database.Default().Create(identity)

	// Log login
	s.authSvc.logLoginHistory(user, userInfo.Provider.ToLoginMethod(), true, "", ipAddress, userAgent)

	return user, true, nil
}

func (s *OAuthService) updateOAuthIdentity(userID uint, userInfo *model.OAuthUserInfo, tokenResp *model.OAuthTokenResponse, provider model.OAuthProvider) error {
	var identity model.OAuthIdentity
	err := database.Default().Where("user_id = ? AND provider = ?", userID, provider).First(&identity).Error

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
		return database.Default().Save(&identity).Error
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Create new identity
		now := time.Now()
		identity = model.OAuthIdentity{
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
			identity.RefreshToken = &tokenResp.RefreshToken
		}

		if !tokenResp.ExpiresAt.IsZero() {
			identity.TokenExpiresAt = &tokenResp.ExpiresAt
		}

		return database.Default().Create(&identity).Error
	}

	return err
}

// GetOAuthConfig returns OAuth config for a provider.
func GetOAuthConfig(provider model.OAuthProvider, clientID, clientSecret, redirectURI string) model.OAuthConfig {
	switch provider {
	case model.OAuthProviderGitHub:
		return model.OAuthConfig{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURI,
			AuthURL:      "https://github.com/login/oauth/authorize",
			TokenURL:     "https://github.com/login/oauth/access_token",
			UserInfoURL:  "https://api.github.com/user",
			Scopes:       []string{"read:user", "user:email"},
		}
	case model.OAuthProviderGitee:
		return model.OAuthConfig{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			RedirectURL:  redirectURI,
			AuthURL:      "https://gitee.com/oauth/authorize",
			TokenURL:     "https://gitee.com/oauth/token",
			UserInfoURL:  "https://gitee.com/api/v5/user",
			Scopes:       []string{"user_info", "emails"},
		}
	default:
		return model.OAuthConfig{}
	}
}
