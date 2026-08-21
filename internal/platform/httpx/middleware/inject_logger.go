package middleware

import (
	"log/slog"
	"net/http"

	"github.com/iamroockie/parterre/internal/platform/logger"
)

func InjectLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logger.NewContext(r.Context(), log)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
