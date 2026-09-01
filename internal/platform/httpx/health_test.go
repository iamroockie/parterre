package httpx_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/httpx"
	"github.com/iamroockie/parterre/internal/platform/httpx/middleware"
	"github.com/iamroockie/parterre/internal/platform/logger/loggertest"
	"github.com/iamroockie/parterre/internal/platform/pg/pgtest"
)

func TestHealth(t *testing.T) {
	r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	httpx.Health(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.Empty(t, w.Body.String())
}

func TestPostgresCheck(t *testing.T) {
	timeout := 20 * time.Millisecond

	t.Run("ok", func(t *testing.T) {
		pool := pgtest.NewTestPool(t)
		r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		httpx.PostgresCheck(pool, timeout).ServeHTTP(w, r)

		require.Equal(t, http.StatusOK, w.Code)
		require.Empty(t, w.Body.String())
	})

	t.Run("unavailable", func(t *testing.T) {
		log, buf := loggertest.NewLogger(t, slog.LevelDebug)
		pool, err := pgxpool.New(t.Context(), "postgres://parterre@127.0.0.1:5432/parterre")
		require.NoError(t, err)
		pool.Close()
		r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		handler := middleware.Logger(log)(httpx.PostgresCheck(pool, timeout))

		handler.ServeHTTP(w, r)

		logs := loggertest.Logs(t, buf)
		require.Len(t, logs, 2)
		require.Contains(t, logs[0]["msg"], "postgres check failed")
		require.Equal(t, http.StatusServiceUnavailable, w.Code)
	})
}
