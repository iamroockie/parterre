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
		LogLevel:        slog.LevelInfo,
		ShutdownTimeout: 30 * time.Second,

		HTTP: config.HTTPConfig{
			Port:         8080,
			IdleTimeout:  1 * time.Minute,
			ReadyTimeout: 2 * time.Second,
		},

		Postgres: config.PostgresConfig{
			MaxConns: 6,
		},
	}
}
