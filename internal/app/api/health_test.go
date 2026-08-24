package api_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/app/api"
	"github.com/iamroockie/parterre/internal/app/api/middleware"
	"github.com/iamroockie/parterre/internal/platform/logger/loggertest"
	"github.com/iamroockie/parterre/internal/platform/postgres/postgrestest"
)

func TestHealth(t *testing.T) {
	r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	api.Health().ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.Empty(t, w.Body.String())
}

func TestPostgresCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping postgres check in short mode")
	}

	timeout := 20 * time.Millisecond

	t.Run("ok", func(t *testing.T) {
		pool := postgrestest.NewTestPool(t)
		r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		api.PostgresCheck(pool, timeout).ServeHTTP(w, r)

		require.Equal(t, http.StatusOK, w.Code)
		require.Empty(t, w.Body.String())
	})

	t.Run("unavailable", func(t *testing.T) {
		log, buf := loggertest.NewLogger(t, slog.LevelDebug)
		pool := postgrestest.NewTestPool(t)
		pool.Close()
		r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		handler := middleware.Logger(log)(api.PostgresCheck(pool, timeout))

		handler.ServeHTTP(w, r)

		logs := loggertest.Logs(t, buf)
		require.Len(t, logs, 2)
		require.Contains(t, logs[0]["msg"], "postgres check failed")
		require.Equal(t, http.StatusServiceUnavailable, w.Code)
		require.Empty(t, w.Body.String())
	})
}
