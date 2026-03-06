package middleware

import (
	"github.com/go-kratos/kratos/v2/middleware"

	"github.com/mgcis-cn/ibookfs/internal/apiserver/biz/accountsecret"
	"github.com/mgcis-cn/ibookfs/pkg/authn/jwt"
	"github.com/mgcis-cn/ibookfs/pkg/log"
)

// Config holds middleware configuration.
type Config struct {
	SkipAuthPaths        []string
	JWTManager           *jwt.Auth
	AccountSecretService accountsecret.AccountSecretBiz
	Logger               log.Logger
}

// NewMiddlewares creates middleware chain for Kratos server.
func NewMiddlewares(cfg *Config) []middleware.Middleware {
	var middlewares []middleware.Middleware

	// Logging middleware (first in chain to capture all requests)
	if cfg.Logger != nil {
		middlewares = append(middlewares, Logging(cfg.Logger))
	}

	// Signature authentication (optional, checks headers first)
	if cfg.AccountSecretService != nil {
		middlewares = append(middlewares, SignatureAuth(&SignatureAuthConfig{
			AccountSecretService: cfg.AccountSecretService,
			SkipAuthPaths:        cfg.SkipAuthPaths,
		}))
	}

	// JWT authentication
	if cfg.JWTManager != nil {
		middlewares = append(middlewares, Auth(&AuthConfig{
			JWTManager:    cfg.JWTManager,
			SkipAuthPaths: cfg.SkipAuthPaths,
		}))
	}

	return middlewares
}
