package internal

import (
	"log"
	"os"
)

type LogLevel string

const (
	Info  LogLevel = "info"
	Error LogLevel = "error"
)

type Logger interface {
	Infof(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

type logger struct {
	enabled bool
	level   LogLevel
	logger  *log.Logger
}

func NewLogger(enabled bool, level LogLevel, outputFile string) *logger {
	var output *os.File

	if outputFile != "" {
		var err error
		output, err = os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			panic(err)
		}
	} else {
		output = os.Stdout
	}

	return &logger{
		enabled: enabled,
		level:   level,
		logger:  log.New(output, "", log.LstdFlags),
	}
}

func (l *logger) Infof(format string, args ...interface{}) {
	if l.enabled && l.level == Info {
		l.logger.Printf("[INFO] "+format, args...)
	}
}

func (l *logger) Errorf(format string, args ...interface{}) {
	if l.enabled {
		l.logger.Printf("[ERROR] "+format, args...)
	}
}
