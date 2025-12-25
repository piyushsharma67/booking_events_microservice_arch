package logger

import (
	"context"
)

type Logger interface {
	Info(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)

	WithContext(ctx context.Context) Logger
	With(fields ...any) Logger
}

type contextKey string

const RequestIDKey contextKey = "request_id"