// Package mysql provides MySQL database source implementation.
package mysql

import (
	"fmt"

	"github.com/mgcis-cn/ibookfs/internal/database/sources"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const SourceKind sources.Kind = "mysql"

func init() {
	sources.Register(SourceKind, func(cfg map[string]string) (sources.Source, error) {
		port := cfg["port"]
		if port == "" {
			port = "3306"
		}
		return New(cfg["host"], port, cfg["user"], cfg["password"], cfg["database"]), nil
	})
}

// Source implements sources.Source for MySQL database.
type Source struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	Charset  string
}

// New creates a new MySQL source.
func New(host, port, user, password, database string) *Source {
	if port == "" {
		port = "3306"
	}
	return &Source{
		Host:     host,
		Port:     port,
		User:     user,
		Password: password,
		Database: database,
		Charset:  "utf8mb4",
	}
}

func (s *Source) Kind() sources.Kind { return SourceKind }
func (s *Source) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=True&loc=Local",
		s.User, s.Password, s.Host, s.Port, s.Database, s.Charset)
}
func (s *Source) Open() (*gorm.DB, error) {
	return gorm.Open(mysql.Open(s.DSN()), sources.DefaultGormConfig())
}
