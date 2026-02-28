package options

import (
	"github.com/spf13/pflag"
)

type NeteaseOptions struct {
	Host     string `json:"host" yaml:"host"`
	Port     int    `json:"port" yaml:"port"`
	User     string `json:"user" yaml:"user"`
	Password string `json:"password" yaml:"password"`
	From     string `json:"from" yaml:"from"`
	FromName string `json:"from_name" yaml:"from_name"`
}

func (o *NeteaseOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.StringVar(&o.Host, FlagNameFunc(fs, "host"), o.Host, "Netease SMTP server host.")
	fs.IntVar(&o.Port, FlagNameFunc(fs, "port"), o.Port, "Netease SMTP server port.")
	fs.StringVar(&o.User, FlagNameFunc(fs, "user"), o.User, "Netease SMTP auth username.")
	fs.StringVar(&o.Password, FlagNameFunc(fs, "password"), o.Password, "Netease SMTP auth password.")
	fs.StringVar(&o.From, FlagNameFunc(fs, "from"), o.From, "Sender email address.")
	fs.StringVar(&o.FromName, FlagNameFunc(fs, "from-name"), o.FromName, "Sender display name.")
}

func (o *NeteaseOptions) Kind() string {
	return "netease"
}

func NewNeteaseOptions() *NeteaseOptions {
	return &NeteaseOptions{}
}
