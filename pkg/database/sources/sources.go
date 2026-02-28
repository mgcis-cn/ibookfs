// Package sources provides database source abstraction for multiple database types.
// Inspired by genai-toolbox's registry pattern for extensible data sources.
package sources

import (
	"context"
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Config is the interface that database configuration structs must implement.
type Config interface {
	Kind() string
}

// Factory creates a Source from configuration.
type Factory func(ctx context.Context, cfg Config) (Source, error)

// registry holds all registered source factories.
var registry = make(map[string]Factory)

// Register registers a source factory for a given kind.
// Returns false if the kind is already registered.
func Register(kind string, factory Factory) bool {
	if _, exists := registry[kind]; exists {
		return false
	}
	registry[kind] = factory
	return true
}

func DecodeConfig(ctx context.Context, cfg Config) (Source, error) {
	kind := cfg.Kind()
	factory, found := registry[kind]
	if !found {
		return nil, fmt.Errorf("unknown source kind: %q", kind)
	}
	source, err := factory(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("unable to parse source %q: %w", kind, err)
	}
	return source, err
}

// Source defines the interface for database sources.
type Source interface {
	Kind() string
	Open() (*gorm.DB, error)
}

// DefaultGormConfig returns the default GORM configuration.
func DefaultGormConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
}
