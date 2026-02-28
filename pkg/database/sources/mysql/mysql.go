// Package mysql provides MySQL database source implementation.
package mysql

import (
	"context"
	"fmt"
	"time"

	"github.com/mgcis-cn/ibookfs/pkg/database/sources"
	"github.com/mgcis-cn/ibookfs/pkg/options"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const SourceKind = "mysql"

func init() {
	sources.Register(SourceKind, func(ctx context.Context, cfg sources.Config) (sources.Source, error) {
		opts, ok := cfg.(*options.MysqlOptions)
		if !ok {
			return nil, fmt.Errorf("mysql: invalid config type %T", cfg)
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

// SourceOption is a functional option for configuring a MySQL source.
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

// source implements sources.Source for MySQL database.
type source struct {
	dsn             string
	maxIdleConns    int
	maxOpenConns    int
	connMaxLifetime time.Duration
	connMaxIdleTime time.Duration
}

// New creates a new MySQL source with functional options.
func New(dsn string, opts ...SourceOption) *source {
	s := &source{dsn: dsn}
	for _, fn := range opts {
		fn(s)
	}
	return s
}

func (s *source) Kind() string { return SourceKind }
func (s *source) Open() (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(s.dsn), sources.DefaultGormConfig())
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
