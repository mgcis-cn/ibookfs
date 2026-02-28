package options

import (
	oauthSources "github.com/mgcis-cn/ibookfs/pkg/authn/oauth/sources"
	"github.com/mgcis-cn/ibookfs/pkg/database/sources"
	emailSources "github.com/mgcis-cn/ibookfs/pkg/email/sources"
	"github.com/mgcis-cn/ibookfs/pkg/options"
	storageSources "github.com/mgcis-cn/ibookfs/pkg/storage/sources"

	"github.com/spf13/pflag"
	"k8s.io/component-base/cli/flag"
)

type ServerRunOptions struct {
	App    *AppOptions     `json:"app" yaml:"app" mapstructure:"app"`
	Server *ServerOptions  `json:"server" yaml:"server" mapstructure:"server"`
	Data   *DataOptions    `json:"data" yaml:"data" mapstructure:"data"`
	Auth   *AuthOptions    `json:"auth" yaml:"auth" mapstructure:"auth"`
	Email  []*EmailOptions `json:"email" yaml:"email" mapstructure:"email"`
}

func defaultDatabaseOptions() *DatabaseOptions {
	return &DatabaseOptions{
		Name:   "default",
		SQLite: options.NewSqliteOptions(),
	}
}

func defaultStorageOptions() *StorageOptions {
	return &StorageOptions{
		Name:  "default",
		Local: options.NewLocalOSOptions(),
	}
}

func defaultEmailOptions() *EmailOptions {
	return &EmailOptions{
		Name:    "default",
		Netease: options.NewNeteaseOptions(),
	}
}

func defaultOAuthOptions() *OAuthOptions {
	return &OAuthOptions{
		Name:  "default",
		Gitee: options.NewGiteeOptions(),
	}
}

func NewServerRunOptions() *ServerRunOptions {
	return &ServerRunOptions{
		App: &AppOptions{},
		Server: &ServerOptions{
			HTTP:   options.NewHttpOptions(),
			Image:  options.NewImageOptions(),
			Upload: options.NewUploadOptions(),
		},
		Data: &DataOptions{
			Database: []*DatabaseOptions{
				defaultDatabaseOptions(),
			},
			Redis: options.NewRedisOptions(),
			Storage: []*StorageOptions{
				defaultStorageOptions(),
			},
		},
		Auth: &AuthOptions{
			JWT: options.NewJWTOptions(),
			OAuth: []*OAuthOptions{
				defaultOAuthOptions(),
			},
		},
		Email: []*EmailOptions{
			defaultEmailOptions(),
		},
	}
}

func (o *ServerRunOptions) AddFlags() (fss flag.NamedFlagSets) {
	o.App.AddFlags(fss.FlagSet("app"))
	o.Server.AddFlags(fss.FlagSet("server"))
	o.Data.AddFlags(fss.FlagSet("data"))
	o.Auth.AddFlags(fss.FlagSet("auth"))
	o.Email[0].AddFlags(fss.FlagSet("email"))
	return fss
}

type AppOptions struct {
	Name    string `json:"name" yaml:"name" mapstructure:"name"`
	Version string `json:"version" yaml:"version" mapstructure:"version"`
}

func (o *AppOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.StringVar(&o.Name, "name", o.Name, "Application name.")
	fs.StringVar(&o.Version, "version", o.Version, "Application version.")
}

type ServerOptions struct {
	HTTP   *options.HttpOptions   `json:"http" yaml:"http" mapstructure:"http"`
	Image  *options.ImageOptions  `json:"image" yaml:"image" mapstructure:"image"`
	Upload *options.UploadOptions `json:"upload" yaml:"upload" mapstructure:"upload"`
}

func (o *ServerOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	var fss flag.NamedFlagSets
	o.HTTP.AddFlags(fss.FlagSet(options.FlagNameFunc(fs, "http")))
	o.Image.AddFlags(fss.FlagSet(options.FlagNameFunc(fs, "image")))
	o.Upload.AddFlags(fss.FlagSet(options.FlagNameFunc(fs, "upload")))
	for _, f := range fss.FlagSets {
		fs.AddFlagSet(f)
	}
}

type DataOptions struct {
	Database []*DatabaseOptions    `json:"database" yaml:"database" mapstructure:"database"`
	Redis    *options.RedisOptions `json:"redis" yaml:"redis" mapstructure:"redis"`
	Storage  []*StorageOptions     `json:"storage" yaml:"storage" mapstructure:"storage"`
}

func (o *DataOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	var fss flag.NamedFlagSets
	if len(o.Database) > 0 {
		o.Database[0].AddFlags(fss.FlagSet(options.FlagNameFunc(fs, "database")))
	}
	o.Redis.AddFlags(fss.FlagSet(options.FlagNameFunc(fs, "redis")))
	if len(o.Storage) > 0 {
		o.Storage[0].AddFlags(fss.FlagSet(options.FlagNameFunc(fs, "storage")))
	}
	for _, f := range fss.FlagSets {
		fs.AddFlagSet(f)
	}
}

