package configtest

import (
	"log/slog"
	"testing"
	"time"

	"github.com/iamroockie/parterre/internal/platform/config"
)

func Config(t testing.TB) config.Config {
	t.Helper()

	return config.Config{
		Env:                     config.EnvLocal,
		LogLevel:                slog.LevelDebug,
		PostgresDSN:             "postgres://",
		GracefulShutdownTimeout: 10 * time.Second,
		HTTP: config.HTTPConfig{
			Host:              "127.0.0.1",
			Port:              10000,
			ReadyTimeout:      2 * time.Second,
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       5 * time.Second,
			WriteTimeout:      5 * time.Second,
			IdleTimeout:       10 * time.Second,
		},
	}
}
