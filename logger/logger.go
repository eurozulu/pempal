package logger

import (
	"fmt"
	"io"
	"os"
	"time"
)

const (
	LevelError LogLevel = iota
	LevelInfo
	LevelWarning
	LevelDebug
)

const timeFormat = time.RFC3339

type LogLevel int

type Logger interface {
	Level() LogLevel
	SetLevel(level LogLevel)
	Log(level LogLevel, format string, a ...any)
	ShowTimeStamp() bool
	SetShowTimeStamp(b bool)
}

type logger struct {
	level         LogLevel
	out           io.Writer
	showTimestamp bool
}

var DefaultLogger Logger = &logger{
	level:         0,
	out:           os.Stdout,
	showTimestamp: true,
}

func (l logger) Level() LogLevel {
	return l.level
}

func (l *logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l logger) Output() io.Writer {
	return l.out
}
func (l *logger) SetOutput(out io.Writer) {
	l.out = out
}

func (l logger) ShowTimeStamp() bool {
	return l.showTimestamp
}

func (l *logger) SetShowTimeStamp(b bool) {
	l.showTimestamp = b
}

func (l logger) Log(level LogLevel, format string, a ...any) {
	if level > l.level {
		return
	}
	if l.showTimestamp {
		fmt.Fprintf(l.out, "%s\t", time.Now().Format(timeFormat))
	}
	fmt.Fprintf(l.out, format, a...)
	fmt.Fprintln(l.out)

}

func NewLogger(out io.Writer, showTimestamp bool, level LogLevel) Logger {
	return &logger{
		level:         level,
		out:           out,
		showTimestamp: showTimestamp,
	}
}

func Debug(format string, a ...any) {
	DefaultLogger.Log(LevelDebug, format, a...)
}

func Warning(format string, a ...any) {
	DefaultLogger.Log(LevelWarning, format, a...)
}

func Info(format string, a ...any) {
	DefaultLogger.Log(LevelInfo, format, a...)
}

func Error(format string, a ...any) {
	DefaultLogger.Log(LevelError, format, a...)
}
