// Package oauth provides OAuth business logic.
package oauth

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/mgcis-cn/ibookfs/internal/apiserver/biz/auth"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/model"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/store"
	"github.com/mgcis-cn/ibookfs/pkg/authn/oauth"
)

// OAuthBiz defines the interface for OAuth business logic.
type OAuthBiz interface {
	GetAuthorizeURL(ctx context.Context, provider model.OAuthProvider, redirectURI, state string) (string, error)
	HandleCallback(ctx context.Context, provider model.OAuthProvider, code, state, redirectURI, ipAddress, userAgent string) (*AuthResponse, error)
	LinkOAuth(ctx context.Context, userID uint, provider model.OAuthProvider, code, redirectURI string) error
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

// oauthBiz is the concrete implementation of OAuthBiz.
type oauthBiz struct {
	authSvc      auth.AuthBiz
	httpClient   *http.Client
	oauthFactory *oauth.Factory
	repo         store.IStore
}

// NewOAuthBiz creates a new OAuth business logic with injected dependencies.
func NewOAuthBiz(authSvc auth.AuthBiz, oauthFactory *oauth.Factory, repo store.IStore) OAuthBiz {
	return &oauthBiz{
		authSvc:      authSvc,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
		oauthFactory: oauthFactory,
		repo:         repo,
	}
}

// GetAuthorizeURL returns the OAuth authorization URL and stores state in database.
func (s *oauthBiz) GetAuthorizeURL(ctx context.Context, provider model.OAuthProvider, redirectURI, state string) (string, error) {
	oauth, err := s.oauthFactory.MustGet(string(provider))
	if err != nil {
		return "", fmt.Errorf("unsupported OAuth provider: %s", provider)
	}

	// Store state in database for validation during callback
	oauthState := &model.OAuthState{
		State:     state,
		Provider:  string(provider),
		Processed: false,
		ExpiresAt: time.Now().Add(10 * time.Minute), // State expires in 10 minutes
	}
	_ = s.repo.User().CreateOAuthState(ctx, oauthState)

	return oauth.GetAuthURL(state), nil
}

// HandleCallback handles OAuth callback and returns auth response.
func (s *oauthBiz) HandleCallback(ctx context.Context, provider model.OAuthProvider, code, state, redirectURI, ipAddress, userAgent string) (*AuthResponse, error) {
	// Check if state has already been processed (prevents duplicate callbacks)
	if s.isStateProcessed(ctx, state) {
		return nil, errors.New("this OAuth authorization has already been processed")
	}

	oauth, err := s.oauthFactory.MustGet(string(provider))
	if err != nil {
		return nil, fmt.Errorf("unsupported OAuth provider: %s", provider)
	}

	// Exchange code for token
	accessToken, err := oauth.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info from provider
	sourcesUserInfo, err := oauth.GetUserInfo(ctx, accessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}

	// Validate user info - must have valid ID and Email
	if sourcesUserInfo.ID == "" || sourcesUserInfo.ID == "0" {
		return nil, errors.New("invalid user info: missing provider user ID")
	}
	if sourcesUserInfo.Email == "" {
		return nil, errors.New("invalid user info: missing email (please make your email public or grant email permission)")
	}

	// Mark state as processed BEFORE creating user to prevent race conditions
	s.markStateProcessed(ctx, state)

	// Convert sources.UserInfo to model.OAuthUserInfo
	userInfo := &model.OAuthUserInfo{
		ID:        sourcesUserInfo.ID,
		Username:  sourcesUserInfo.Name,
		Email:     sourcesUserInfo.Email,
		Name:      sourcesUserInfo.Name,
		AvatarURL: sourcesUserInfo.Avatar,
		Provider:  provider,
	}

	// Find or create user (also handles OAuth identity creation/update)
	user, isNew, err := s.findOrCreateUser(ctx, userInfo, accessToken, ipAddress, userAgent)
	if err != nil {
		return nil, fmt.Errorf("failed to find or create user: %w", err)
	}

	// Generate tokens
	tokenPair, err := s.authSvc.GenerateTokenPair(ctx, user.ID, user.Email, string(user.Role))
	if err != nil {
		return nil, fmt.Errorf("failed to generate tokens: %w", err)
	}

	// Create session for token validation
	_ = s.authSvc.CreateSession(ctx, user, tokenPair.Token, ipAddress, userAgent)

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
func (s *oauthBiz) LinkOAuth(ctx context.Context, userID uint, provider model.OAuthProvider, code, redirectURI string) error {
	oauth, err := s.oauthFactory.MustGet(string(provider))
	if err != nil {
		return fmt.Errorf("unsupported OAuth provider: %s", provider)
	}

	userStore := s.repo.User()

	// Check if already linked
	_, err = userStore.GetOAuthIdentityByUserIDAndProvider(ctx, userID, provider.String())
	if err == nil {
		return errors.New("该OAuth账号已绑定")
	}

	// Exchange code for token
	accessToken, err := oauth.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to exchange code for token: %w", err)
	}

	// Get user info from provider
	sourcesUserInfo, err := oauth.GetUserInfo(ctx, accessToken)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}

	// Check if OAuth identity is linked to another user
	_, err = userStore.GetOAuthIdentityByProviderAndProviderUserID(ctx, provider.String(), sourcesUserInfo.ID)
	if err == nil {
		return errors.New("该OAuth账号已绑定到其他用户")
	}

	// Create OAuth identity
	now := time.Now()
	identity := &model.OAuthIdentity{
		UserID:           userID,
		Provider:         provider.String(),
		ProviderUserID:   sourcesUserInfo.ID,
		ProviderUsername: &sourcesUserInfo.Name,
		ProviderEmail:    &sourcesUserInfo.Email,
		AccessToken:      &accessToken,
		IsPrimary:        false,
		LinkedAt:         &now,
	}

	return userStore.CreateOAuthIdentity(ctx, identity)
}

