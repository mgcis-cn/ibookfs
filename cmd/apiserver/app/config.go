// Package app provides Kratos-based configuration loading with support for
// multiple environments (dev, test, prod) using application.yaml + application-{env}.yaml pattern.
package app

import (
	"flag"
	"os"
	"path/filepath"

	"github.com/go-kratos/kratos/v2/config"
	"github.com/go-kratos/kratos/v2/config/file"
	"github.com/spf13/pflag"
)

var (
	flagConfDir string
	flagEnv     string
)

func init() {
	flag.StringVar(&flagConfDir, "conf", "configs", "config directory path")
}

// AddConfigFlags registers config-related flags on the given pflag.FlagSet.
func AddConfigFlags(fs *pflag.FlagSet) {
	fs.StringVar(&flagEnv, "env", "", "runtime environment (e.g. dev, test, prod)")
}

// Load loads configuration with environment overlay.
// Priority: application-{env}.yaml > application.yaml
func Load() (config.Config, error) {
	baseConfig := filepath.Join(flagConfDir, "application.yaml")
	envConfig := filepath.Join(flagConfDir, "application-"+flagEnv+".yaml")

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
