package options

import (
	"github.com/spf13/pflag"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MysqlOptions struct {
	Source          string          `json:"source" yaml:"source" mapstructure:"source"`
	MaxIdleConns    int             `json:"max_idle_conns" yaml:"max_idle_conns" mapstructure:"max_idle_conns"`
	MaxOpenConns    int             `json:"max_open_conns" yaml:"max_open_conns" mapstructure:"max_open_conns"`
	ConnMaxLifetime metav1.Duration `json:"conn_max_lifetime" yaml:"conn_max_lifetime" mapstructure:"conn_max_lifetime"`
	ConnMaxIdleTime metav1.Duration `json:"conn_max_idle_time" yaml:"conn_max_idle_time" mapstructure:"conn_max_idle_time"`
}

func NewMysqlOptions() *MysqlOptions {
	return &MysqlOptions{}
}

func (o *MysqlOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.StringVar(&o.Source, FlagNameFunc(fs, "source"), o.Source, "MySQL DSN connection string.")
	fs.IntVar(&o.MaxIdleConns, FlagNameFunc(fs, "max-idle-conns"), o.MaxIdleConns, "Maximum number of idle connections in the pool.")
	fs.IntVar(&o.MaxOpenConns, FlagNameFunc(fs, "max-open-conns"), o.MaxOpenConns, "Maximum number of open connections to the database.")
	fs.DurationVar(&o.ConnMaxLifetime.Duration, FlagNameFunc(fs, "conn-max-lifetime"), o.ConnMaxLifetime.Duration, "Maximum amount of time a connection may be reused.")
	fs.DurationVar(&o.ConnMaxIdleTime.Duration, FlagNameFunc(fs, "conn-max-idle-time"), o.ConnMaxIdleTime.Duration, "Maximum amount of time a connection may be idle.")
}

func (o *MysqlOptions) Kind() string {
	return "mysql"
}
