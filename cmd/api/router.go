package main

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/iamroockie/parterre/internal/catalog"
	catalogPg "github.com/iamroockie/parterre/internal/catalog/infra/postgres"
	catalogRest "github.com/iamroockie/parterre/internal/catalog/transport/rest"
	catalogUC "github.com/iamroockie/parterre/internal/catalog/usecase"
	"github.com/iamroockie/parterre/internal/platform/config"
	"github.com/iamroockie/parterre/internal/platform/httpx"
	"github.com/iamroockie/parterre/internal/platform/httpx/middleware"
)

type dependencies struct {
	cfg  config.Config
	log  *slog.Logger
	pool *pgxpool.Pool
}

func router(deps dependencies) *chi.Mux {
	r := chi.NewRouter()

	r.NotFound(notFoundHandleFunc)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(deps.log, "/healthz", "/readyz"))
	r.Use(middleware.Recover)

	venueStore := catalogPg.NewVenueStore(deps.pool)
	hallStore := catalogPg.NewHallStore(deps.pool)

	catalogMod := catalog.Module{
		GetVenue:    catalogUC.NewGetVenue(venueStore).Execute,
		CreateVenue: catalogUC.NewCreateVenue(venueStore).Execute,
		GetHall:     catalogUC.NewGetHall(hallStore).Execute,
		CreateHall:  catalogUC.NewCreateHall(hallStore).Execute,
	}

	r.Get("/healthz", httpx.Health)
	r.Get("/readyz", httpx.PostgresCheck(deps.pool, deps.cfg.HTTP.ReadyTimeout))

	r.Route("/api", func(r chi.Router) {
		catalogRest.RegisterRoutes(r, catalogMod)
	})

	return r
}

func notFoundHandleFunc(w http.ResponseWriter, r *http.Request) {
	httpx.WriteError(w, r, httpx.RouteNotFoundError(r.URL.Path))
}
