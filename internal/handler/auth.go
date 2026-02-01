// Package handler contains HTTP request handlers.
package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/mgcis-cn/ibookfs/internal/config"
	"github.com/mgcis-cn/ibookfs/internal/middleware"
	"github.com/mgcis-cn/ibookfs/internal/model"
	"github.com/mgcis-cn/ibookfs/internal/service"
)

// AuthHandler handles authentication HTTP requests.
type AuthHandler struct {
	authSvc     *service.AuthService
	oauthSvc    *service.OAuthService
	oauthConfig map[model.OAuthProvider]model.OAuthConfig
	emailConfig *config.EmailConfig
}

// NewAuthHandler creates a new auth handler.
func NewAuthHandler(
	authSvc *service.AuthService,
	oauthSvc *service.OAuthService,
	oauthConfig map[model.OAuthProvider]model.OAuthConfig,
	emailConfig *config.EmailConfig,
) *AuthHandler {
	return &AuthHandler{
		authSvc:     authSvc,
		oauthSvc:    oauthSvc,
		oauthConfig: oauthConfig,
		emailConfig: emailConfig,
	}
}

// APIResponse represents a standard API response.
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
}

// SendCode handles POST /auth/send-code
func (h *AuthHandler) SendCode(c *gin.Context) {
	var req service.SendCodeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	// Validate type
	validTypes := map[string]bool{
		"login":          true,
		"register":       true,
		"reset_password": true,
		"bind_email":     true,
	}
	if !validTypes[req.Type] {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid code type",
		})
		return
	}

	if err := h.authSvc.SendCode(req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{
			"expires_in": 300,
			"message":    "验证码已发送",
		},
	})
}

// Login handles POST /auth/login
func (h *AuthHandler) Login(c *gin.Context) {
	var req service.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	resp, err := h.authSvc.Login(req, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    resp,
	})
}

// Register handles POST /auth/register
func (h *AuthHandler) Register(c *gin.Context) {
	var req service.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")

	resp, err := h.authSvc.Register(req, ipAddress, userAgent)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Success: true,
		Data:    resp,
	})
}

// OAuthAuthorize handles GET /auth/oauth/authorize
func (h *AuthHandler) OAuthAuthorize(c *gin.Context) {
	provider := c.Query("provider")
	redirectURI := c.Query("redirect_uri")
	state := c.Query("state")

	if provider == "" {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Provider is required",
		})
		return
	}

	if redirectURI == "" {
		redirectURI = c.GetHeader("Origin") + "/auth/callback"
	}

	if state == "" {
		state = "state_" + time.Now().Format("20060102150405")
	}

	authorizeURL, err := h.oauthSvc.GetAuthorizeURL(model.OAuthProvider(provider), redirectURI, state)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{
			"authorize_url": authorizeURL,
			"state":         state,
		},
	})
}

// OAuthCallback handles POST /auth/oauth/callback
func (h *AuthHandler) OAuthCallback(c *gin.Context) {
	var req struct {
		Provider string `json:"provider" binding:"required"`
		Code     string `json:"code" binding:"required"`
		State    string `json:"state" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	ipAddress := c.ClientIP()
	userAgent := c.GetHeader("User-Agent")
	redirectURI := c.GetHeader("Origin") + "/auth/callback"

	resp, err := h.oauthSvc.HandleCallback(
		model.OAuthProvider(req.Provider),
		req.Code,
		req.State,
		redirectURI,
		ipAddress,
		userAgent,
	)

	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data:    resp,
	})
}

// RefreshToken handles POST /auth/refresh
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	tokenPair, err := h.authSvc.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{
			"token":         tokenPair.AccessToken,
			"refresh_token": tokenPair.RefreshToken,
			"expires_in":    int64(time.Until(tokenPair.ExpiresAt).Seconds()),
		},
	})
}

// Logout handles POST /auth/logout
func (h *AuthHandler) Logout(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Not authenticated",
		})
		return
	}

	token := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

	if err := h.authSvc.Logout(userID, token); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to logout",
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{
			"message": "Logged out successfully",
		},
	})
}

// GetCurrentUser handles GET /auth/me
func (h *AuthHandler) GetCurrentUser(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Not authenticated",
		})
		return
	}

	user, err := h.authSvc.GetCurrentUser(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, APIResponse{
			Success: false,
			Message: "User not found",
		})
		return
	}

	// Load linked accounts
	linkedAccounts, _ := h.authSvc.GetLinkedAccounts(userID)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{
			"user":            user,
			"linked_accounts": linkedAccounts,
		},
	})
}

// GetLinkedAccounts handles GET /auth/accounts
func (h *AuthHandler) GetLinkedAccounts(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Not authenticated",
		})
		return
	}

	accounts, err := h.authSvc.GetLinkedAccounts(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to get linked accounts",
		})
		return
	}

	// Add email account if user has email
	user, _ := h.authSvc.GetCurrentUser(userID)
	result := []gin.H{}

	if user.Email != "" {
		result = append(result, gin.H{
			"provider":       "email",
			"provider_email": user.Email,
			"is_primary":     true,
			"linked_at":      user.CreatedAt,
		})
	}

	for _, acc := range accounts {
		result = append(result, gin.H{
			"provider":          acc.Provider,
			"provider_user_id":  &acc.ProviderUserID,
			"provider_username": acc.ProviderUsername,
			"provider_email":    acc.ProviderEmail,
			"is_primary":        acc.IsPrimary,
			"linked_at":         acc.LinkedAt,
			"last_used_at":      acc.LastUsedAt,
		})
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{
			"accounts": result,
		},
	})
}

// LinkOAuth handles POST /auth/oauth/link
func (h *AuthHandler) LinkOAuth(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Not authenticated",
		})
		return
	}

	var req struct {
		Provider string `json:"provider" binding:"required"`
		Code     string `json:"code" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	redirectURI := c.GetHeader("Origin") + "/settings/accounts/link"

	err := h.oauthSvc.LinkOAuth(userID, model.OAuthProvider(req.Provider), req.Code, redirectURI)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	// Get updated identity list
	c.Request.URL.Path = "/auth/accounts"
	h.GetLinkedAccounts(c)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{
			"message": "OAuth account linked successfully",
		},
	})
}

