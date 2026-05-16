package logger

import (
	"context"
	"strings"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/slog"
)

type LogrusHandler struct {
	logger *logrus.Logger
}

func NewLogrusHandler(logger *logrus.Logger) *LogrusHandler {
	return &LogrusHandler{
		logger: logger,
	}
}

func ConvertLogLevel(level string) logrus.Level {
	var l logrus.Level

	switch strings.ToLower(level) {
	case "error":
		l = logrus.ErrorLevel
	case "warm", "warn": // добавил warn на случай опечатки в конфиге
		l = logrus.WarnLevel
	case "info":
		l = logrus.InfoLevel
	case "debug":
		l = logrus.DebugLevel
	default:
		l = logrus.InfoLevel
	}

	return l
}

func (h *LogrusHandler) Enabled(ctx context.Context, level slog.Level) bool {
	// Поддерживаем все уровни логирования
	return true
}

func (h *LogrusHandler) Handle(ctx context.Context, rec slog.Record) error {
	fields := make(map[string]interface{}, rec.NumAttrs())

	rec.Attrs(func(a slog.Attr) {
		fields[a.Key] = a.Value.Any()
	})

	entry := h.logger.WithFields(fields)

	switch rec.Level {
	case slog.LevelDebug:
		entry.Debug(rec.Message)
	case slog.LevelInfo:
		entry.Info(rec.Message)
	case slog.LevelWarn:
		entry.Warn(rec.Message)
	case slog.LevelError:
		entry.Error(rec.Message)
	}

	return nil
}

func (h *LogrusHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	// Возвращаем h для соответствия интерфейсу
	return h
}

func (h *LogrusHandler) WithGroup(name string) slog.Handler {
	// Возвращаем h для соответствия интерфейсу
	return h
}
