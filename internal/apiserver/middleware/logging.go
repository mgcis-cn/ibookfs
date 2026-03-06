// Package middleware provides Kratos middleware.
package middleware

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"

	"github.com/mgcis-cn/ibookfs/pkg/log"
)

// Logging creates request logging middleware with injected logger.
func Logging(logger log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			startTime := time.Now()

			// Get transport info
			var (
				operation string
				method    string
				path      string
				clientIP  string
			)

			if tr, ok := transport.FromServerContext(ctx); ok {
				operation = tr.Operation()
				if ht, ok := tr.(*http.Transport); ok {
					method = ht.Request().Method
					path = ht.Request().URL.Path
					clientIP = ht.Request().RemoteAddr
				}
			}

			// Process request
			resp, err := handler(ctx, req)

			// Calculate duration
			duration := time.Since(startTime)

			// Log based on result
			if err != nil {
				logger.Error("http_request",
					"method", method,
					"path", path,
					"operation", operation,
					"client_ip", clientIP,
					"duration_ms", duration.Milliseconds(),
					"error", err.Error(),
				)
			} else {
				logger.Info("http_request",
					"method", method,
					"path", path,
					"operation", operation,
					"client_ip", clientIP,
					"duration_ms", duration.Milliseconds(),
				)
			}

			return resp, err
		}
	}
}
