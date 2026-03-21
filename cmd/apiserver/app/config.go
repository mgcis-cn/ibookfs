// Package app provides Kratos-based configuration loading with support for
// multiple environments (dev, test, prod) using application.yaml + application-{env}.yaml pattern.
package app

import (
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	kratoslog "github.com/go-kratos/kratos/v2/log"
	"github.com/spf13/pflag"
)

var (
	flagConfDir string
	flagEnv     string
)

// AddConfigFlags registers config-related flags on the given pflag.FlagSet.
func AddConfigFlags(fs *pflag.FlagSet) {
	fs.StringVar(&flagConfDir, "conf", "configs", "config directory path")
	fs.StringVar(&flagEnv, "env", "", "runtime environment (e.g. dev, test, prod)")
}

// Load loads configuration with environment overlay.
// Priority: application-{env}.yaml > application.yaml
func Load() (config.Config, error) {
	baseConfig := filepath.Join(flagConfDir, "application.yaml")
	envConfig := filepath.Join(flagConfDir, "application-"+flagEnv+".yaml")

	// Initialize config files if base config does not exist
	if _, err := os.Stat(baseConfig); os.IsNotExist(err) {
		if err := initConfigFiles(flagConfDir); err != nil {
			return nil, err
		}
	}

	sources := []config.Source{
		file.NewSource(baseConfig),
	}

	// Add environment config if exists
	if _, err := os.Stat(envConfig); err == nil {
		sources = append(sources, file.NewSource(envConfig))
	}

	c := config.New(config.WithSource(sources...))
	if err := c.Load(); err != nil {
		return nil, err
	}

	return c, nil
}

// initConfigFiles creates the config directory and initializes default configuration files.
func initConfigFiles(confDir string) error {
	// Create config directory if not exists
	if err := os.MkdirAll(confDir, 0755); err != nil {
		return err
	}

	// Define config files to create
	configFiles := map[string]string{
		"application.yaml":      defaultApplicationConfig,
		"application-dev.yaml":  defaultEnvConfigTemplate("dev"),
		"application-test.yaml": defaultEnvConfigTemplate("test"),
		"application-prod.yaml": defaultEnvConfigTemplate("prod"),
	}

	logger := kratoslog.NewHelper(kratoslog.With(kratoslog.NewStdLogger(os.Stdout), "caller", kratoslog.DefaultCaller))

	for filename, content := range configFiles {
		filePath := filepath.Join(confDir, filename)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			return err
		}
		logger.Infof("Created config file: %s", filePath)
	}

	logger.Warn("========================================")
	logger.Warn("Configuration files have been initialized with default values.")
	logger.Warn("Please review and modify the configuration files before restarting:")
	logger.Warnf("  - %s/application.yaml (base configuration)", confDir)
	logger.Warnf("  - %s/application-{env}.yaml (environment-specific overrides)", confDir)
	logger.Warn("Important: Update JWT secret, database settings, and storage paths as needed.")
	logger.Warn("========================================")

	return nil
}

// defaultApplicationConfig contains the default base configuration.
const defaultApplicationConfig = `# iBookFS API Server Configuration
# This is the base configuration file. Environment-specific settings can override these values.

app:
  name: ibookfs
  version: 1.0.0

server:
  http:
    addr: ":8080"
    timeout: 2h
  image:
    processing:
      blurhash_enabled: true
      variants:
        - name: small
          max_width: 400
          max_height: 400
        - name: medium
          max_width: 1200
          max_height: 1200
  upload:
    max_file_size: 50Mi
    allowed_types:
      - image/jpeg
      - image/png
      - image/gif
      - image/webp
      - application/pdf
      - application/epub+zip

data:
  database:
    - name: default
      sqlite:
        source: "ibookfs.db"
  redis:
    addr: ""
    dial_timeout: 0s
    read_timeout: 0s
    write_timeout: 0s
  storage:
    - name: default
      local:
        endpoint: "http://localhost:8080"
        bucket: "storages"

auth:
  jwt:
    secret: "please-change-this-secret-key"
    expired: 2h
  oauth: []

email: []
`

// defaultEnvConfigTemplate returns environment-specific configuration template.
func defaultEnvConfigTemplate(env string) string {
	return `# iBookFS API Server - ` + env + ` Environment Configuration
# This file overrides values from application.yaml for the ` + env + ` environment.
# Only include settings that differ from the base configuration.

# Example overrides:
# server:
#   http:
#     addr: ":8080"
# data:
#   database:
#     - name: default
#       sqlite:
#         source: "ibookfs-` + env + `.db"
`
}
