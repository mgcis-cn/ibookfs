// Package sources provides database source abstraction for multiple database types.
// Inspired by genai-toolbox's registry pattern for extensible data sources.
package sources

import (
	"fmt"
	"log"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Kind represents the type of database source.
type Kind string

// Factory creates a Source from configuration.
type Factory func(cfg map[string]string) (Source, error)

// registry holds all registered source factories.
var (
	registry   = make(map[Kind]Factory)
	registryMu sync.RWMutex
)

// Register registers a source factory for a given kind.
// Returns false if the kind is already registered.
func Register(kind Kind, factory Factory) bool {
	registryMu.Lock()
	defer registryMu.Unlock()
	if _, exists := registry[kind]; exists {
		return false
	}
	registry[kind] = factory
	log.Printf("Registered database source: %s", kind)
	return true
}

// Create creates a source using the registered factory.
func Create(kind Kind, cfg map[string]string) (Source, error) {
	registryMu.RLock()
	factory, found := registry[kind]
	registryMu.RUnlock()
	if !found {
		return nil, fmt.Errorf("unknown source kind: %q", kind)
	}
	return factory(cfg)
}

// Source defines the interface for database sources.
type Source interface {
	Kind() Kind
	Open() (*gorm.DB, error)
	DSN() string
}

// DefaultGormConfig returns the default GORM configuration.
func DefaultGormConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}
}
