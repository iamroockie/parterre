package api

import (
	"log/slog"

	"github.com/go-chi/chi/v5"

	"github.com/iamroockie/parterre/internal/app/api/middleware"
)

func newRouter(log *slog.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(log))
	r.Use(middleware.Recover)

	return r
}
