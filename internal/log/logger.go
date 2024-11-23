package log

import (
	"log/slog"
	"os"
)

type Level int

type LoggerInterface interface {
	Info(msg string, args ...any)
	Debug(msg string, args ...any)
	Warn(msg string, args ...any)
	Error(msg string, args ...any)
	Handler() slog.Handler
}

func NewLogger() LoggerInterface {
	return slog.New(slog.NewJSONHandler(os.Stdout, nil))
}
