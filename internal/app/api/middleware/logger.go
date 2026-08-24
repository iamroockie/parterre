package middleware

import (
	"log/slog"
	"net/http"
	"time"

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

			now := time.Now()
			next.ServeHTTP(w, r)
			elapsed := time.Since(now)

			log.Info("http request",
				"duration_ms", float64(elapsed.Microseconds())/1000,
			)
		})
	}
}
