package httpx

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/iamroockie/parterre/internal/platform/httpx/middleware"
)

func NewRouter(log *slog.Logger) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.InjectLogger(log))
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	r.Get("/readyz", func(w http.ResponseWriter, _ *http.Request) {
		// TODO: Check DB
		w.WriteHeader(http.StatusOK)
	})

	return r
}
