package options

import (
	"github.com/spf13/pflag"
)

type GithubOptions struct {
	ClientID     string `json:"client_id" yaml:"client_id"`
	ClientSecret string `json:"client_secret" yaml:"client_secret"`
	RedirectURL  string `json:"redirect_url" yaml:"redirect_url"`
}

func (o *GithubOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.StringVar(&o.ClientID, FlagNameFunc(fs, "client-id"), o.ClientID, "GitHub OAuth client ID.")
	fs.StringVar(&o.ClientSecret, FlagNameFunc(fs, "client-secret"), o.ClientSecret, "GitHub OAuth client secret.")
	fs.StringVar(&o.RedirectURL, FlagNameFunc(fs, "redirect-url"), o.RedirectURL, "GitHub OAuth redirect URL.")
}

func NewGithubOptions() *GithubOptions {
	return &GithubOptions{}
}

func (o *GithubOptions) Kind() string {
	return "github"
}
