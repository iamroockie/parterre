package httpx

import (
	"log/slog"

	"github.com/go-chi/chi/v5"

	"github.com/iamroockie/parterre/internal/platform/httpx/middleware"
)

func NewRouter(log *slog.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.InjectLogger(log))
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	return r
}
