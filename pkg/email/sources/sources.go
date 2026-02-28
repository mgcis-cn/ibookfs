package sources

import (
	"context"
	"fmt"
	"log"
)

type Kind string

// Config is the interface that email configuration structs must implement.
type Config interface {
	Kind() string
}

type Factory func(cfg Config) (Email, error)

var registry = make(map[Kind]Factory)

func Register(kind Kind, factory Factory) bool {
	if _, exists := registry[kind]; exists {
		return false
	}
	registry[kind] = factory
	log.Printf("Registered email source: %s", kind)
	return true
}

func Create(cfg Config) (Email, error) {
	kind := Kind(cfg.Kind())
	factory, found := registry[kind]
	if !found {
		return nil, fmt.Errorf("unknown email kind: %q", kind)
	}
	return factory(cfg)
}

// Email defines the interface for email operations.
// Implementations can support different email providers (Netease, SendGrid, SES, etc.)
type Email interface {
	// Send sends an email with the given parameters.
	Send(ctx context.Context, to, subject, body string) error

	// IsEnabled checks if the email provider is properly configured.
	IsEnabled() bool

	// Kind returns the email provider type identifier.
	Kind() Kind
}
