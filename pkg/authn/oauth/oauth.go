package oauth

import (
	"errors"

	"github.com/mgcis-cn/ibookfs/cmd/apiserver/app/options"
	"github.com/mgcis-cn/ibookfs/pkg/authn/oauth/sources"
)

// Factory manages multiple named OAuth instances.
type Factory struct {
	items map[string]sources.OAuth
}

func WithNameFunc(opts *options.OAuthOptions) string {
	return opts.Name
}

// NewFactory creates a Factory from OAuth options.
func NewFactory(configs []*options.OAuthOptions, nameFunc func(*options.OAuthOptions) string) (*Factory, error) {
	f := &Factory{items: make(map[string]sources.OAuth)}
	for _, c := range configs {
		cfg := c.ActiveConfig()
		if cfg == nil {
			continue
		}
		o, err := sources.Create(cfg)
		if err != nil {
			return nil, err
		}
		f.items[nameFunc(c)] = o
	}
	return f, nil
}

// Get returns OAuth by name.
func (f *Factory) Get(name string) (sources.OAuth, bool) {
	v, ok := f.items[name]
	return v, ok
}

// MustGet returns OAuth by name or error if not found.
func (f *Factory) MustGet(name string) (sources.OAuth, error) {
	v, ok := f.items[name]
	if !ok {
		return nil, errors.New("oauth not found: " + name)
	}
	return v, nil
}

// GetProviderKinds returns a list of all configured OAuth provider kinds.
func (f *Factory) GetProviderKinds() []string {
	kinds := make([]string, 0, len(f.items))
	for name := range f.items {
		kinds = append(kinds, f.items[name].Kind())
	}
	return kinds
}