// Helper functions

func (s *oauthBiz) findOrCreateUser(ctx context.Context, userInfo *model.OAuthUserInfo, accessToken, ipAddress, userAgent string) (*model.User, bool, error) {
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

		// Update OAuth identity with new access token
		identity.AccessToken = &accessToken
		identity.LastUsedAt = &now
		_ = userStore.UpdateOAuthIdentityModel(ctx, identity)

		// Log login
		s.logLoginHistory(ctx, user, userInfo.Provider.ToLoginMethod(), true, "", ipAddress, userAgent)

		return user, false, nil
	}

	// Not found by OAuth, check if email exists
	var user *model.User
	if userInfo.Email != "" {
		user, err = userStore.GetUserByEmail(ctx, userInfo.Email)
		if err == nil {
			// Found user by email, link this OAuth identity
			return s.linkOAuthToUser(ctx, user, userInfo, accessToken, ipAddress, userAgent)
		}
	}

	// Create new user
	return s.createUserFromOAuth(ctx, userInfo, accessToken, ipAddress, userAgent)
}

func (s *oauthBiz) linkOAuthToUser(ctx context.Context, user *model.User, userInfo *model.OAuthUserInfo, accessToken, ipAddress, userAgent string) (*model.User, bool, error) {
	userStore := s.repo.User()
	now := time.Now()

	user.LastLoginAt = &now
	user.LastLoginIP = &ipAddress
	if userInfo.AvatarURL != "" {
		user.AvatarURL = &userInfo.AvatarURL
	}
	_ = userStore.UpdateUserModel(ctx, user)

	// Create OAuth identity with access token
	identity := &model.OAuthIdentity{
		UserID:           user.ID,
		Provider:         userInfo.Provider.String(),
		ProviderUserID:   userInfo.ID,
		ProviderUsername: &userInfo.Username,
		ProviderEmail:    &userInfo.Email,
		AccessToken:      &accessToken,
		IsPrimary:        false,
		LinkedAt:         &now,
		LastUsedAt:       &now,
	}

	_ = userStore.CreateOAuthIdentity(ctx, identity)

	// Log login
	s.logLoginHistory(ctx, user, userInfo.Provider.ToLoginMethod(), true, "", ipAddress, userAgent)

	return user, false, nil
}

func (s *oauthBiz) createUserFromOAuth(ctx context.Context, userInfo *model.OAuthUserInfo, accessToken, ipAddress, userAgent string) (*model.User, bool, error) {
	userStore := s.repo.User()

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

	if err := userStore.CreateUser(ctx, user); err != nil {
		return nil, false, err
	}

	// Create OAuth identity as primary with access token
	now := time.Now()
	identity := &model.OAuthIdentity{
		UserID:           user.ID,
		Provider:         userInfo.Provider.String(),
		ProviderUserID:   userInfo.ID,
		ProviderUsername: &userInfo.Username,
		ProviderEmail:    &userInfo.Email,
		AccessToken:      &accessToken,
		IsPrimary:        true,
		LinkedAt:         &now,
		LastUsedAt:       &now,
	}

	_ = userStore.CreateOAuthIdentity(ctx, identity)

	// Log login
	s.logLoginHistory(ctx, user, userInfo.Provider.ToLoginMethod(), true, "", ipAddress, userAgent)

	return user, true, nil
}

// updateOAuthIdentityWithToken updates or creates OAuth identity with access token.
func (s *oauthBiz) updateOAuthIdentityWithToken(ctx context.Context, userID uint, userInfo *model.OAuthUserInfo, accessToken string, provider model.OAuthProvider) error {
	userStore := s.repo.User()

	identity, err := userStore.GetOAuthIdentityByUserIDAndProvider(ctx, userID, provider.String())
	if err == nil {
		// Update existing identity
		now := time.Now()
		identity.AccessToken = &accessToken
		identity.LastUsedAt = &now
		return userStore.UpdateOAuthIdentityModel(ctx, identity)
	}

	// Create new identity
	now := time.Now()
	newIdentity := &model.OAuthIdentity{
		UserID:           userID,
		Provider:         provider.String(),
		ProviderUserID:   userInfo.ID,
		ProviderUsername: &userInfo.Username,
		ProviderEmail:    &userInfo.Email,
		AccessToken:      &accessToken,
		IsPrimary:        true,
		LinkedAt:         &now,
		LastUsedAt:       &now,
	}

	return userStore.CreateOAuthIdentity(ctx, newIdentity)
}

func (s *oauthBiz) logLoginHistory(ctx context.Context, user *model.User, method model.LoginMethod, success bool, failureReason, ipAddress, userAgent string) {
	userStore := s.repo.User()
	history := &model.LoginHistory{
		UserID:        &user.ID,
		LoginMethod:   method,
		Success:       success,
		FailureReason: &failureReason,
		IPAddress:     &ipAddress,
		UserAgent:     &userAgent,
	}
	_ = userStore.CreateLoginHistory(ctx, history)
}

// isStateProcessed checks if an OAuth state has already been processed using database
func (s *oauthBiz) isStateProcessed(ctx context.Context, stateToken string) bool {
	userStore := s.repo.User()
	state, err := userStore.GetOAuthState(ctx, stateToken)
	if err != nil {
		return false // State not found, not processed
	}
	return state.Processed || time.Now().After(state.ExpiresAt)
}

// markStateProcessed marks an OAuth state as processed in database
func (s *oauthBiz) markStateProcessed(ctx context.Context, stateToken string) {
	userStore := s.repo.User()
	_ = userStore.MarkOAuthStateProcessed(ctx, stateToken)
}
