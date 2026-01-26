// Package sqlite provides SQLite database source implementation.
package sqlite

import (
	"github.com/mgcis/ibookfs/internal/database/sources"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const SourceKind sources.Kind = "sqlite"

func init() {
	sources.Register(SourceKind, func(cfg map[string]string) (sources.Source, error) {
		path := cfg["path"]
		if path == "" {
			path = "ibookfs.db"
		}
		return New(path), nil
	})
}

// Source implements sources.Source for SQLite database.
type Source struct {
	Path string
}

// New creates a new SQLite source.
func New(path string) *Source {
	if path == "" {
		path = "ibookfs.db"
	}
	return &Source{Path: path}
}

func (s *Source) Kind() sources.Kind { return SourceKind }
func (s *Source) DSN() string        { return s.Path }
func (s *Source) Open() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open(s.Path), sources.DefaultGormConfig())
}
