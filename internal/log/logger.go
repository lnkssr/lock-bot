package logger

import (
	"fmt"
	"lockbot/internal/config"
	"log"
	"os"
	"strings"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

func parcerLevel(env string) Level {
	switch strings.ToUpper(env) {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	default:
		return INFO 
	}
}

var (
	currentLevel = parcerLevel(config.LogLevel) 
	stdLogger    = log.New(os.Stdout, "", log.LstdFlags)
)

func SetLevel(l Level) {
	currentLevel = l
}

func logf(l Level, prefix string, args ...any) {
	if l >= currentLevel {
		stdLogger.Printf("%s %s", prefix, fmt.Sprintln(args...))
	}
}

func Debug(args ...any) { logf(DEBUG, "[DEBUG]", args...) }
func Info(args ...any)  { logf(INFO, "[INFO]", args...) }
func Warn(args ...any)  { logf(WARN, "[WARN]", args...) }
func Error(args ...any) { logf(ERROR, "[ERROR]", args...) }
