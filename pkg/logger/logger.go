package logger

import (
	"electra/internal/config"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

type appLogger struct {
	*slog.Logger
}

// New создает новый экземпляр логгера
func NewLogger(cfg config.Config) (Logger, error) {
	var level slog.Level
	switch cfg.Logger.Level {
	case "debug":
		level = slog.LevelDebug
	case "info":
		level = slog.LevelInfo
	case "error":
		level = slog.LevelError
	case "warn":
		level = slog.LevelWarn
	default:
		level = slog.LevelInfo
	}

	if cfg.Logger.FilePath != "" {
		logDir := filepath.Dir(cfg.Logger.FilePath)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return nil, err
		}
	}

	lumberjackLogger := &lumberjack.Logger{
		Filename: cfg.Logger.FilePath,
		MaxSize:  cfg.Logger.MaxSize,
	}

	opts := &slog.HandlerOptions{
		Level:     level,
		AddSource: true,
	}

	var handler slog.Handler
	handler = slog.NewTextHandler(lumberjackLogger, opts)

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return &appLogger{Logger: logger}, nil
}
