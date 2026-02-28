package bootstrap

import "github.com/go-kratos/kratos/v2/log"

func NewLogger(info AppInfo) log.Logger {
	return log.With(
		log.DefaultLogger,
		"ts", log.DefaultTimestamp,
		"caller", log.DefaultCaller,
		"service.id", info.Id,
		"service.name", info.Name,
		"service.version", info.Version,
	)
}
