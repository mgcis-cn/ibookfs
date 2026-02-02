// Package config provides application configuration with support for
// YAML files, environment variables, and default values.
// Inspired by genai-toolbox's configuration pattern.
package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config holds all configuration for the application.
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	JWT      JWTConfig      `yaml:"jwt"`
	OAuth    OAuthConfig    `yaml:"oauth"`
	Email    EmailConfig    `yaml:"email"`
	Image    ImageConfig    `yaml:"image"`
	Storage  StorageConfig  `yaml:"storage"`
}

// ServerConfig holds server-related configuration.
type ServerConfig struct {
	Address        string   `yaml:"address"`
	BaseURL        string   `yaml:"baseURL"` // Base URL for static assets (e.g., https://your-domain.com)
	AllowedOrigins []string `yaml:"allowedOrigins"`
	SkipAuthPaths  []string `yaml:"skipAuthPaths"`
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

// JWTConfig holds JWT-related configuration.
type JWTConfig struct {
	Secret     string `yaml:"secret"`
	Expiration int    `yaml:"expiration"` // hours
}

// OAuthConfig holds OAuth-related configuration.
type OAuthConfig struct {
	GitHub OAuthProviderConfig `yaml:"github"`
	Gitee  OAuthProviderConfig `yaml:"gitee"`
}

// OAuthProviderConfig holds configuration for a single OAuth provider.
type OAuthProviderConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	RedirectURL  string `yaml:"redirect_url"`
}

type EmailConfer interface {
	IsEmailEnabled() bool
	GetEmailConfig() *EmailConfig
}

