package ctxmeta

import (
	"context"
)

// Неэкспортируемый тип ключа — защита от коллизий.
type key int

const (
	keyRequestID key = iota
	keyUserID
	keyTraceID
)

// ------- RequestID -------
func WithRequestID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyRequestID, id)
}

func RequestID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(keyRequestID).(string)
	return v, ok && v != ""
}

// ------- UserID (int64) -------
func WithUserID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, keyUserID, id)
}
func UserID(ctx context.Context) (int64, bool) {
	v, ok := ctx.Value(keyUserID).(int64)
	return v, ok && v != 0
}

// ------- TraceID (если используется трейсинг вручную) -------
func WithTraceID(ctx context.Context, id string) context.Context {
	return context.WithValue(ctx, keyTraceID, id)
}
func TraceID(ctx context.Context) (string, bool) {
	v, ok := ctx.Value(keyTraceID).(string)
	return v, ok && v != ""
}
