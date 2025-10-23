package http

import "context"

// Интерфейс Logger определяет контракт для логгера сервера HTTP.
type Logger interface {
	Info(msg string, v ...any)
	Error(msg string, v ...any)
	Debug(msg string, v ...any)
	Warn(msg string, v ...any)

	InfoCtx(ctx context.Context, msg string, v ...any)
	ErrorCtx(ctx context.Context, msg string, v ...any)
	DebugCtx(ctx context.Context, msg string, v ...any)
	WarnCtx(ctx context.Context, msg string, v ...any)
}
