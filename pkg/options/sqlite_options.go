package options

import "github.com/spf13/pflag"

type SqliteOptions struct {
	Source string `json:"source" yaml:"source" mapstructure:"source"`
}

func NewSqliteOptions() *SqliteOptions {
	return &SqliteOptions{}
}

func (o *SqliteOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.StringVar(&o.Source, FlagNameFunc(fs, "source"), o.Source, "SQLite database file path.")
}

func (o *SqliteOptions) Kind() string {
	return "sqlite"
}
