package health_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/health"
	"github.com/iamroockie/parterre/internal/platform/httpx/middleware"
	"github.com/iamroockie/parterre/internal/platform/httpx/middleware/mwtest"
)

type checker func(ctx context.Context) error

func (c checker) Ping(ctx context.Context) error {
	return c(ctx)
}

func TestLive(t *testing.T) {
	r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	health.Live().ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.Empty(t, w.Body.String())
}

func TestReady(t *testing.T) {
	timeout := 20 * time.Millisecond

	t.Run("empty checkers", func(t *testing.T) {
		r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		health.Ready(timeout, map[string]health.Checker{}).ServeHTTP(w, r)

		require.Equal(t, http.StatusOK, w.Code)
		require.Empty(t, w.Body.String())
	})

	t.Run("all checkers ready", func(t *testing.T) {
		r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		checkers := map[string]health.Checker{
			"first_service": checker(func(_ context.Context) error {
				return nil
			}),
			"second_service": checker(func(_ context.Context) error {
				return nil
			}),
		}

		health.Ready(timeout, checkers).ServeHTTP(w, r)

		require.Equal(t, http.StatusOK, w.Code)
		require.Empty(t, w.Body.String())
	})

	t.Run("checker timed out", func(t *testing.T) {
		log, _ := mwtest.NewTestLogger(t)
		ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
		defer cancel()
		r := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		checkers := map[string]health.Checker{
			"slow_service": checker(func(ctx context.Context) error {
				<-ctx.Done()
				return errors.New("unavailable")
			}),
		}
		handler := middleware.InjectLogger(log)(health.Ready(timeout, checkers))

		start := time.Now()
		handler.ServeHTTP(w, r)
		elapsed := time.Since(start)

		require.Equal(t, http.StatusServiceUnavailable, w.Code)
		require.Empty(t, w.Body.String())
		require.Less(t, elapsed, 500*time.Millisecond)
	})

	t.Run("checker ignores context", func(t *testing.T) {
		log, buf := mwtest.NewTestLogger(t)
		ctx, cancel := context.WithTimeout(t.Context(), 2*time.Second)
		defer cancel()
		r := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		release := make(chan struct{})
		t.Cleanup(func() { close(release) })

		checkers := map[string]health.Checker{
			"rude_service": checker(func(_ context.Context) error {
				<-release

				return nil
			}),
		}
		handler := middleware.InjectLogger(log)(health.Ready(timeout, checkers))

		start := time.Now()
		handler.ServeHTTP(w, r)
		elapsed := time.Since(start)
		lines := mwtest.LogLines(t, buf)

		require.Equal(t, http.StatusServiceUnavailable, w.Code)
		require.Less(t, elapsed, 500*time.Millisecond)
		require.Len(t, lines, 1)
		require.Equal(t, []any{"rude_service"}, lines[0]["checks"])
		require.Contains(t, lines[0]["error"], context.DeadlineExceeded.Error())
	})

	t.Run("checkers unavailabled", func(t *testing.T) {
		log, buf := mwtest.NewTestLogger(t)
		r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		checkers := map[string]health.Checker{
			"first_service": checker(func(_ context.Context) error {
				return nil
			}),
			"second_service": checker(func(_ context.Context) error {
				return errors.New("unavailable")
			}),
			"third_service": checker(func(_ context.Context) error {
				return errors.New("unavailable")
			}),
		}
		handler := middleware.InjectLogger(log)(health.Ready(timeout, checkers))

		handler.ServeHTTP(w, r)
		lines := mwtest.LogLines(t, buf)

		require.Len(t, lines, 1)
		require.Equal(t, []any{"second_service", "third_service"}, lines[0]["checks"])
		require.Equal(t, http.StatusServiceUnavailable, w.Code)
		require.Empty(t, w.Body.String())
	})
}
