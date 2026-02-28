package options

import (
	"github.com/spf13/pflag"
)

type GiteeOptions struct {
	ClientID     string `json:"client_id" yaml:"client_id"`
	ClientSecret string `json:"client_secret" yaml:"client_secret"`
	RedirectURL  string `json:"redirect_url" yaml:"redirect_url"`
}

func (o *GiteeOptions) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.StringVar(&o.ClientID, FlagNameFunc(fs, "client-id"), o.ClientID, "Gitee OAuth client ID.")
	fs.StringVar(&o.ClientSecret, FlagNameFunc(fs, "client-secret"), o.ClientSecret, "Gitee OAuth client secret.")
	fs.StringVar(&o.RedirectURL, FlagNameFunc(fs, "redirect-url"), o.RedirectURL, "Gitee OAuth redirect URL.")
}

func NewGiteeOptions() *GiteeOptions {
	return &GiteeOptions{}
}

func (o *GiteeOptions) Kind() string {
	return "gitee"
}
