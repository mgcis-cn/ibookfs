// Package database provides database connection using atomic.Pointer pattern
// similar to Go's log/slog for thread-safe default database access.
package database

import (
	"log"
	"sync/atomic"

	"github.com/mgcis/ibookfs/internal/config"
	"github.com/mgcis/ibookfs/internal/database/sources"
	"github.com/mgcis/ibookfs/internal/database/sources/mysql"
	"github.com/mgcis/ibookfs/internal/database/sources/sqlite"
	"github.com/mgcis/ibookfs/internal/model"
	"gorm.io/gorm"
)

// defaultDB holds the default database instance using atomic.Pointer
// for thread-safe access, similar to log/slog's defaultLogger pattern.
var defaultDB atomic.Pointer[gorm.DB]

// Default returns the default database instance.
// If no database has been set, it initializes a SQLite database.
func Default() *gorm.DB {
	db := defaultDB.Load()
	if db == nil {
		// Initialize default SQLite database
		source := sqlite.New("ibookfs.db")
		if err := OpenSource(source); err != nil {
			log.Printf("Failed to initialize default SQLite database: %v", err)
			return nil
		}
		db = defaultDB.Load()
	}
	return db
}

// SetDefault sets the default database instance.
func SetDefault(db *gorm.DB) {
	defaultDB.Store(db)
}

// OpenSource opens a database connection from a Source and sets it as default.
func OpenSource(source sources.Source) error {
	db, err := source.Open()
	if err != nil {
		return err
	}
	log.Printf("Database connected: %s", source.Kind())
	SetDefault(db)
	return nil
}

// AutoMigrate runs auto migration for all models.
func AutoMigrate(db *gorm.DB) error {
	if err := db.AutoMigrate(&model.Book{}); err != nil {
		return err
	}
	log.Println("Database migration completed")
	return nil
}

// Init initializes the database based on config.
// Uses the database kind from config to select the appropriate source.
func Init(cfg *config.Config) error {
	var source sources.Source
	if cfg.IsMySQL() {
		source = mysql.New(cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Password, cfg.Database.Name)
	} else {
		path := cfg.Database.Path
		if path == "" {
			path = "ibookfs.db"
		}
		source = sqlite.New(path)
	}

	if err := OpenSource(source); err != nil {
		return err
	}

	return AutoMigrate(Default())
}
