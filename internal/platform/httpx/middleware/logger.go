package middleware

import (
	"log/slog"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/iamroockie/parterre/internal/platform/logger"
)

func Logger(baseLogger *slog.Logger) func(http.Handler) http.Handler {
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
			log.Info("http request",
				"status", status,
				"duration_ms", float64(elapsed.Microseconds())/1000,
			)
		})
	}
}
