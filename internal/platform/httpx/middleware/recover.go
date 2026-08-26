package middleware

import (
	"errors"
	"net/http"
	"runtime/debug"

	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/iamroockie/parterre/internal/platform/httpx/response"
	"github.com/iamroockie/parterre/internal/platform/logger"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ww, ok := w.(chimw.WrapResponseWriter)
		if !ok {
			ww = chimw.NewWrapResponseWriter(w, r.ProtoMajor)
		}

		defer func() {
			if rvr := recover(); rvr != nil {
				if err, ok := rvr.(error); ok && errors.Is(err, http.ErrAbortHandler) {
					panic(rvr)
				}

				log := logger.FromContext(r.Context())
				log.Error("recovered from panic", "panic", rvr, "stack", string(debug.Stack()))

				if ww.Status() != 0 {
					return
				}

				response.Error(ww, http.StatusInternalServerError, "Internal error")
			}
		}()

		next.ServeHTTP(ww, r)
	})
}
