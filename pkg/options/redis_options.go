package options

import (
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type RedisOptions struct {
	Addr         string          `json:"addr" yaml:"addr"`
	DialTimeout  metav1.Duration `json:"dial_timeout" yaml:"dial_timeout"`
	ReadTimeout  metav1.Duration `json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout metav1.Duration `json:"write_timeout" yaml:"write_timeout"`
}

func (o *RedisOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.StringVar(&o.Addr, FlagNameFunc(fs, "addr"), o.Addr, "Redis server address.")
	fs.DurationVar(&o.DialTimeout.Duration, FlagNameFunc(fs, "dial-timeout"), o.DialTimeout.Duration, "Redis dial timeout.")
	fs.DurationVar(&o.ReadTimeout.Duration, FlagNameFunc(fs, "read-timeout"), o.ReadTimeout.Duration, "Redis read timeout.")
	fs.DurationVar(&o.WriteTimeout.Duration, FlagNameFunc(fs, "write-timeout"), o.WriteTimeout.Duration, "Redis write timeout.")
}

func NewRedisOptions() *RedisOptions {
	return &RedisOptions{}
}
