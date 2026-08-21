package middleware

import (
	"context"
	"net/http"
	"uuid"

	"github.com/iamroockie/parterre/internal/platform/logger"
)

type reqIDCtxKey struct{}

const HeaderRequestID = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := resolveRequestID(r)

		log := logger.FromContext(r.Context()).With("request_id", id)

		ctx := context.WithValue(r.Context(), reqIDCtxKey{}, id)
		ctx = logger.NewContext(ctx, log)

		w.Header().Set(HeaderRequestID, id.String())

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RequestIDFromContext(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value(reqIDCtxKey{}).(uuid.UUID); ok {
		return id
	}

	return uuid.Nil()
}

func resolveRequestID(r *http.Request) uuid.UUID {
	reqID, err := uuid.Parse(r.Header.Get(HeaderRequestID))
	if err != nil || reqID == uuid.Nil() {
		return uuid.NewV4()
	}

	return reqID
}
