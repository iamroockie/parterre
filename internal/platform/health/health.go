package health

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"slices"
	"time"

	"github.com/iamroockie/parterre/internal/platform/logger"
)

type Checker interface {
	Ping(context.Context) error
}

func Live() http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
}

func Ready(timeout time.Duration, checkers map[string]Checker) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), timeout)
		defer cancel()
		log := logger.FromContext(ctx)

		names, errs := pingCheckers(ctx, checkers)

		if len(names) > 0 {
			log.Warn("readiness check failed", "checks", names, "error", errors.Join(errs...))
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

type checkResult struct {
	name string
	err  error
}

func pingCheckers(ctx context.Context, checkers map[string]Checker) ([]string, []error) {
	results := make(chan checkResult, len(checkers))

	pending := make(map[string]struct{}, len(checkers))
	for name, ch := range checkers {
		pending[name] = struct{}{}

		go func() {
			results <- checkResult{name: name, err: ch.Ping(ctx)}
		}()
	}

	failed := make(map[string]error)

collect:
	for len(pending) > 0 {
		select {
		case res := <-results:
			delete(pending, res.name)

			if res.err != nil {
				failed[res.name] = fmt.Errorf("%s: %w", res.name, res.err)
			}
		case <-ctx.Done():
			for name := range pending {
				failed[name] = fmt.Errorf("%s: %w", name, ctx.Err())
			}

			break collect
		}
	}

	names := slices.Sorted(maps.Keys(failed))
	errs := make([]error, len(names))

	for i, n := range names {
		errs[i] = failed[n]
	}

	return names, errs
}
