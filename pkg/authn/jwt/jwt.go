package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	defaultSecret = "ibookfs!#888"
)

type options struct {
	tokenType     string
	expired       time.Duration
	signingMethod jwt.SigningMethod
	signingKey    any
	keyfunc       jwt.Keyfunc

	issuer      string
	tokenHeader map[string]any
}

type Option func(*options)

var defaultOptions = options{
	tokenType:     "Bearer",
	expired:       2 * time.Hour,
	signingMethod: jwt.SigningMethodHS256,
	signingKey:    []byte(defaultSecret),
	keyfunc: func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return []byte(defaultSecret), nil
	},
}

// WithSigningMethod set signature method.
func WithSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

// WithIssuer set token issuer which is identifies the principal that issued the JWTOptions.
func WithIssuer(issuer string) Option {
	return func(o *options) {
		o.issuer = issuer
	}
}

// WithSigningKey set the signature key.
func WithSigningKey(key any) Option {
	return func(o *options) {
		o.signingKey = key
	}
}

// WithKeyfunc set the callback function for verifying the key.
func WithKeyfunc(keyFunc jwt.Keyfunc) Option {
	return func(o *options) {
		o.keyfunc = keyFunc
	}
}

// WithExpired set the token expiration time (in seconds, default 2h).
func WithExpired(expired time.Duration) Option {
	return func(o *options) {
		o.expired = expired
	}
}

// WithTokenHeader set the customer tokenHeader for client side.
func WithTokenHeader(header map[string]any) Option {
	return func(o *options) {
		o.tokenHeader = header
	}
}

// Claims represents JWTOptions claims.
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// TokenPair represents access and refresh tokens.
type TokenPair struct {
	Token     string `json:"token"`
	Type      string `json:"type"`
	ExpiresAt int64  `json:"expires_at"`
}

type Auth struct {
	opts *options
}

// New create a authentication instance.
func New(opts ...Option) *Auth {
	o := defaultOptions
	for _, opt := range opts {
		opt(&o)
	}

	// Rebuild keyfunc to use the actual signingKey (not the hardcoded default)
	if o.keyfunc == nil || o.signingKey != nil {
		signingKey := o.signingKey
		o.keyfunc = func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrTokenSignatureInvalid
			}
			return signingKey, nil
		}
	}

	return &Auth{opts: &o}
}

// Sign generates access and refresh tokens.
func (a *Auth) Sign(userID string, email string, role string) (*TokenPair, error) {
	now := time.Now()
	expiresAt := now.Add(a.opts.expired)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    a.opts.issuer,
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Subject:   userID,
		},
	})

	if a.opts.tokenHeader != nil {
		for k, v := range a.opts.tokenHeader {
			token.Header[k] = v
		}
	}
	accessToken, err := token.SignedString(a.opts.signingKey)
	if err != nil {
		return nil, fmt.Errorf("failed to sign access token: %w", err)
	}

	return &TokenPair{
		Token:     accessToken,
		Type:      a.opts.tokenType,
		ExpiresAt: expiresAt.Unix(),
	}, nil
}

// ParseClaims parse the token and return the claims.
func (a *Auth) ParseClaims(ctx context.Context, tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, a.opts.keyfunc)
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if token.Method != a.opts.signingMethod {
		return nil, fmt.Errorf("invalid token")
	}

	return token.Claims.(*Claims), nil
}
