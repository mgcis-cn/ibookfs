package options

import (
	"time"

	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type HttpOptions struct {
	Addr    string          `json:"addr" yaml:"addr" mapstructure:"addr"`
	Timeout metav1.Duration `json:"timeout" yaml:"timeout" mapstructure:"timeout"`
}

func (o *HttpOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.StringVar(&o.Addr, FlagNameFunc(fs, "addr"), o.Addr, "HTTP server listen address.")
	fs.DurationVar(&o.Timeout.Duration, FlagNameFunc(fs, "timeout"), o.Timeout.Duration, "HTTP server request timeout.")
}

func NewHttpOptions() *HttpOptions {
	return &HttpOptions{
		Addr:    ":8080",
		Timeout: metav1.Duration{Duration: 2 * time.Hour},
	}
}