// EmailConfig holds email-related configuration.
type EmailConfig struct {
	Kind     string `yaml:"kind"` // netease, smtp, disabled
	Host     string `yaml:"host"`
	Port     int    `yaml:"port"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	From     string `yaml:"from"`
	FromName string `yaml:"from_name"`
}

// ImageConfig holds image-related configuration.
type ImageConfig struct {
	Processing ImageProcessingConfig `yaml:"processing"`
	Upload     ImageUploadConfig     `yaml:"upload"`
}

// StorageConfig holds storage configuration (root level).
type StorageConfig struct {
	Default string            `yaml:"default"` // local, oss, s3, cos
	Local   LocalStorageConfig `yaml:"local"`
	OSS     OSSConfig         `yaml:"oss"`
	S3      S3Config          `yaml:"s3"`
	COS     COSConfig         `yaml:"cos"`
}

// LocalStorageConfig holds local storage configuration.
type LocalStorageConfig struct {
	BasePath string `yaml:"base_path"`
	BaseURL  string `yaml:"base_url"`
}

// OSSConfig holds Alibaba Cloud OSS configuration.
type OSSConfig struct {
	Endpoint        string `yaml:"endpoint"`
	AccessKeyID     string `yaml:"access_key_id"`
	AccessKeySecret string `yaml:"access_key_secret"`
	BucketName      string `yaml:"bucket_name"`
}

// S3Config holds AWS S3 configuration.
type S3Config struct {
	Region          string `yaml:"region"`
	AccessKeyID     string `yaml:"access_key_id"`
	SecretAccessKey string `yaml:"secret_access_key"`
	Bucket          string `yaml:"bucket"`
	Endpoint        string `yaml:"endpoint"` // Optional, for S3-compatible services
}

// COSConfig holds Tencent Cloud COS configuration.
type COSConfig struct {
	SecretID     string `yaml:"secret_id"`
	SecretKey    string `yaml:"secret_key"`
	BucketName   string `yaml:"bucket_name"`
	BucketRegion string `yaml:"bucket_region"`
}

// ImageProcessingConfig holds image processing configuration.
type ImageProcessingConfig struct {
	BlurHashEnabled bool           `yaml:"blurhash_enabled"`
	Variants        []VariantSpec  `yaml:"variants"`
}

// VariantSpec defines a variant specification.
type VariantSpec struct {
	Name      string `yaml:"name"`
	MaxWidth  int    `yaml:"max_width"`
	MaxHeight int    `yaml:"max_height"`
}

// ImageUploadConfig holds image upload configuration.
type ImageUploadConfig struct {
	MaxFileSize  int64    `yaml:"max_file_size"`  // bytes
	AllowedTypes []string `yaml:"allowed_types"`
}

// defaults returns a Config with default values.
func defaults() *Config {
	return &Config{
		Server: ServerConfig{
			Address:        ":8080",
			AllowedOrigins: []string{"http://localhost:3000", "http://localhost:5173"},
			SkipAuthPaths: []string{
				"/health",
				"/api/v1/auth/send-code",
				"/api/v1/auth/login",
				"/api/v1/auth/register",
				"/api/v1/auth/oauth/authorize",
				"/api/v1/auth/oauth/callback",
				"/api/v1/auth/refresh",
			},
		},
		Database: DatabaseConfig{
			Kind: "sqlite",
			Path: "ibookfs.db",
			Port: "3306",
		},
		JWT: JWTConfig{
			Secret:     "your-secret-key-change-in-production",
			Expiration: 24,
		},
		Storage: StorageConfig{
			Default: "local",
			Local: LocalStorageConfig{
				BasePath: "./storages",
				BaseURL:  "/storages",
			},
		},
		Image: ImageConfig{
			Processing: ImageProcessingConfig{
				BlurHashEnabled: true,
				Variants: []VariantSpec{
					{Name: "small", MaxWidth: 400, MaxHeight: 400},
					{Name: "medium", MaxWidth: 1200, MaxHeight: 1200},
				},
			},
			Upload: ImageUploadConfig{
				MaxFileSize: 50 * 1024 * 1024, // 50MB
				AllowedTypes: []string{
					"image/jpeg",
					"image/png",
					"image/gif",
					"image/webp",
				},
			},
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
	if v := os.Getenv("SKIP_AUTH_PATHS"); v != "" {
		c.Server.SkipAuthPaths = strings.Split(v, ",")
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
	// JWT
	if v := os.Getenv("JWT_SECRET"); v != "" {
		c.JWT.Secret = v
	}
	// OAuth GitHub
	if v := os.Getenv("GITHUB_CLIENT_ID"); v != "" {
		c.OAuth.GitHub.ClientID = v
	}
	if v := os.Getenv("GITHUB_CLIENT_SECRET"); v != "" {
		c.OAuth.GitHub.ClientSecret = v
	}
	if v := os.Getenv("GITHUB_REDIRECT_URL"); v != "" {
		c.OAuth.GitHub.RedirectURL = v
	}
	// OAuth Gitee
	if v := os.Getenv("GITEE_CLIENT_ID"); v != "" {
		c.OAuth.Gitee.ClientID = v
	}
	if v := os.Getenv("GITEE_CLIENT_SECRET"); v != "" {
		c.OAuth.Gitee.ClientSecret = v
	}
	if v := os.Getenv("GITEE_REDIRECT_URL"); v != "" {
		c.OAuth.Gitee.RedirectURL = v
	}
	// Email
	if v := os.Getenv("EMAIL_KIND"); v != "" {
		c.Email.Kind = v
	}
	if v := os.Getenv("EMAIL_HOST"); v != "" {
		c.Email.Host = v
	}
	if v := os.Getenv("EMAIL_PORT"); v != "" {
		port, _ := strconv.ParseInt(v, 10, 64)
		c.Email.Port = int(port)
	}
	if v := os.Getenv("EMAIL_USER"); v != "" {
		c.Email.User = v
	}
	if v := os.Getenv("EMAIL_PASSWORD"); v != "" {
		c.Email.Password = v
	}
	if v := os.Getenv("EMAIL_FROM"); v != "" {
		c.Email.From = v
	}
	if v := os.Getenv("EMAIL_FROM_NAME"); v != "" {
		c.Email.FromName = v
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

// IsEmailEnabled returns true if email sending is enabled.
func (c *EmailConfig) IsEmailEnabled() bool {
	return c.Kind != "" && c.Kind != "disabled"
}

// GetEmailConfig returns SMTP configuration based on kind.
func (c *EmailConfig) GetEmailConfig() *EmailConfig {
	return &EmailConfig{
		Kind:     c.Kind,
		Host:     c.Host,
		Port:     c.Port,
		User:     c.User,
		Password: c.Password,
		From:     c.From,
		FromName: c.FromName,
	}
}
