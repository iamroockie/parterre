package middleware

import (
	"log/slog"
	"net/http"
	"time"

	chimw "github.com/go-chi/chi/v5/middleware"

	"github.com/iamroockie/parterre/internal/platform/logger"
)

func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		now := time.Now()

		ww := chimw.NewWrapResponseWriter(w, r.ProtoMajor)

		next.ServeHTTP(ww, r)

		elapsed := time.Since(now)
		log := logger.FromContext(r.Context())

		status := ww.Status()
		if status == 0 {
			status = http.StatusOK
		}

		var level slog.Level
		switch {
		case status < 400:
			level = slog.LevelInfo
		case status < 500:
			level = slog.LevelWarn
		default:
			level = slog.LevelError
		}

		log.Log(
			r.Context(),
			level,
			"http request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", status,
			"duration_ms", float64(elapsed.Microseconds())/1000)
	})
}
