package logger

import (
	"log/slog"
	"os"
	"strings"

	"github.com/lmittmann/tint"
)

// Init initializes a new slog logger based on the provided format.
// Supported formats: "json", "text". Defaults to "text".
func Init(format string) *slog.Logger {
	var handler slog.Handler

	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	switch strings.ToLower(format) {
	case "json":
		handler = slog.NewJSONHandler(os.Stdout, opts)
	default:
		handler = tint.NewHandler(os.Stdout, &tint.Options{
			Level:      slog.LevelInfo,
			TimeFormat: "15:04:05",
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}
