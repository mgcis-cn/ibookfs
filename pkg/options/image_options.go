package options

import (
	"github.com/spf13/pflag"
	"k8s.io/component-base/cli/flag"
)

// ImageProcessingConfig holds image processing configuration.
type ImageProcessingConfig struct {
	BlurHashEnabled bool          `json:"blurhash_enabled" yaml:"blurhash_enabled"`
	Variants        []VariantSpec `json:"variants" yaml:"variants"`
}

func (o *ImageProcessingConfig) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.BoolVar(&o.BlurHashEnabled, FlagNameFunc(fs, "blurhash-enabled"), o.BlurHashEnabled, "Enable BlurHash generation for images.")
}

// VariantSpec defines a variant specification.
type VariantSpec struct {
	Name      string `json:"name" yaml:"name"`
	MaxWidth  int    `json:"max_width" yaml:"max_width"`
	MaxHeight int    `json:"max_height" yaml:"max_height"`
}

// ImageOptions holds image-related configuration.
type ImageOptions struct {
	Processing ImageProcessingConfig `json:"processing" yaml:"processing"`
}

func NewImageOptions() *ImageOptions {
	return &ImageOptions{}
}

func (o *ImageOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	var fss flag.NamedFlagSets
	o.Processing.AddFlags(fss.FlagSet(FlagNameFunc(fs, "processing")))
	for _, f := range fss.FlagSets {
		fs.AddFlagSet(f)
	}
}
