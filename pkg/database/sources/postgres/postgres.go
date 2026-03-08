// Package postgres provides PostgreSQL database source implementation.
package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/mgcis-cn/ibookfs/pkg/database/sources"
	"github.com/mgcis-cn/ibookfs/pkg/options"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const SourceKind = "postgres"

func init() {
	sources.Register(SourceKind, func(ctx context.Context, cfg sources.Config) (sources.Source, error) {
		opts, ok := cfg.(*options.PostgresOptions)
		if !ok {
			return nil, fmt.Errorf("postgres: invalid config type %T", cfg)
		}
		return New(
			opts.Source,
			WithMaxIdleConns(opts.MaxIdleConns),
			WithMaxOpenConns(opts.MaxOpenConns),
			WithConnMaxLifetime(opts.ConnMaxLifetime.Duration),
			WithConnMaxIdleTime(opts.ConnMaxIdleTime.Duration),
		), nil
	})
}

// SourceOption is a functional option for configuring a PostgreSQL source.
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

// source implements sources.Source for PostgreSQL database.
type source struct {
	dsn             string
	maxIdleConns    int
	maxOpenConns    int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
}

// New creates a new PostgreSQL source with functional options.
func New(dsn string, opts ...SourceOption) *source {
	s := &source{dsn: dsn}
	for _, fn := range opts {
		fn(s)
	}
	return s
}

func (s *source) Kind() string { return SourceKind }
func (s *source) Open() (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(s.dsn), sources.DefaultGormConfig())
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
