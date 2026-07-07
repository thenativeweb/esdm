package logging

import (
	"log/slog"
	"os"
)

var logger *slog.Logger

func init() {
	// By initializing the logger in the init function,
	// it can be used as soon as the package is being
	// imported.
	handler := slog.NewJSONHandler(os.Stdout, nil)
	logger = slog.New(handler)
}

func Debug(message string, args ...any) {
	logger.Debug(message, args...)
}

func Info(message string, args ...any) {
	logger.Info(message, args...)
}

func Warn(message string, args ...any) {
	logger.Warn(message, args...)
}

func Error(message string, args ...any) {
	logger.Error(message, args...)
}

func Fatal(message string, args ...any) {
	logger.Error(message, args...)
	os.Exit(1)
}
