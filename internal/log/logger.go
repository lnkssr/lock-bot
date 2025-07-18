package logger

import (
	"os"
	"strings"

	"lockbot/internal/config"

	"github.com/charmbracelet/log"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

func parseLevel(env string) log.Level {
	switch strings.ToUpper(env) {
	case "DEBUG":
		return log.DebugLevel
	case "INFO":
		return log.InfoLevel
	case "WARN":
		return log.WarnLevel
	case "ERROR":
		return log.ErrorLevel
	default:
		return log.InfoLevel
	}
}

var (
	stdLogger = log.NewWithOptions(os.Stderr, log.Options{
		ReportTimestamp: true,
		Prefix:          "api client log:",
	})
)

func init() {
	stdLogger.SetLevel(parseLevel(config.LogLevel))
}

func SetLevel(l Level) {
	switch l {
	case DEBUG:
		stdLogger.SetLevel(log.DebugLevel)
	case INFO:
		stdLogger.SetLevel(log.InfoLevel)
	case WARN:
		stdLogger.SetLevel(log.WarnLevel)
	case ERROR:
		stdLogger.SetLevel(log.ErrorLevel)
	}
}

func Debug(args ...any) { stdLogger.Debug(args[0], args[1:]...) }
func Info(args ...any)  { stdLogger.Info(args[0], args[1:]...) }
func Warn(args ...any)  { stdLogger.Warn(args[0], args[1:]...) }
func Error(args ...any) { stdLogger.Error(args[0], args[1:]...) }
