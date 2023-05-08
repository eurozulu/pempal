package logger

import (
	"fmt"
	"io"
	"log"
	"os"
)

const (
	Error LogLevel = iota
	Warning
	Info
	Debug
)

type LogLevel int

type Logger interface {
	Level() LogLevel
	SetLevel(level LogLevel)
	Log(level LogLevel, format string, a ...any)
}

type logger struct {
	level LogLevel
	out   io.Writer
}

var DefaultLogger Logger = &logger{
	level: 0,
	out:   os.Stdout,
}

func (l logger) Level() LogLevel {
	return l.level
}

func (l *logger) SetLevel(level LogLevel) {
	l.level = level
}

func (l logger) Log(level LogLevel, format string, a ...any) {
	if level > l.level {
		return
	}
	if level == Error {
		fmt.Fprintf(os.Stderr, format, a...)
		fmt.Fprintln(os.Stderr)
	} else {
		log.Printf(format, a...)
	}
}

func Log(level LogLevel, format string, a ...any) {
	DefaultLogger.Log(level, format, a...)
}
