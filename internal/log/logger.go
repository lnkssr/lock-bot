package logger

import (
	"fmt"
	"log"
	"os"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARN
	ERROR
)

var (
	currentLevel = DEBUG
	stdLogger    = log.New(os.Stdout, "", log.LstdFlags)
)

func SetLevel(l Level) {
	currentLevel = l
}

func logf(l Level, prefix string, args ...any) {
	if l >= currentLevel {
		stdLogger.Printf("%s %s", prefix, fmt.Sprint(args...))
	}
}

func Debug(args ...any) { logf(DEBUG, "[DEBUG]", args...) }
func Info(args ...any)  { logf(INFO, "[INFO]", args...) }
func Warn(args ...any)  { logf(WARN, "[WARN]", args...) }
func Error(args ...any) { logf(ERROR, "[ERROR]", args...) }
