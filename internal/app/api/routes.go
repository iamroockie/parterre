package api

import (
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/iamroockie/parterre/internal/platform/health"
)

type deps struct {
	readyTimeout time.Duration
	checkers     map[string]health.Checker
}

func routes(r chi.Router, d deps) {
	r.Get("/healthz", health.Live())
	r.Get("/readyz", health.Ready(d.readyTimeout, d.checkers))
}
