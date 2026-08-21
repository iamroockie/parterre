package logger_test

import (
	"bytes"
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"log/slog"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/logger"
)

func TestNew(t *testing.T) {
	t.Run("json format", func(t *testing.T) {
		var buf bytes.Buffer

		log := logger.New(&buf, logger.FormatJSON, slog.LevelInfo)
		log.Info("info msg")
		log.Warn("warn msg")

		msgs := strings.Split(strings.TrimSpace(buf.String()), "\n")
		require.Len(t, msgs, 2)

		tests := []struct{ level, msg string }{
			{"INFO", "info msg"},
			{"WARN", "warn msg"},
		}

		for i, tt := range tests {
			var rec map[string]any
			require.NoError(t, json.Unmarshal([]byte(msgs[i]), &rec))
			require.Equal(t, tt.level, rec["level"])
			require.Equal(t, tt.msg, rec["msg"])
		}
	})

	t.Run("debug below level info", func(t *testing.T) {
		var buf bytes.Buffer

		log := logger.New(&buf, logger.FormatText, slog.LevelInfo)
		log.Debug("hidden message")

		require.Empty(t, buf.String())
	})

	t.Run("text format", func(t *testing.T) {
		var buf bytes.Buffer

		log := logger.New(&buf, logger.FormatText, slog.LevelInfo)
		log.Info("message")

		require.False(t, jsontext.Value(buf.Bytes()).IsValid())
		require.Contains(t, buf.String(), "level=INFO")
	})

	t.Run("unknown format", func(t *testing.T) {
		var buf bytes.Buffer

		log := logger.New(&buf, "invalid format", slog.LevelInfo)
		require.NotNil(t, log)

		log.Info("test message")

		require.True(t, jsontext.Value(buf.Bytes()).IsValid())
	})

	t.Run("level filter", func(t *testing.T) {
		var buf bytes.Buffer

		log := logger.New(&buf, logger.FormatJSON, slog.LevelWarn)
		log.Info("no write")
		require.Zero(t, buf.Len())

		log.Warn("write")
		require.NotZero(t, buf.Len())
	})
}

func TestContext(t *testing.T) {
	t.Run("logger from context", func(t *testing.T) {
		var buf bytes.Buffer
		ctx := context.Background()

		log := logger.New(&buf, logger.FormatText, slog.LevelInfo)
		log.Info("test message")

		ctxWithLogger := logger.NewContext(ctx, log)

		require.Same(t, log, logger.FromContext(ctxWithLogger))
	})

	t.Run("empty context", func(t *testing.T) {
		ctx := context.Background()

		log := logger.FromContext(ctx)

		require.Same(t, slog.Default(), log)
	})
}
