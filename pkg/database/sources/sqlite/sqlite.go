// Package sqlite provides SQLite database source implementation.
package sqlite

import (
	"context"
	"fmt"
	"time"

	"github.com/mgcis-cn/ibookfs/pkg/database/sources"
	"github.com/mgcis-cn/ibookfs/pkg/options"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const SourceKind string = "sqlite"

func init() {
	sources.Register(SourceKind, func(ctx context.Context, cfg sources.Config) (sources.Source, error) {
		opts, ok := cfg.(*options.SqliteOptions)
		if !ok {
			return nil, fmt.Errorf("sqlite: invalid config type %T", cfg)
		}
		path := opts.Source
		if path == "" {
			path = "sqlite.db"
		}
		return New(path), nil
	})
}

// SourceOption is a functional option for configuring a SQLite source.
type SourceOption func(*source)

func WithMaxIdleConns(n int) SourceOption {
	return func(s *source) { s.maxIdleConns = n }
}

func WithMaxOpenConns(n int) SourceOption {
	return func(s *source) { s.maxOpenConns = n }
}

func WithConnMaxLifetime(d time.Duration) SourceOption {
	return func(s *source) { s.connMaxLifetime = d }
}

func WithConnMaxIdleTime(d time.Duration) SourceOption {
	return func(s *source) { s.connMaxIdleTime = d }
}

// source implements sources.Source for SQLite database.
type source struct {
	path            string
	maxIdleConns    int
	maxOpenConns    int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
}

// New creates a new SQLite source with functional options.
func New(path string, opts ...SourceOption) *source {
	if path == "" {
		path = "sqlite.db"
	}
	s := &source{path: path}
	for _, fn := range opts {
		fn(s)
	}
	return s
}

func (s *source) Kind() string { return SourceKind }
func (s *source) Open() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(s.path), sources.DefaultGormConfig())
	if err != nil {
		return nil, err
	}
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	if s.maxIdleConns > 0 {
		sqlDB.SetMaxIdleConns(s.maxIdleConns)
	}
	if s.maxOpenConns > 0 {
		sqlDB.SetMaxOpenConns(s.maxOpenConns)
	}
	if s.connMaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(s.connMaxLifetime)
	}
	if s.connMaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(s.connMaxIdleTime)
	}
	return db, nil
}
