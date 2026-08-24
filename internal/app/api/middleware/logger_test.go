package middleware_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/app/api/middleware"
	"github.com/iamroockie/parterre/internal/platform/logger/loggertest"
)

func TestLogger(t *testing.T) {
	sleep := 10 * time.Millisecond
	h := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		time.Sleep(sleep)
	})
	log, buf := loggertest.NewLogger(t, slog.LevelDebug)
	r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()
	handler := middleware.RequestID(middleware.Logger(log)(h))

	handler.ServeHTTP(w, r)

	logs := loggertest.Logs(t, buf)
	require.Len(t, logs, 1)
	msg := logs[0]
	reqID := w.Result().Header.Get(middleware.HeaderRequestID)
	require.Equal(t, reqID, msg["request_id"])
	require.Equal(t, http.MethodGet, msg["method"])
	require.Equal(t, "/test", msg["path"])
	require.GreaterOrEqual(t, msg["duration_ms"], 10.0)
	require.Less(t, msg["duration_ms"], 1000.0)
}
