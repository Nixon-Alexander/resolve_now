package config

import (
	"log/slog"
	"os"
)

type Log struct {
	log *slog.Logger
}

func InitLog() *Log {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{AddSource: true}))

	return &Log{
		log: logger,
	}
}

func (l *Log) Info(message string, args ...any) {
	l.log.Info(message, args...)
}

func (l *Log) Warn(message string, args ...any) {
	l.log.Warn(message, args...)
}

func (l *Log) Error(message string, args ...any) {
	l.log.Error(message, args...)
}
