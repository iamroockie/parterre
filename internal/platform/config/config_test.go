package config_test

import (
	"log/slog"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/config"
)

func TestLoad(t *testing.T) {
	env := config.EnvLocal
	logLevel := slog.LevelDebug.String()
	postgresDSN := "postgres://"
	graceShutdown := 30 * time.Second
	httpPort := uint16(8080)
	httpHost := "127.0.0.1"
	httpReadHeaderTimeout := 5 * time.Second
	httpReadTimeout := 10 * time.Second
	httpWriteTimeout := 10 * time.Second
	httpIdleTimeout := 60 * time.Second

	t.Setenv("APP_ENV", string(env))
	t.Setenv("LOG_LEVEL", logLevel)
	t.Setenv("HTTP_PORT", strconv.FormatUint(uint64(httpPort), 10))
	t.Setenv("HTTP_HOST", httpHost)
	t.Setenv("HTTP_READ_TIMEOUT", httpReadTimeout.String())
	t.Setenv("HTTP_WRITE_TIMEOUT", httpWriteTimeout.String())
	t.Setenv("HTTP_READ_HEADER_TIMEOUT", httpReadHeaderTimeout.String())
	t.Setenv("HTTP_IDLE_TIMEOUT", httpIdleTimeout.String())
	t.Setenv("GRACEFUL_SHUTDOWN_TIMEOUT", graceShutdown.String())
	t.Setenv("POSTGRES_DSN", postgresDSN)

	t.Run("success", func(t *testing.T) {
		cfg, err := config.Load()

		require.NoError(t, err)
		require.Equal(t, env, cfg.Env)
		require.Equal(t, logLevel, cfg.LogLevel.String())
		require.Equal(t, postgresDSN, cfg.PostgresDSN)
		require.Equal(t, graceShutdown, cfg.GracefulShutdownTimeout)
		require.Equal(t, httpPort, cfg.HTTP.Port)
		require.Equal(t, httpHost, cfg.HTTP.Host)
		require.Equal(t, httpReadHeaderTimeout, cfg.HTTP.ReadHeaderTimeout)
		require.Equal(t, httpReadTimeout, cfg.HTTP.ReadTimeout)
		require.Equal(t, httpWriteTimeout, cfg.HTTP.WriteTimeout)
		require.Equal(t, httpIdleTimeout, cfg.HTTP.IdleTimeout)
	})

	t.Run("fails", func(t *testing.T) {
		tests := map[string]struct {
			key, value string
		}{
			"invalid ENV":                        {"APP_ENV", "asd"},
			"postgres DSN empty":                 {"POSTGRES_DSN", ""},
			"http port zero":                     {"HTTP_PORT", "0"},
			"http read header timeout negative":  {"HTTP_READ_HEADER_TIMEOUT", "-1s"},
			"http read header timeout zero":      {"HTTP_READ_HEADER_TIMEOUT", "0"},
			"http read timeout negative":         {"HTTP_READ_TIMEOUT", "-1s"},
			"http read timeout zero":             {"HTTP_READ_TIMEOUT", "0"},
			"http write timeout negative":        {"HTTP_WRITE_TIMEOUT", "-1s"},
			"http write timeout zero":            {"HTTP_WRITE_TIMEOUT", "0"},
			"http idle timeout negative":         {"HTTP_IDLE_TIMEOUT", "-1s"},
			"http idle timeout zero":             {"HTTP_IDLE_TIMEOUT", "0"},
			"graceful shutdown timeout negative": {"GRACEFUL_SHUTDOWN_TIMEOUT", "-1s"},
			"graceful shutdown timeout zero":     {"GRACEFUL_SHUTDOWN_TIMEOUT", "0"},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				t.Setenv(tt.key, tt.value)

				_, err := config.Load()

				require.Error(t, err)
				require.ErrorContains(t, err, tt.key)
			})
		}
	})
}
