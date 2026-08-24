package api

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/logger/loggertest"
)

func TestOrderMiddleware(t *testing.T) {
	h := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic("test panic")
	})
	log, buf := loggertest.NewLogger(t, slog.LevelDebug)
	router := newRouter(log)
	router.Get("/test-panic", h)
	r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/test-panic", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)

	logs := loggertest.Logs(t, buf)
	require.Len(t, logs, 2)
	require.NotEqual(t, uuid.Nil().String(), logs[0]["request_id"])
	require.NotEqual(t, uuid.Nil().String(), logs[1]["request_id"])
	require.Equal(t, logs[0]["request_id"], logs[1]["request_id"])
}
