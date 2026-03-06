// Package log provides a unified logging interface based on Kratos log.
package log

import (
	"os"

	kratoslog "github.com/go-kratos/kratos/v2/log"
)

// Logger defines the logging interface for dependency injection.
type Logger interface {
	Log(level kratoslog.Level, keyvals ...any) error
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Warn(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// kratosLogger implements Logger using Kratos log.
type kratosLogger struct {
	logger kratoslog.Logger
	helper *kratoslog.Helper
}

// NewLogger creates a new Logger with service info.
func NewLogger(serviceName, serviceID, version string) Logger {
	logger := kratoslog.With(
		kratoslog.NewStdLogger(os.Stdout),
		"ts", kratoslog.DefaultTimestamp,
		"caller", kratoslog.DefaultCaller,
		"service.name", serviceName,
		"service.id", serviceID,
		"service.version", version,
	)
	return &kratosLogger{
		logger: logger,
		helper: kratoslog.NewHelper(logger),
	}
}

// NewLoggerFromKratos creates a Logger from an existing Kratos logger.
func NewLoggerFromKratos(logger kratoslog.Logger) Logger {
	return &kratosLogger{
		logger: logger,
		helper: kratoslog.NewHelper(logger),
	}
}

// Log implements kratos log.Logger interface.
func (l *kratosLogger) Log(level kratoslog.Level, keyvals ...interface{}) error {
	return l.logger.Log(level, keyvals...)
}

func (l *kratosLogger) Debug(msg string, keyvals ...interface{}) {
	l.helper.Debugw(append([]interface{}{"msg", msg}, keyvals...)...)
}

func (l *kratosLogger) Info(msg string, keyvals ...interface{}) {
	l.helper.Infow(append([]interface{}{"msg", msg}, keyvals...)...)
}

func (l *kratosLogger) Warn(msg string, keyvals ...interface{}) {
	l.helper.Warnw(append([]interface{}{"msg", msg}, keyvals...)...)
}

func (l *kratosLogger) Error(msg string, keyvals ...interface{}) {
	l.helper.Errorw(append([]interface{}{"msg", msg}, keyvals...)...)
}

func (l *kratosLogger) Debugf(format string, args ...interface{}) {
	l.helper.Debugf(format, args...)
}

func (l *kratosLogger) Infof(format string, args ...interface{}) {
	l.helper.Infof(format, args...)
}

func (l *kratosLogger) Warnf(format string, args ...interface{}) {
	l.helper.Warnf(format, args...)
}

func (l *kratosLogger) Errorf(format string, args ...interface{}) {
	l.helper.Errorf(format, args...)
}

// NopLogger is a no-op logger for testing or when logging is disabled.
type NopLogger struct{}

func (NopLogger) Debug(msg string, keyvals ...interface{})                {}
func (NopLogger) Info(msg string, keyvals ...interface{})                 {}
func (NopLogger) Warn(msg string, keyvals ...interface{})                 {}
func (NopLogger) Error(msg string, keyvals ...interface{})                {}
func (NopLogger) Debugf(format string, args ...interface{})               {}
func (NopLogger) Infof(format string, args ...interface{})                {}
func (NopLogger) Warnf(format string, args ...interface{})                {}
func (NopLogger) Errorf(format string, args ...interface{})               {}
func (NopLogger) Log(level kratoslog.Level, keyvals ...interface{}) error { return nil }
