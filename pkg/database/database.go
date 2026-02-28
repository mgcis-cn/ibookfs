// Package database provides database connection with interface-based dependency injection.
package database

import (
	"context"
	"errors"
	"io"
	"log"

	"github.com/mgcis-cn/ibookfs/cmd/apiserver/app/options"
	"github.com/mgcis-cn/ibookfs/internal/apiserver/model"
	"github.com/mgcis-cn/ibookfs/pkg/database/sources"
	"gorm.io/gorm"
)

// Transactor defines transactional database operations.
type Transactor interface {
	TX(ctx context.Context, fn func(ctx context.Context) error) error
}

// DBConnector provides database connection access.
type DBConnector interface {
	Conn(ctx context.Context) *gorm.DB
}

// Database combines Transactor and DBConnector interfaces.
type Database interface {
	Transactor
	DBConnector
	io.Closer
}

// transactionKey is the context key for storing transactional DB.
type transactionKey struct{}

// datastore implements Database interface wrapping GORM.
type datastore struct {
	core *gorm.DB
}

// Close closes the database connection
func (d *datastore) Close() error {
	sqlDB, err := d.core.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// Conn retrieves the current transactional DB instance if it exists
// in context or falls back to the main database.
func (d *datastore) Conn(ctx context.Context) *gorm.DB {
	tx, ok := ctx.Value(transactionKey{}).(*gorm.DB)
	if ok {
		return tx
	}
	return d.core
}

// TX starts a transaction using the main DB context
// and passes the transactional context to the provided function.
func (d *datastore) TX(ctx context.Context, fn func(ctx context.Context) error) error {
	return d.core.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			ctx = context.WithValue(ctx, transactionKey{}, tx)
			return fn(ctx)
		},
	)
}

// New creates a new Database instance from a GORM database.
func New(db *gorm.DB) *datastore {
	return &datastore{core: db}
}

// AutoMigrate runs auto migration for all models.
func AutoMigrate(db Database) error {
	models := []interface{}{
		&model.Book{},
		&model.User{},
		&model.OAuthIdentity{},
		&model.EmailVerificationCode{},
		&model.Session{},
		&model.LoginHistory{},
		&model.Image{},
		&model.ImageVariant{},
		&model.ImageGroup{},
		&model.ImageGroupMember{},
	}

	for _, m := range models {
		if err := db.Conn(context.Background()).AutoMigrate(m); err != nil {
			return err
		}
	}
	log.Println("DatabaseOptions migration completed")
	return nil
}

type Factory struct {
	items map[string]Database
}

func WithNameFunc(opts *options.DatabaseOptions) string {
	return opts.Name
}

func NewFactory(ctx context.Context, configs []*options.DatabaseOptions, nameFunc func(dataOptions *options.DatabaseOptions) string) (*Factory, error) {
	f := &Factory{items: make(map[string]Database)}
	for _, c := range configs {
		cfg := c.ActiveConfig()
		if cfg == nil {
			continue
		}
		db, err := NewWith(ctx, cfg)
		if err != nil {
			return nil, err
		}
		f.items[nameFunc(c)] = db
	}
	return f, nil
}

func (f *Factory) Get(name string) (Database, bool) {
	v, ok := f.items[name]
	return v, ok
}

func (f *Factory) MustGet(name string) (Database, error) {
	v, ok := f.items[name]
	if !ok {
		return nil, errors.New("not found")
	}
	return v, nil
}

func (f *Factory) Close() (err error) {
	for _, db := range f.items {
		if db != nil {
			err = db.Close()
		}
	}
	return nil
}

// NewWith initializes the database based on config.
func NewWith(ctx context.Context, cfg sources.Config) (Database, error) {
	source, err := sources.DecodeConfig(ctx, cfg)
	if err != nil {
		return nil, err
	}
	db, err := source.Open()
	if err != nil {
		return nil, err
	}
	log.Printf("DatabaseOptions connected: %s", source.Kind())
	return New(db), nil
}
