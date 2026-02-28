package app

import (
	"strings"

	"github.com/mgcis-cn/ibookfs/cmd/apiserver/app/options"
	"github.com/mgcis-cn/ibookfs/internal/apiserver"
	_ "github.com/mgcis-cn/ibookfs/pkg/authn/oauth/sources/gitee"
	_ "github.com/mgcis-cn/ibookfs/pkg/authn/oauth/sources/github"
	_ "github.com/mgcis-cn/ibookfs/pkg/database/sources/mysql"
	_ "github.com/mgcis-cn/ibookfs/pkg/database/sources/postgres"
	_ "github.com/mgcis-cn/ibookfs/pkg/database/sources/sqlite"
	_ "github.com/mgcis-cn/ibookfs/pkg/email/sources/netease"
	_ "github.com/mgcis-cn/ibookfs/pkg/storage/sources/local"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"k8s.io/component-base/cli/flag"
	"k8s.io/component-base/term"
)

const appName = "apiserver"

func NewAPIServerCommand() *cobra.Command {
	runOpts := options.NewServerRunOptions()
	cmd := &cobra.Command{
		Use:   appName,
		Short: "iBookFS API Server",
		Long: `iBookFS API Server provides HTTP APIs for book and image management,
user authentication (JWT, OAuth), email services, and file storage.

Configuration is loaded from YAML files in the config directory (--conf).
Use --env to specify the runtime environment (dev, test, prod) which overlays
application-{env}.yaml on top of the base application.yaml.`,
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := Load()
			if err != nil {
				return err
			}
			if err := c.Scan(runOpts); err != nil {
				return err
			}
			serverCfg := apiserver.NewServerConfig(runOpts)
			app, cleanup, err := apiserver.New(serverCfg)
			if err != nil {
				return err
			}
			defer cleanup()
			return app.Run()
		},
	}

	fs := cmd.Flags()
	AddConfigFlags(fs)
	fss := runOpts.AddFlags()
	for _, f := range fss.FlagSets {
		fs.AddFlagSet(f)
	}

	cols, _, _ := term.TerminalSize(cmd.OutOrStdout())
	fs.SetNormalizeFunc(normalizeFunc)
	flag.SetUsageAndHelpFunc(cmd, fss, cols)
	return cmd
}

func normalizeFunc(f *pflag.FlagSet, name string) pflag.NormalizedName {
	from := []string{"-", "."}
	to := "_"
	for _, sep := range from {
		name = strings.Replace(name, sep, to, -1)
	}
	return pflag.NormalizedName(name)
}
