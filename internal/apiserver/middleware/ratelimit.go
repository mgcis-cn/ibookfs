package middleware

import (
	"context"

	"github.com/go-kratos/aegis/ratelimit"
	"github.com/go-kratos/aegis/ratelimit/bbr"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

// RateLimitConfig holds rate limiter configuration.
type RateLimitConfig struct {
	// BBR options
	Opts []bbr.Option
}

// RateLimit creates a BBR adaptive rate limiting middleware.
// BBR adjusts limits based on CPU usage and request latency (Little's Law).
// When the system is overloaded, it returns HTTP 429 Too Many Requests.
func RateLimit(cfg *RateLimitConfig) middleware.Middleware {
	var opts []bbr.Option
	if cfg != nil {
		opts = cfg.Opts
	}
	limiter := bbr.NewLimiter(opts...)

	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			done, err := limiter.Allow()
			if err != nil {
				return nil, errors.New(429, "RATE_LIMITED", "服务过载，请稍后重试")
			}

			resp, err := handler(ctx, req)
			done(ratelimit.DoneInfo{Err: err})
			return resp, err
		}
	}
}
