package logging

import (
	"fmt"
	"io"
	"os"
	"time"
)

const timeformat = "2006-01-02 15:04:05.111"

const (
	LogError LogLevel = iota
	LogWarning
	LogInfo
	LogDebug
	LogTrace
)

// DefaultLogger is the logger used for the package.
var DefaultLogger Logger = &logger{
	level: LogError,
	out:   os.Stdout,
}

var LogLevelNames []string = []string{
	"Error",
	"Warning",
	"Info",
	"Debug",
	"Trace",
}

// LogLevel sets the level of logging written to the log
type LogLevel int

type Logger interface {
	Log(level LogLevel, msg string, a ...any)
	LogLevel() LogLevel
	SetLogLevel(level LogLevel)
	SetOutput(out io.Writer)
}

type logger struct {
	level LogLevel
	out   io.Writer
}

func (l logger) Log(level LogLevel, msg string, a ...any) {
	if level > l.level {
		return
	}

	fmt.Fprintf(l.out, "[%s] [%s] %s\n", time.Now().Format(timeformat), level, fmt.Sprintf(msg, a...))
}

func (l logger) LogLevel() LogLevel {
	return l.level
}

func (l *logger) SetLogLevel(level LogLevel) {
	l.level = level
}

func (l *logger) SetOutput(out io.Writer) {
	l.out = out
}

func (ll LogLevel) String() string {
	if ll < 0 || int(ll) >= len(LogLevelNames) {
		return ""
	}
	return LogLevelNames[ll]
}

func Error(msg string, a ...any) {
	DefaultLogger.Log(LogError, msg, a...)
}

func Warning(msg string, a ...any) {
	DefaultLogger.Log(LogWarning, msg, a...)
}

func Info(msg string, a ...any) {
	DefaultLogger.Log(LogInfo, msg, a...)
}

func Debug(msg string, a ...any) {
	DefaultLogger.Log(LogDebug, msg, a...)
}

func Trace(msg string, a ...any) {
	DefaultLogger.Log(LogTrace, msg, a...)
}

func IsLogLevel(level LogLevel) bool {
	return DefaultLogger.LogLevel() >= level
}
