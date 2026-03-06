package bootstrap

import (
	"github.com/mgcis-cn/ibookfs/pkg/log"
)

// NewLogger creates a pkg/log.Logger from AppInfo.
func NewLogger(info AppInfo) log.Logger {
	return log.NewLogger(info.Name, info.Id, info.Version)
}
