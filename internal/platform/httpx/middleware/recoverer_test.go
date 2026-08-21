package middleware_test

import (
	"encoding/json/jsontext"
	"net/http"
	"net/http/httptest"
	"testing"

	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/httpx/middleware"
	"github.com/iamroockie/parterre/internal/platform/httpx/middleware/mwtest"
)

func TestRecoverer(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		h := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			panic("test_panic")
		})
		log, buf := mwtest.NewTestLogger(t)
		r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		handler := middleware.InjectLogger(log)(middleware.Recoverer(h))

		handler.ServeHTTP(w, r)
		lines := mwtest.LogLines(t, buf)

		require.Equal(t, http.StatusInternalServerError, w.Code)
		require.Equal(t, "application/json", w.Result().Header.Get("Content-Type"))
		require.Equal(t, `{"error":"Internal server error"}`, w.Body.String())

		require.Len(t, lines, 1)
		msg := lines[0]
		require.Equal(t, "test_panic", msg["panic"])
		require.Contains(t, msg["stack"], "goroutine")
	})

	t.Run("ErrAbortHandler", func(t *testing.T) {
		defer func() {
			rvr := recover()
			require.NotNil(t, rvr)
			err, ok := rvr.(error)
			require.True(t, ok)
			require.ErrorIs(t, err, http.ErrAbortHandler)
		}()

		h := http.HandlerFunc(func(_ http.ResponseWriter, _ *http.Request) {
			panic(http.ErrAbortHandler)
		})
		r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		handler := middleware.Recoverer(h)

		handler.ServeHTTP(w, r)
	})

	t.Run("handler write header and panic", func(t *testing.T) {
		h := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"status": "success", "data": `))

			panic("test_panic")
		})
		r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)
		handler := middleware.Recoverer(h)

		handler.ServeHTTP(ww, r)

		require.Equal(t, http.StatusOK, w.Code)
		require.Empty(t, w.Result().Header.Get("Content-Type"))
		require.False(t, jsontext.Value(w.Body.String()).IsValid())
	})
}
