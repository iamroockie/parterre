package middleware

import (
	"log/slog"
	"net/http"
	"slices"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/iamroockie/parterre/internal/platform/logger"
)

func Logger(baseLogger *slog.Logger, quiet ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqID := GetRequestID(r.Context())
			log := baseLogger.With(
				"request_id", reqID,
				"method", r.Method,
				"path", r.URL.Path,
			)
			ctx := logger.NewContext(r.Context(), log)
			r = r.WithContext(ctx)
			ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)

			now := time.Now()
			next.ServeHTTP(ww, r)
			elapsed := time.Since(now)

			status := ww.Status()
			if status == 0 {
				status = http.StatusOK
			}
			if status != http.StatusOK || !slices.Contains(quiet, r.URL.Path) {
				log.Info("http request",
					"status", status,
					"duration_ms", float64(elapsed.Microseconds())/1000,
				)
			}
		})
	}
}
