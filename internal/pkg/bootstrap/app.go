package bootstrap

import (
	"os"

	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport"
)

type AppInfo struct {
	Id      string
	Name    string
	Version string
}

func NewAppInfo(id, name, version string) AppInfo {
	if id == "" {
		id, _ = os.Hostname()
	}
	return AppInfo{
		Id:      id,
		Name:    name,
		Version: version,
	}
}

type AppConfig struct {
	Info   AppInfo
	Logger log.Logger
}

func NewApp(c AppConfig, servers ...transport.Server) *kratos.App {
	return kratos.New(
		kratos.ID(c.Info.Id),
		kratos.Name(c.Info.Name),
		kratos.Version(c.Info.Version),
		kratos.Logger(c.Logger),
		kratos.Server(servers...),
	)
}
