package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/iamroockie/parterre/internal/platform/config"
)

type routesDeps struct {
	cfg  config.Config
	pool *pgxpool.Pool
}

func routes(r chi.Router, deps routesDeps) {
	r.Get("/healthz", Health())
	r.Get("/readyz", PostgresCheck(deps.pool, deps.cfg.HTTP.ReadyTimeout))
}
