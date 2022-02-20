package logger

import (
	"fmt"
	"log"
)

type Logger interface {
	Errorf(format string, args ...interface{})
	Warningf(format string, args ...interface{})
	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
}
type LogLevel int8

const (
	Error   LogLevel = 1
	Warning LogLevel = 2
	Info    LogLevel = 3
	Debug   LogLevel = 4
)

type logger struct {
	Level  LogLevel
	Source string
}

func New(level LogLevel, source string) logger {
	return logger{
		Level:  level,
		Source: source,
	}
}

func (l logger) logf(logType string, format string, data ...interface{}) {
	line := fmt.Sprintf(format, data...)
	log.Printf("%s: %s%s\n", l.Source, logType, line)
}

func (l logger) Errorf(format string, args ...interface{}) {
	// showing errors can't be disabled
	l.logf("ERROR  ", format, args...)
}

func (l logger) Warningf(format string, args ...interface{}) {
	if l.Level >= Info {
		l.logf("WARNING  ", format, args...)
	}
}

func (l logger) Debugf(format string, args ...interface{}) {
	if l.Level >= Debug {
		l.logf("DEBUG  ", format, args...)
	}
}

func (l logger) Infof(format string, args ...interface{}) {
	if l.Level >= Info {
		l.logf("INFO  ", format, args...)
	}
}
