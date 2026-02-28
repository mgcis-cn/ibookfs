package email

import (
	"errors"

	"github.com/mgcis-cn/ibookfs/cmd/apiserver/app/options"
	"github.com/mgcis-cn/ibookfs/pkg/email/sources"
)

// Factory manages multiple named email instances.
type Factory struct {
	items map[string]sources.Email
}

func WithNameFunc(opts *options.EmailOptions) string {
	return opts.Name
}

// NewFactory creates a Factory from email options.
func NewFactory(configs []*options.EmailOptions, nameFunc func(*options.EmailOptions) string) (*Factory, error) {
	f := &Factory{items: make(map[string]sources.Email)}
	for _, c := range configs {
		cfg := c.ActiveConfig()
		if cfg == nil {
			continue
		}
		s, err := sources.Create(cfg)
		if err != nil {
			return nil, err
		}
		f.items[nameFunc(c)] = s
	}
	return f, nil
}

// Get returns email by name.
func (f *Factory) Get(name string) (sources.Email, bool) {
	v, ok := f.items[name]
	return v, ok
}

// MustGet returns email by name or error if not found.
func (f *Factory) MustGet(name string) (sources.Email, error) {
	v, ok := f.items[name]
	if !ok {
		return nil, errors.New("email not found: " + name)
	}
	return v, nil
}
