package options

import (
	"github.com/spf13/pflag"
	"k8s.io/apimachinery/pkg/api/resource"
)

// UploadOptions holds upload configuration.
type UploadOptions struct {
	MaxFileSize  resource.Quantity `json:"max_file_size" yaml:"max_file_size" mapstructure:"max_file_size"`
	AllowedTypes []string          `json:"allowed_types" yaml:"allowed_types" mapstructure:"allowed_types"`
}

func (o *UploadOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.Var(&resource.QuantityValue{Quantity: o.MaxFileSize}, FlagNameFunc(fs, "max-file-size"), "Maximum upload file size (e.g. 50Mi).")
	fs.StringSliceVar(&o.AllowedTypes, FlagNameFunc(fs, "allowed-types"), o.AllowedTypes, "Allowed upload MIME types.")
}

func NewUploadOptions() *UploadOptions {
	return &UploadOptions{}
}
