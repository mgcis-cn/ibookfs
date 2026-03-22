package middleware

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/go-kratos/aegis/ratelimit"
	"github.com/go-kratos/kratos/v2/errors"
	"github.com/go-kratos/kratos/v2/middleware"
)

// mockLimiter implements ratelimit.Limiter for testing.
type mockLimiter struct {
	allowCount int64 // number of calls allowed before rejecting
	called     int64
}

func (m *mockLimiter) Allow() (ratelimit.DoneFunc, error) {
	n := atomic.AddInt64(&m.called, 1)
	if m.allowCount > 0 && n > m.allowCount {
		return nil, ratelimit.ErrLimitExceed
	}
	return func(ratelimit.DoneInfo) {}, nil
}

// rateLimitWithLimiter creates middleware with a custom limiter (for testing).
func rateLimitWithLimiter(limiter ratelimit.Limiter) middleware.Middleware {
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

func TestRateLimit_AllowsNormalRequests(t *testing.T) {
	mw := RateLimit(nil)
	handler := mw(func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	})

	resp, err := handler(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp != "ok" {
		t.Fatalf("expected 'ok', got %v", resp)
	}
}

func TestRateLimit_NilConfig(t *testing.T) {
	// Should not panic with nil config
	mw := RateLimit(nil)
	if mw == nil {
		t.Fatal("expected middleware, got nil")
	}

	handler := mw(func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	})
	_, err := handler(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error with nil config, got %v", err)
	}
}

func TestRateLimit_RejectsWhenLimited(t *testing.T) {
	limiter := &mockLimiter{allowCount: 1}
	mw := rateLimitWithLimiter(limiter)

	handler := mw(func(ctx context.Context, req interface{}) (interface{}, error) {
		return "ok", nil
	})

	// First request should pass
	resp, err := handler(context.Background(), nil)
	if err != nil {
		t.Fatalf("first request: expected no error, got %v", err)
	}
	if resp != "ok" {
		t.Fatalf("first request: expected 'ok', got %v", resp)
	}

	// Second request should be rejected
	_, err = handler(context.Background(), nil)
	if err == nil {
		t.Fatal("second request: expected error, got nil")
	}

	se := errors.FromError(err)
	if se.Code != 429 {
		t.Errorf("expected HTTP 429, got %d", se.Code)
	}
	if se.Reason != "RATE_LIMITED" {
		t.Errorf("expected reason RATE_LIMITED, got %s", se.Reason)
	}
}

func TestRateLimit_PropagatesHandlerError(t *testing.T) {
	limiter := &mockLimiter{allowCount: 0} // allow all
	mw := rateLimitWithLimiter(limiter)

	expectedErr := errors.New(500, "INTERNAL", "something failed")
	handler := mw(func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, expectedErr
	})

	_, err := handler(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	se := errors.FromError(err)
	if se.Code != 500 {
		t.Errorf("expected HTTP 500, got %d", se.Code)
	}
}

func TestRateLimit_ConcurrentRequests(t *testing.T) {
	// Allow first 10 requests, reject rest
	limiter := &mockLimiter{allowCount: 10}
	mw := rateLimitWithLimiter(limiter)

	var (
		successCount int64
		rejectCount  int64
		wg           sync.WaitGroup
	)

	handler := mw(func(ctx context.Context, req interface{}) (interface{}, error) {
		time.Sleep(time.Millisecond) // simulate work
		return "ok", nil
	})

	total := 20
	wg.Add(total)
	for i := 0; i < total; i++ {
		go func() {
			defer wg.Done()
			_, err := handler(context.Background(), nil)
			if err != nil {
				atomic.AddInt64(&rejectCount, 1)
			} else {
				atomic.AddInt64(&successCount, 1)
			}
		}()
	}
	wg.Wait()

	if successCount != 10 {
		t.Errorf("expected 10 successes, got %d", successCount)
	}
	if rejectCount != 10 {
		t.Errorf("expected 10 rejections, got %d", rejectCount)
	}
}

func TestRateLimit_DoneCalledWithError(t *testing.T) {
	var doneErr error
	var doneCalled bool

	customLimiter := &trackingLimiter{
		onDone: func(info ratelimit.DoneInfo) {
			doneCalled = true
			doneErr = info.Err
		},
	}

	mw := rateLimitWithLimiter(customLimiter)
	handlerErr := errors.New(400, "BAD_REQUEST", "bad")
	handler := mw(func(ctx context.Context, req interface{}) (interface{}, error) {
		return nil, handlerErr
	})

	handler(context.Background(), nil)

	if !doneCalled {
		t.Fatal("expected done to be called")
	}
	if doneErr != handlerErr {
		t.Errorf("expected done error to be handler error, got %v", doneErr)
	}
}

// trackingLimiter tracks done calls for testing.
type trackingLimiter struct {
	onDone func(ratelimit.DoneInfo)
}

func (l *trackingLimiter) Allow() (ratelimit.DoneFunc, error) {
	return func(info ratelimit.DoneInfo) {
		if l.onDone != nil {
			l.onDone(info)
		}
	}, nil
}
