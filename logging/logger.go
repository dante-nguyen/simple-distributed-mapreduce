package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

const timestampFormat = time.RFC3339

type Logger struct {
	logger *log.Logger
	Name   string
	level  Level
}

func NewLogger(name string, level Level) Logger {
	assertValidLevel(level)
	return Logger{
		logger: log.New(os.Stderr, "", 0),
		Name:   name,
		level:  level,
	}
}

func (l *Logger) SetLevel(level Level) {
	assertValidLevel(level)
	l.level = level
}

func (l *Logger) Debug(format string, args ...any) {
	l.logAt(DEBUG, format, args...)
}

func (l *Logger) Info(format string, args ...any) {
	l.logAt(INFO, format, args...)
}

func (l *Logger) Warning(format string, args ...any) {
	l.logAt(WARNING, format, args...)
}

func (l *Logger) Error(format string, args ...any) {
	l.logAt(ERROR, format, args...)
}

func (l *Logger) Fatal(format string, args ...any) {
	l.logAt(FATAL, format, args...)
	os.Exit(1)
}

func (l *Logger) logAt(level Level, format string, args ...any) {
	if l.level.Allow(level) {
		msg := fmt.Sprintf(format, args...)
		timestamp := time.Now().Format(timestampFormat)
		levelName := level.Name()
		l.logger.Printf("%s - [%s] [%s]: %s\n", timestamp, l.Name, levelName, msg)
	}
}
