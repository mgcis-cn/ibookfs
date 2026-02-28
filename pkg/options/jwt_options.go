package options

import (
	"time"

	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type JWTOptions struct {
	Secret  string          `json:"secret" yaml:"secret" mapstructure:"secret"`
	Expired metav1.Duration `json:"expired" yaml:"expired" mapstructure:"expired"`
}

func (o *JWTOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.StringVar(&o.Secret, FlagNameFunc(fs, "secret"), o.Secret, "JWT signing secret key.")
	fs.DurationVar(&o.Expired.Duration, FlagNameFunc(fs, "expired"), o.Expired.Duration, "JWT token expiration duration.")
}

func NewJWTOptions() *JWTOptions {
	return &JWTOptions{
		Secret:  "",
		Expired: metav1.Duration{Duration: 2 * time.Hour},
	}
}
