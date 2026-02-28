package apiserver

import (
	"github.com/mgcis-cn/ibookfs/cmd/apiserver/app/options"
)

// Config wraps ServerRunOptions for internal use.
type Config struct {
	*options.ServerRunOptions
}

// NewConfig creates a new Config from ServerRunOptions.
func NewConfig(opts *options.ServerRunOptions) *Config {
	return &Config{ServerRunOptions: opts}
}
