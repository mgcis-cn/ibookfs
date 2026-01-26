// Package config provides application configuration with support for
// YAML files, environment variables, and default values.
// Inspired by genai-toolbox's configuration pattern.
package config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
}

// ServerConfig holds server-related configuration.
type ServerConfig struct {
	Address        string   `yaml:"address"`
	AllowedOrigins []string `yaml:"allowedOrigins"`
}

// DatabaseConfig holds database-related configuration.
type DatabaseConfig struct {
	Kind     string `yaml:"kind"` // sqlite, mysql
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Name     string `yaml:"name"`
	Path     string `yaml:"path"` // for sqlite
}

// defaults returns a Config with default values.
func defaults() *Config {
	return &Config{
		Server: ServerConfig{
			Address:        ":8080",
			AllowedOrigins: []string{"http://localhost:3000", "http://localhost:5173"},
		},
		Database: DatabaseConfig{
			Kind: "sqlite",
			Path: "ibookfs.db",
			Port: "3306",
		},
	}
}

// Load reads configuration from YAML file (if exists), then overrides with environment variables.
func Load() *Config {
	cfg := defaults()

	// Try to load from config file
	if data, err := os.ReadFile("config.yaml"); err == nil {
		_ = yaml.Unmarshal(data, cfg)
	}

	// Override with environment variables
	cfg.applyEnv()

	return cfg
}

// LoadFromFile loads configuration from a specific YAML file.
func LoadFromFile(path string) (*Config, error) {
	cfg := defaults()

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	cfg.applyEnv()
	return cfg, nil
}

// applyEnv overrides config values with environment variables.
func (c *Config) applyEnv() {
	if v := os.Getenv("SERVER_ADDR"); v != "" {
		c.Server.Address = v
	}
	if v := os.Getenv("ALLOWED_ORIGINS"); v != "" {
		c.Server.AllowedOrigins = strings.Split(v, ",")
	}
	if v := os.Getenv("DB_KIND"); v != "" {
		c.Database.Kind = v
	}
	if v := os.Getenv("DB_HOST"); v != "" {
		c.Database.Host = v
	}
	if v := os.Getenv("DB_PORT"); v != "" {
		c.Database.Port = v
	}
	if v := os.Getenv("DB_USER"); v != "" {
		c.Database.User = v
	}
	if v := os.Getenv("DB_PASSWORD"); v != "" {
		c.Database.Password = v
	}
	if v := os.Getenv("DB_NAME"); v != "" {
		c.Database.Name = v
	}
	if v := os.Getenv("DB_PATH"); v != "" {
		c.Database.Path = v
	}
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.Server.Address == "" {
		return fmt.Errorf("server address is required")
	}

	switch c.Database.Kind {
	case "sqlite":
		if c.Database.Path == "" {
			return fmt.Errorf("database path is required for sqlite")
		}
	case "mysql":
		if c.Database.Host == "" {
			return fmt.Errorf("database host is required for mysql")
		}
		if c.Database.Name == "" {
			return fmt.Errorf("database name is required for mysql")
		}
	default:
		return fmt.Errorf("unsupported database kind: %s", c.Database.Kind)
	}

	return nil
}

// DSN returns the database connection string for MySQL.
func (c *Config) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.Database.User, c.Database.Password, c.Database.Host, c.Database.Port, c.Database.Name)
}

// IsSQLite returns true if using SQLite database.
func (c *Config) IsSQLite() bool {
	return c.Database.Kind == "sqlite" || c.Database.Kind == ""
}

// IsMySQL returns true if using MySQL database.
func (c *Config) IsMySQL() bool {
	return c.Database.Kind == "mysql"
}
