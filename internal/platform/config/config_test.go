package config_test

import (
	"log/slog"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/config"
	"github.com/iamroockie/parterre/internal/platform/config/configtest"
)

func TestLoad_Default(t *testing.T) {
	cfg := configtest.Default(t)
	cfg.AppEnv = config.EnvLocal
	cfg.Postgres.DSN = "pg//"
	t.Setenv("APP_ENV", cfg.AppEnv.String())
	t.Setenv("LOG_LEVEL", cfg.LogLevel.String())
	t.Setenv("HTTP_PORT", strconv.FormatUint(uint64(cfg.HTTP.Port), 10))
	t.Setenv("HTTP_IDLE_TIMEOUT", cfg.HTTP.IdleTimeout.String())
	t.Setenv("HTTP_READY_TIMEOUT", cfg.HTTP.ReadyTimeout.String())
	t.Setenv("HTTP_SHUTDOWN_TIMEOUT", cfg.HTTP.ShutdownTimeout.String())
	t.Setenv("POSTGRES_DSN", cfg.Postgres.DSN)
	t.Setenv("POSTGRES_MAX_CONNS", strconv.FormatInt(int64(cfg.Postgres.MaxConns), 10))

	got, err := config.Load()

	require.NoError(t, err)
	require.Equal(t, cfg, got)
}

func TestLoad_Custom(t *testing.T) {
	t.Setenv("APP_ENV", "prod")
	t.Setenv("LOG_LEVEL", "debug")
	t.Setenv("HTTP_PORT", "3000")
	t.Setenv("HTTP_IDLE_TIMEOUT", "10ms")
	t.Setenv("HTTP_READY_TIMEOUT", "10ms")
	t.Setenv("HTTP_SHUTDOWN_TIMEOUT", "10ms")
	t.Setenv("POSTGRES_DSN", "pg//test")
	t.Setenv("POSTGRES_MAX_CONNS", "4")

	cfg, err := config.Load()

	require.NoError(t, err)
	require.Equal(t, slog.LevelDebug, cfg.LogLevel)
	require.Equal(t, config.EnvProd, cfg.AppEnv)
	require.Equal(t, uint16(3000), cfg.HTTP.Port)
	require.Equal(t, 10*time.Millisecond, cfg.HTTP.IdleTimeout)
	require.Equal(t, 10*time.Millisecond, cfg.HTTP.ReadyTimeout)
	require.Equal(t, 10*time.Millisecond, cfg.HTTP.ShutdownTimeout)
	require.Equal(t, "pg//test", cfg.Postgres.DSN)
	require.Equal(t, int32(4), cfg.Postgres.MaxConns)
}

func TestLoad_FailParse(t *testing.T) {
	t.Setenv("APP_ENV", "prod")
	t.Setenv("POTGRES_DSN", "pg//")
	tests := map[string]struct {
		key, value string
	}{
		"invalid APP_ENV":   {"APP_ENV", "unknown"},
		"invalid LOG_LEVEL": {"LOG_LEVEL", "unknown"},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			t.Setenv(tt.key, tt.value)

			_, err := config.Load()

			require.Error(t, err)
			require.ErrorContains(t, err, "parse config")
		})
	}
}
