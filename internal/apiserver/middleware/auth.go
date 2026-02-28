// Package middleware provides Kratos middleware.
package middleware

import (
	"context"
	"strconv"
	"strings"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/mgcis-cn/ibookfs/internal/apiserver/model"
	"github.com/mgcis-cn/ibookfs/pkg/authn/jwt"
	pkgctx "github.com/mgcis-cn/ibookfs/pkg/context"
)

// AuthConfig represents authentication configuration.
type AuthConfig struct {
	SkipAuthPaths []string
	JWTManager    *jwt.Auth
}

// Auth creates JWT authentication middleware.
func Auth(cfg *AuthConfig) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// Get transport info
			tr, ok := transport.FromClientContext(ctx)
			if !ok {
				return handler(ctx, req)
			}

			// Check if path should skip auth
			path := tr.Operation()
			for _, skipPath := range cfg.SkipAuthPaths {
				if path == skipPath {
					return handler(ctx, req)
				}
			}

			// Get Authorization header
			var authHeader string
			if ht, ok := tr.(*http.Transport); ok {
				authHeader = ht.RequestHeader().Get("Authorization")
			}

			if authHeader == "" {
				return nil, errors.Unauthorized("UNAUTHORIZED", "Missing authorization header")
			}

			// Extract Bearer token
			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return nil, errors.Unauthorized("INVALID_TOKEN_FORMAT", "Invalid authorization header format")
			}

			tokenString := parts[1]

			// Validate token
			claims, err := cfg.JWTManager.ParseClaims(ctx, tokenString)
			if err != nil {
				return nil, errors.Unauthorized("TOKEN_EXPIRED", "Invalid or expired token")
			}

			// Set user info in context
			userID, _ := strconv.Atoi(claims.UserID)
			ctx = pkgctx.WithUserId(ctx, userID)
			ctx = pkgctx.WithUserEmail(ctx, claims.Email)
			ctx = pkgctx.WithUserRole(ctx, claims.Role)

			return handler(ctx, req)
		}
	}
}

// RequireRole checks if user has required role.
func RequireRole(roles ...model.UserRole) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			userRole, ok := pkgctx.UserRole(ctx)
			if !ok {
				return nil, errors.Unauthorized("UNAUTHORIZED", "User not authenticated")
			}

			// Check if user has required role
			hasRole := false
			for _, role := range roles {
				if string(role) == userRole {
					hasRole = true
					break
				}
			}

			if !hasRole {
				return nil, errors.Forbidden("FORBIDDEN", "Insufficient permissions")
			}

			return handler(ctx, req)
		}
	}
}
