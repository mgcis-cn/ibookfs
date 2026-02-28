package options

import "github.com/spf13/pflag"

type LocalOSOptions struct {
	Endpoint string `json:"endpoint" yaml:"endpoint"`
	Bucket   string `json:"bucket" yaml:"bucket"`
}

func NewLocalOSOptions() *LocalOSOptions {
	return &LocalOSOptions{}
}

func (o *LocalOSOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.StringVar(&o.Endpoint, FlagNameFunc(fs, "endpoint"), o.Endpoint, "Local storage endpoint URL.")
	fs.StringVar(&o.Bucket, FlagNameFunc(fs, "bucket"), o.Bucket, "Local storage bucket directory.")
}

func (o *LocalOSOptions) Kind() string {
	return "local"
}
