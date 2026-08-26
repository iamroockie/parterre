package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/iamroockie/parterre/internal/catalog"
	"github.com/iamroockie/parterre/internal/platform/config"
)

type routesDeps struct {
	cfg  config.Config
	pool *pgxpool.Pool
}

func routes(r chi.Router, deps routesDeps) {
	r.Get("/healthz", Health())
	r.Get("/readyz", PostgresCheck(deps.pool, deps.cfg.HTTP.ReadyTimeout))

	venueStore := catalog.NewVenueStore(deps.pool)
	venueHandler := catalog.NewVenueHandler(venueStore)

	r.Route("/v1/venues", func(r chi.Router) {
		r.Post("/", venueHandler.Create)
		r.Get("/{id}", venueHandler.Get)
	})
}
