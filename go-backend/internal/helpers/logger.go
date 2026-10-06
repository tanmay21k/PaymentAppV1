package helpers

import (
	"log/slog"
	"os"
)

func New() *slog.Logger {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	logger := slog.NewJSONHandler(os.Stdout, opts)

	return slog.New(logger)
}
