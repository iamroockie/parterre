package middleware_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/httpx/middleware"
	"github.com/iamroockie/parterre/internal/platform/httpx/middleware/mwtest"
)

func TestLogger(t *testing.T) {
	tests := map[string]struct {
		method      string
		path        string
		status      int
		level       slog.Level
		writeHeader bool
	}{
		"no handler set": {
			method:      http.MethodGet,
			path:        "/",
			status:      http.StatusOK,
			level:       slog.LevelInfo,
			writeHeader: false,
		},
		"2xx -> info": {
			method:      http.MethodPost,
			path:        "/users",
			status:      http.StatusCreated,
			level:       slog.LevelInfo,
			writeHeader: true,
		},
		"4xx -> warn": {
			method:      http.MethodDelete,
			path:        "/users/1",
			status:      http.StatusNotFound,
			level:       slog.LevelWarn,
			writeHeader: true,
		},
		"5xx -> error": {
			method:      http.MethodPatch,
			path:        "/users/1",
			status:      http.StatusServiceUnavailable,
			level:       slog.LevelError,
			writeHeader: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			h := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				time.Sleep(10 * time.Millisecond)
				if tt.writeHeader {
					w.WriteHeader(tt.status)
				}
			})
			log, buf := mwtest.NewTestLogger(t)
			r := httptest.NewRequestWithContext(t.Context(), tt.method, tt.path, nil)
			w := httptest.NewRecorder()
			handler := middleware.InjectLogger(log)(middleware.RequestID(middleware.Logger(h)))

			handler.ServeHTTP(w, r)
			lines := mwtest.LogLines(t, buf)

			require.Len(t, lines, 1)
			msg := lines[0]
			require.Equal(t, tt.level.String(), msg["level"])
			require.Equal(t, tt.method, msg["method"])
			require.Equal(t, tt.path, msg["path"])
			require.Equal(t, float64(tt.status), msg["status"])
			require.GreaterOrEqual(t, msg["duration_ms"], 10.0)
			require.Less(t, msg["duration_ms"], 1000.0)
			require.NotNil(t, msg["request_id"])
		})
	}
}
