package sources

import (
	"context"
	"fmt"
	"log"
)

// Config is the interface that OAuth configuration structs must implement.
type Config interface {
	Kind() string
}

type Factory func(cfg Config) (OAuth, error)

var registry = make(map[string]Factory)

func Register(kind string, factory Factory) bool {
	if _, exists := registry[kind]; exists {
		return false
	}
	registry[kind] = factory
	log.Printf("Registered oauth source: %s", kind)
	return true
}

func Create(cfg Config) (OAuth, error) {
	factory, found := registry[cfg.Kind()]
	if !found {
		return nil, fmt.Errorf("unknown oauth kind: %q", cfg.Kind())
	}
	return factory(cfg)
}

// UserInfo represents the user information returned by OAuth providers.
type UserInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Avatar   string `json:"avatar"`
	Provider string `json:"provider"`
}

// OAuth defines the interface for OAuth authentication operations.
type OAuth interface {
	// GetAuthURL returns the OAuth authorization URL.
	GetAuthURL(state string) string

	// Exchange exchanges the authorization code for an access token.
	Exchange(ctx context.Context, code string) (string, error)

	// GetUserInfo retrieves user information using the access token.
	GetUserInfo(ctx context.Context, token string) (*UserInfo, error)

	// Kind returns the OAuth provider type identifier.
	Kind() string
}
