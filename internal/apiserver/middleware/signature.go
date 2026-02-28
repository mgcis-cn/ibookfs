// Package middleware provides signature-based authentication middleware.
package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/mgcis-cn/ibookfs/internal/apiserver/biz/accountsecret"
	pkgctx "github.com/mgcis-cn/ibookfs/pkg/context"
)

// SignatureAuthConfig represents signature authentication configuration.
type SignatureAuthConfig struct {
	AccountSecretService accountsecret.AccountSecretBiz
	SkipAuthPaths        []string
}

// SignatureAuth creates signature-based authentication middleware.
// This middleware validates requests signed with AccountKey + SecretKey.
// Similar to cloud providers' signature authentication (e.g., Alibaba Cloud OSS).
func SignatureAuth(cfg *SignatureAuthConfig) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			// Get transport info
			tr, ok := transport.FromServerContext(ctx)
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

			// Get signature parameters from headers
			var accountKey, signature, expiresStr string
			if ht, ok := tr.(*http.Transport); ok {
				accountKey = ht.RequestHeader().Get("X-Account-Key")
				signature = ht.RequestHeader().Get("X-Signature")
				expiresStr = ht.RequestHeader().Get("X-Expires")
			}

			// If signature headers are present, use signature auth
			if accountKey != "" && signature != "" && expiresStr != "" {
				// Parse expires
				expires, err := strconv.ParseInt(expiresStr, 10, 64)
				if err != nil {
					return nil, errors.BadRequest("INVALID_EXPIRES", "Invalid expires timestamp")
				}

				// Check if expired
				if time.Now().Unix() > expires {
					return nil, errors.Unauthorized("SIGNATURE_EXPIRED", "Signature has expired")
				}

				// Get account secret by account key
				accountSecret, err := cfg.AccountSecretService.GetByAccountKey(ctx, accountKey)
				if err != nil {
					return nil, errors.Unauthorized("INVALID_CREDENTIALS", "Invalid account key")
				}

				// Check if secret is valid
				if !accountSecret.IsValid() {
					return nil, errors.Unauthorized("SECRET_DISABLED", "Account secret is disabled or expired")
				}

				// Get request info for signature verification
				var method, reqPath, query string
				if ht, ok := tr.(*http.Transport); ok {
					r := ht.Request()
					method = r.Method
					reqPath = r.URL.Path
					query = r.URL.RawQuery
				}

				// Verify signature
				expectedSignature := generateSignature(method, reqPath, query, accountSecret.SecretKey, expires)
				if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
					return nil, errors.Unauthorized("INVALID_SIGNATURE", "Signature verification failed")
				}

				// Update last used at time
				_ = cfg.AccountSecretService.UpdateLastUsedAt(ctx, accountSecret.ID)

				// Set signature context
				ctx = pkgctx.WithSignatureAuth(ctx, &pkgctx.SignatureContext{
					AccountSecret: accountSecret,
					AccountKey:    accountKey,
				})
				ctx = pkgctx.WithAuthType(ctx, "signature")
				ctx = pkgctx.WithOwnerID(ctx, accountSecret.OwnerID)

				return handler(ctx, req)
			}

			// No signature auth, proceed to next middleware (JWT may handle it)
			return handler(ctx, req)
		}
	}
}

// RequireSignature requires signature authentication (JWT is not accepted).
func RequireSignature() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			authType, ok := pkgctx.AuthType(ctx)
			if !ok || authType != "signature" {
				return nil, errors.Unauthorized("SIGNATURE_REQUIRED", "Signature authentication is required")
			}
			return handler(ctx, req)
		}
	}
}

// generateSignature generates HMAC-SHA256 signature for the request.
// Signature format: HMAC-SHA256(secret_key, string_to_sign)
// String to sign: method:path:query:expires
func generateSignature(method, path, query, secretKey string, expires int64) string {
	// Normalize query string (sort parameters)
	if query != "" {
		query = sortQueryString(query)
	}

	stringToSign := fmt.Sprintf("%s:%s:%s:%d", method, path, query, expires)

	// Calculate HMAC-SHA256
	h := hmac.New(sha256.New, []byte(secretKey))
	h.Write([]byte(stringToSign))
	sig := h.Sum(nil)

	return hex.EncodeToString(sig)
}

// sortQueryString sorts query parameters for consistent signature generation.
func sortQueryString(query string) string {
	params := strings.Split(query, "&")
	sort.Strings(params)
	return strings.Join(params, "&")
}
