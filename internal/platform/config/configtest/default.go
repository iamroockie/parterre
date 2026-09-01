package configtest

import (
	"log/slog"
	"testing"
	"time"

	"github.com/iamroockie/parterre/internal/platform/config"
)

func Default(t testing.TB) config.Config {
	t.Helper()

	return config.Config{
		AppEnv:   config.EnvLocal,
		LogLevel: slog.LevelInfo,

		HTTP: config.HTTPConfig{
			Port:            8080,
			IdleTimeout:     1 * time.Minute,
			ReadyTimeout:    2 * time.Second,
			ShutdownTimeout: 20 * time.Second,
		},

		// nolint:exhaustruct_v5
		Postgres: config.PostgresConfig{
			MaxConns: 6,
		},
	}
}