// UnlinkOAuth handles DELETE /auth/oauth/unlink
func (h *AuthHandler) UnlinkOAuth(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Not authenticated",
		})
		return
	}

	var req struct {
		Provider       string  `json:"provider" binding:"required"`
		ProviderUserID *string `json:"provider_user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	err := h.authSvc.UnlinkOAuth(userID, req.Provider, req.ProviderUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{
			"message": "OAuth account unlinked successfully",
		},
	})
}

// BindEmail handles POST /auth/bind-email
func (h *AuthHandler) BindEmail(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Not authenticated",
		})
		return
	}

	var req struct {
		Email string `json:"email" binding:"required,email"`
		Code  string `json:"code" binding:"required,len=6"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	err := h.authSvc.BindEmail(userID, req.Email, req.Code)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	user, _ := h.authSvc.GetCurrentUser(userID)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{
			"user": user,
		},
	})
}

// SetPrimaryAccount handles PUT /auth/accounts/primary
func (h *AuthHandler) SetPrimaryAccount(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Not authenticated",
		})
		return
	}

	var req struct {
		Provider       string  `json:"provider" binding:"required"`
		ProviderUserID *string `json:"provider_user_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

	err := h.authSvc.SetPrimaryAccount(userID, req.Provider, req.ProviderUserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{
			"message": "Primary account updated",
		},
	})
}

// UpdateProfile handles PUT /auth/me
func (h *AuthHandler) UpdateProfile(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, APIResponse{
			Success: false,
			Message: "Not authenticated",
		})
		return
	}

	var req struct {
		FirstName   *string `json:"first_name"`
		LastName    *string `json:"last_name"`
		DisplayName *string `json:"display_name"`
		Language    *string `json:"language"`
		Timezone    *string `json:"timezone"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Success: false,
			Message: "Invalid request format",
		})
		return
	}

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

	if err := h.authSvc.UpdateProfile(userID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Success: false,
			Message: "Failed to update profile",
		})
		return
	}

	user, _ := h.authSvc.GetCurrentUser(userID)

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{
			"user": user,
		},
	})
}

// GetConfig handles GET /auth/config
// Returns the authentication configuration including supported OAuth providers and email status
func (h *AuthHandler) GetConfig(c *gin.Context) {
	// Build list of supported OAuth providers
	var supportedProviders []string
	for provider, cfg := range h.oauthConfig {
		if cfg.ClientID != "" && cfg.ClientSecret != "" {
			supportedProviders = append(supportedProviders, string(provider))
		}
	}

	// Check if email is enabled
	emailEnabled := h.emailConfig.IsEmailEnabled()

	// Helper function to safely check if provider is configured
	isConfigured := func(provider model.OAuthProvider) bool {
		cfg, exists := h.oauthConfig[provider]
		return exists && cfg.ClientID != "" && cfg.ClientSecret != ""
	}

	// Check if any OAuth provider is enabled
	oauthEnabled := len(supportedProviders) > 0

	c.JSON(http.StatusOK, APIResponse{
		Success: true,
		Data: gin.H{
			"oauth_providers": supportedProviders,
			"email_enabled":   emailEnabled,
			"oauth_enabled":   oauthEnabled,
			"oauth": gin.H{
				"github": gin.H{"configured": isConfigured(model.OAuthProviderGitHub)},
				"gitee":  gin.H{"configured": isConfigured(model.OAuthProviderGitee)},
				"icloud": gin.H{"configured": isConfigured(model.OAuthProviderICloud)},
				"google": gin.H{"configured": isConfigured(model.OAuthProviderGoogle)},
				"wechat": gin.H{"configured": isConfigured(model.OAuthProviderWechat)},
			},
		},
	})
}

// RegisterAuthRoutes registers all auth routes.
func RegisterAuthRoutes(r *gin.RouterGroup, authHandler *AuthHandler, authMiddleware gin.HandlerFunc) {
	auth := r.Group("/auth")
	{
		// Public routes
		auth.GET("/config", authHandler.GetConfig)
		auth.POST("/send-code", authHandler.SendCode)
		auth.POST("/login", authHandler.Login)
		auth.POST("/register", authHandler.Register)
		auth.GET("/oauth/authorize", authHandler.OAuthAuthorize)
		auth.POST("/oauth/callback", authHandler.OAuthCallback)
		auth.POST("/refresh", authHandler.RefreshToken)

		// Protected routes
		protected := auth.Group("")
		protected.Use(authMiddleware)
		{
			protected.POST("/logout", authHandler.Logout)
			protected.GET("/me", authHandler.GetCurrentUser)
			protected.GET("/accounts", authHandler.GetLinkedAccounts)
			protected.POST("/oauth/link", authHandler.LinkOAuth)
			protected.DELETE("/oauth/unlink", authHandler.UnlinkOAuth)
			protected.POST("/bind-email", authHandler.BindEmail)
			protected.PUT("/accounts/primary", authHandler.SetPrimaryAccount)
			protected.PUT("/me", authHandler.UpdateProfile)
		}
	}
}
