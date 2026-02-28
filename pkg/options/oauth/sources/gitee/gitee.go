package gitee

import (
	"github.com/mgcis-cn/ibookfs/pkg/options"

	"github.com/spf13/pflag"
)

const SourceKind = "gitee"

type Source struct {
	ClientID     string `json:"client_id" yaml:"client_id"`
	ClientSecret string `json:"client_secret" yaml:"client_secret"`
	RedirectURL  string `json:"redirect_url" yaml:"redirect_url"`
}

func (o *Source) AddFlags(fs *pflag.FlagSet) {
	if o == nil {
		return
	}
	fs.StringVar(&o.ClientID, options.FlagNameFunc(fs, "client-id"), o.ClientID, "Gitee OAuth client ID.")
	fs.StringVar(&o.ClientSecret, options.FlagNameFunc(fs, "client-secret"), o.ClientSecret, "Gitee OAuth client secret.")
	fs.StringVar(&o.RedirectURL, options.FlagNameFunc(fs, "redirect-url"), o.RedirectURL, "Gitee OAuth redirect URL.")
}

func New() *Source {
	return &Source{}
}