type DatabaseOptions struct {
	Name     string                   `json:"name" yaml:"name" mapstructure:"name"`
	SQLite   *options.SqliteOptions   `json:"sqlite" yaml:"sqlite" mapstructure:"sqlite"`
	MySQL    *options.MysqlOptions    `json:"mysql" yaml:"mysql" mapstructure:"mysql"`
	Postgres *options.PostgresOptions `json:"postgres" yaml:"postgres" mapstructure:"postgres"`
}

func (o *DatabaseOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	var fss flag.NamedFlagSets
	o.SQLite.AddFlags(fss.FlagSet(options.FlagNameFunc(fs, "sqlite")))
	o.MySQL.AddFlags(fss.FlagSet(options.FlagNameFunc(fs, "mysql")))
	o.Postgres.AddFlags(fss.FlagSet(options.FlagNameFunc(fs, "postgres")))
	for _, f := range fss.FlagSets {
		fs.AddFlagSet(f)
	}
}

type StorageOptions struct {
	Name      string                    `json:"name" yaml:"name" mapstructure:"name"`
	Local     *options.LocalOSOptions   `json:"local" yaml:"local" mapstructure:"local"`
	AliyunOSS *options.AliyunOSSOptions `json:"aliyunoss" yaml:"aliyunoss" mapstructure:"aliyunoss"`
}

func (o *StorageOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	var fss flag.NamedFlagSets
	o.Local.AddFlags(fss.FlagSet(options.FlagNameFunc(fs, "local")))
	o.AliyunOSS.AddFlags(fss.FlagSet(options.FlagNameFunc(fs, "aliyunoss")))
	for _, f := range fss.FlagSets {
		fs.AddFlagSet(f)
	}
}

type AuthOptions struct {
	JWT   *options.JWTOptions `json:"jwt" yaml:"jwt" mapstructure:"jwt"`
	OAuth []*OAuthOptions     `json:"oauth" yaml:"oauth" mapstructure:"oauth"`
}

func (o *AuthOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	var fss flag.NamedFlagSets
	o.JWT.AddFlags(fss.FlagSet(options.FlagNameFunc(fs, "jwt")))
	if len(o.OAuth) > 0 {
		o.OAuth[0].AddFlags(fss.FlagSet(options.FlagNameFunc(fs, "oauth")))
	}
	for _, f := range fss.FlagSets {
		fs.AddFlagSet(f)
	}
}

type OAuthOptions struct {
	Name   string                 `json:"name" yaml:"name" mapstructure:"name"`
	GitHub *options.GithubOptions `json:"github" yaml:"github" mapstructure:"github"`
	Gitee  *options.GiteeOptions  `json:"gitee" yaml:"gitee" mapstructure:"gitee"`
}

func (o *OAuthOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.StringVar(&o.Name, options.FlagNameFunc(fs, "name"), o.Name, "OAuth provider name.")
	var fss flag.NamedFlagSets
	o.GitHub.AddFlags(fss.FlagSet(options.FlagNameFunc(fs, "github")))
	o.Gitee.AddFlags(fss.FlagSet(options.FlagNameFunc(fs, "gitee")))
	for _, f := range fss.FlagSets {
		fs.AddFlagSet(f)
	}
}

// ActiveConfig returns the active OAuth configuration implementing oauth sources.Config.
func (o *OAuthOptions) ActiveConfig() oauthSources.Config {
	if o.GitHub != nil && o.GitHub.ClientID != "" {
		return o.GitHub
	}
	if o.Gitee != nil && o.Gitee.ClientID != "" {
		return o.Gitee
	}
	return nil
}

type EmailOptions struct {
	Name    string                  `json:"name" yaml:"name" mapstructure:"name"`
	Netease *options.NeteaseOptions `json:"netease" yaml:"netease" mapstructure:"netease"`
}

func (o *EmailOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.StringVar(&o.Name, options.FlagNameFunc(fs, "name"), o.Name, "Email provider name (e.g. netease).")
	var fss flag.NamedFlagSets
	o.Netease.AddFlags(fss.FlagSet(options.FlagNameFunc(fs, "netease")))
	for _, f := range fss.FlagSets {
		fs.AddFlagSet(f)
	}
}

// ActiveConfig returns the active email configuration implementing email sources.Config.
func (o *EmailOptions) ActiveConfig() emailSources.Config {
	if o.Netease != nil && o.Netease.Host != "" {
		return o.Netease
	}
	return nil
}

// ActiveConfig returns the active database configuration implementing sources.Config.
func (o *DatabaseOptions) ActiveConfig() sources.Config {
	if o.MySQL != nil && o.MySQL.Source != "" {
		return o.MySQL
	}
	if o.Postgres != nil && o.Postgres.Source != "" {
		return o.Postgres
	}
	if o.SQLite != nil && o.SQLite.Source != "" {
		return o.SQLite
	}
	return nil
}

// ActiveConfig returns the active storage configuration implementing storage sources.Config.
func (o *StorageOptions) ActiveConfig() storageSources.Config {
	if o.Local != nil && o.Local.Bucket != "" {
		return o.Local
	}
	if o.AliyunOSS != nil && o.AliyunOSS.Bucket != "" {
		return o.AliyunOSS
	}
	return nil
}
