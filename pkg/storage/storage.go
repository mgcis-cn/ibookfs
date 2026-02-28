package storage

import (
	"errors"

	"github.com/mgcis-cn/ibookfs/cmd/apiserver/app/options"
	"github.com/mgcis-cn/ibookfs/pkg/storage/sources"
)

// Factory manages multiple named storage instances.
type Factory struct {
	items map[string]sources.Storage
}

// NewFactory creates a Factory from storage options.
func NewFactory(configs []*options.StorageOptions, nameFunc func(*options.StorageOptions) string) (*Factory, error) {
	f := &Factory{items: make(map[string]sources.Storage)}
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

// Get returns storage by name.
func (f *Factory) Get(name string) (sources.Storage, bool) {
	v, ok := f.items[name]
	return v, ok
}

// MustGet returns storage by name or error if not found.
func (f *Factory) MustGet(name string) (sources.Storage, error) {
	v, ok := f.items[name]
	if !ok {
		return nil, errors.New("storage not found: " + name)
	}
	return v, nil
}
