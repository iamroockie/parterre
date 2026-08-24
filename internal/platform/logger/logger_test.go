package logger_test

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/logger"
)

func TestJSON(t *testing.T) {
	var buf bytes.Buffer

	log := logger.New(&buf, slog.LevelInfo)
	log.Info("test")

	var rec map[string]any
	require.NoError(t, json.Unmarshal(buf.Bytes(), &rec))
}

func TestContext(t *testing.T) {
	t.Run("with logger", func(t *testing.T) {
		var buf bytes.Buffer
		log := logger.New(&buf, slog.LevelInfo)
		ctx := logger.NewContext(context.Background(), log)

		got := logger.FromContext(ctx)

		require.Same(t, log, got)
	})

	t.Run("without logger", func(t *testing.T) {
		got := logger.FromContext(context.Background())

		require.Same(t, slog.Default(), got)
	})
}
