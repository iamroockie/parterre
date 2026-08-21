package health_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/health"
)

type checker func(ctx context.Context) error

func (c checker) Ping(ctx context.Context) error {
	return c(ctx)
}

func TestHealth(t *testing.T) {
	r := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	health.Health().ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	require.Empty(t, w.Body)
}

func TestReady(t *testing.T) {
	t.Run("all checkers return nil", func(t *testing.T) {
		r := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		checkers := []health.Checker{
			checker(func(_ context.Context) error {
				return nil
			}),
			checker(func(_ context.Context) error {
				return nil
			}),
		}

		health.Ready(checkers).ServeHTTP(w, r)

		require.Equal(t, http.StatusOK, w.Code)
		require.Empty(t, w.Body)
	})

	t.Run("one checker return error", func(t *testing.T) {
		r := httptest.NewRequestWithContext(context.Background(), http.MethodGet, "/", nil)
		w := httptest.NewRecorder()
		checkers := []health.Checker{
			checker(func(_ context.Context) error {
				return nil
			}),
			checker(func(_ context.Context) error {
				return errors.New("unavailable")
			}),
			checker(func(_ context.Context) error {
				return nil
			}),
		}

		health.Ready(checkers).ServeHTTP(w, r)

		require.Equal(t, http.StatusServiceUnavailable, w.Code)
		require.Empty(t, w.Body)
	})
}
