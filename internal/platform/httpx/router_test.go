package httpx_test

import (
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/httpx"
	"github.com/iamroockie/parterre/internal/platform/httpx/middleware/mwtest"
)

func TestRouter(t *testing.T) {
	tests := map[string]struct {
		path       string
		wantStatus int
		wantLevel  slog.Level
	}{
		"healthz": {
			path:       "/healthz",
			wantStatus: http.StatusOK,
			wantLevel:  slog.LevelInfo,
		},
		"readyz": {
			path:       "/readyz",
			wantStatus: http.StatusOK,
			wantLevel:  slog.LevelInfo,
		},

		"unknown-path": {
			path:       "/unknown-path",
			wantStatus: http.StatusNotFound,
			wantLevel:  slog.LevelWarn,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			log, buf := mwtest.NewTestLogger(t)
			router := httpx.NewRouter(log)
			r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, r)
			lines := mwtest.LogLines(t, buf)

			require.Equal(t, tt.wantStatus, w.Code)
			require.Len(t, lines, 1)
			require.Equal(t, tt.wantLevel.String(), lines[0]["level"])
		})
	}
}

func TestOrderMiddleware(t *testing.T) {
	h := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
		panic("test panic")
	})
	log, buf := mwtest.NewTestLogger(t)
	router := httpx.NewRouter(log)
	router.Get("/test-panic", h)
	r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/test-panic", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, r)
	lines := mwtest.LogLines(t, buf)

	require.Len(t, lines, 2)
	require.NotEqual(t, uuid.Nil().String(), lines[0]["request_id"])
	require.NotEqual(t, uuid.Nil().String(), lines[1]["request_id"])
	require.Equal(t, lines[0]["request_id"], lines[1]["request_id"])
}
