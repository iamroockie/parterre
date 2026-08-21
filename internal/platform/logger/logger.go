package logger

import (
	"context"
	"io"
	"log/slog"
)

type ctxKey struct{}

type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
)

func New(w io.Writer, format Format, l slog.Level) *slog.Logger {
	var handler slog.Handler
	opts := slog.HandlerOptions{Level: l}

	switch format {
	case FormatText:
		handler = slog.NewTextHandler(w, &opts)
	case FormatJSON:
		handler = slog.NewJSONHandler(w, &opts)
	default:
		handler = slog.NewJSONHandler(w, &opts)
	}

	return slog.New(handler)
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
