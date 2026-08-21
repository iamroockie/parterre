package middleware

import (
	"encoding/json/v2"
	"errors"
	"net/http"
	"runtime/debug"

	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/iamroockie/parterre/internal/platform/logger"
)

func Recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rvr := recover(); rvr != nil {
				if err, ok := rvr.(error); ok && errors.Is(err, http.ErrAbortHandler) {
					panic(rvr)
				}

				log := logger.FromContext(r.Context())
				log.Error("recovered from panic", "panic", rvr, "stack", string(debug.Stack()))

				if ww, ok := w.(chimw.WrapResponseWriter); ok && ww.Status() != 0 {
					return
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				jsonBody, _ := json.Marshal(map[string]string{
					"error": "Internal server error",
				})
				_, _ = w.Write(jsonBody)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
