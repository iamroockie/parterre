package logger

import (
	"context"
	"io"
	"log/slog"
)

type ctxKey struct{}

func New(w io.Writer, level slog.Level) *slog.Logger {
	// nolint:exhaustruct_v5
	return slog.New(slog.NewJSONHandler(w, &slog.HandlerOptions{Level: level}))
}

func NewContext(ctx context.Context, log *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxKey{}, log)
}

func FromContext(ctx context.Context) *slog.Logger {
	if log, ok := ctx.Value(ctxKey{}).(*slog.Logger); ok {
		return log
	}

	return slog.Default()
}
