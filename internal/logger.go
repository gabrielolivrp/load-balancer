package internal

import (
	"log"
	"os"
	"path/filepath"
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
	output := getLogOutput(outputFile)
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

func getLogOutput(outputFile string) *os.File {
	if outputFile == "" {
		return os.Stdout
	}

	if _, err := os.Stat(outputFile); os.IsExist(err) {
		panic(err)
	}

	if err := os.MkdirAll(filepath.Dir(outputFile), 0770); err != nil {
		panic(err)
	}

	output, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		panic(err)
	}

	return output
}
