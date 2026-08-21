package httpx_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/httpx"
	"github.com/iamroockie/parterre/internal/platform/httpx/middleware/mwtest"
)

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
