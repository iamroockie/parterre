package api

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/iamroockie/parterre/internal/platform/httpx/response"
)

func notFound(w http.ResponseWriter, _ *http.Request) {
	response.Error(w, http.StatusNotFound, "Not found")
}

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	if allow := allowedMethods(r); allow != "" {
		w.Header().Set("Allow", allow)
	}
	response.Error(w, http.StatusMethodNotAllowed, "Method not allowed")
}

func allowedMethods(r *http.Request) string {
	rctx := chi.RouteContext(r.Context())
	if rctx == nil {
		return ""
	}

	all := []string{
		http.MethodGet, http.MethodHead, http.MethodPost, http.MethodPut,
		http.MethodPatch, http.MethodDelete, http.MethodOptions,
	}

	methods := make([]string, 0, len(all))
	for _, m := range all {
		if rctx.Routes.Match(chi.NewRouteContext(), m, r.URL.Path) {
			methods = append(methods, m)
		}
	}

	return strings.Join(methods, ", ")
}
