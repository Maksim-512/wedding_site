package logger

import (
	"log/slog"
	"os"
	"wedding_website/internal/lib/logger/sl"
)

func SetupLogger() *slog.Logger {
	myLogger := slog.New(
		slog.NewTextHandler(
			os.Stdout,
			&slog.HandlerOptions{Level: slog.LevelDebug},
		),
	)

	return myLogger
}

func ConfLogger(baseLogger *slog.Logger, op string) *slog.Logger {
	return baseLogger.With(
		sl.Op(op),
	)
}
