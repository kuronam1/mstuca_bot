package logger

import (
	"log/slog"
	"os"
)

type Logger interface {
	Fatal(msg string, args ...any)
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Error(msg string, args ...any)
	Warn(msg string, args ...any)
}

type logger struct {
	log *slog.Logger
}

func New(slogger *slog.Logger) Logger {
	return &logger{
		log: slogger,
	}
}

func (l *logger) Fatal(msg string, args ...any) {
	l.log.Error(msg, args...)
	os.Exit(1)
}

func (l *logger) Info(msg string, args ...any) {
	l.log.Info(msg, args...)
}

func (l *logger) Debug(msg string, args ...any) {
	l.log.Debug(msg, args...)
}

func (l *logger) Error(msg string, args ...any) {
	l.log.Error(msg, args...)
}

func (l *logger) Warn(msg string, args ...any) {
	l.log.Warn(msg, args...)
}
