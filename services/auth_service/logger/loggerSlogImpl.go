package logger

import (
	"context"
	"log/slog"
	"os"
)

type slogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger(serviceName, env string, level slog.Level) Logger {
	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	})

	base := slog.New(handler).With(
		"service", serviceName,
		"env", env,
	)

	return &slogLogger{logger: base}
}

// Ensure at compile-time that slogLogger implements Logger
var _ Logger = (*slogLogger)(nil)

func (l *slogLogger) Info(msg string, keysAndValues ...any) {
	l.logger.Info(msg, keysAndValues...)
}

func (l *slogLogger) Warn(msg string, keysAndValues ...any) {
	l.logger.Warn(msg, keysAndValues...)
}

func (l *slogLogger) Error(msg string, keysAndValues ...any) {
	l.logger.Error(msg, keysAndValues...)
}

func (l *slogLogger) With(fields ...any) Logger {
	return &slogLogger{
		logger: l.logger.With(fields...),
	}
}

func (l *slogLogger) WithContext(ctx context.Context) Logger {
	if ctx == nil {
		return l
	}

	var fields []any

	if reqID, ok := ctx.Value(RequestIDKey).(string); ok {
		fields = append(fields, "request_id", reqID)
	}

	if len(fields) == 0 {
		return l
	}

	return l.With(fields...)
}
