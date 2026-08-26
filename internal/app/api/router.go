package api

import (
	"log/slog"

	"github.com/go-chi/chi/v5"

	"github.com/iamroockie/parterre/internal/platform/httpx/middleware"
)

func newRouter(log *slog.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.NotFound(notFound)
	r.MethodNotAllowed(methodNotAllowed)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger(log, "/healthz", "/readyz"))
	r.Use(middleware.Recover)

	return r
}
