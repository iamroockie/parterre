package middleware

import (
	"context"
	"net/http"
	"uuid"

	"github.com/iamroockie/parterre/internal/platform/identity"
)

type ctxKeyRequestID struct{}

const HeaderRequestID = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id, err := uuid.Parse(r.Header.Get(HeaderRequestID))
		if err != nil || id == uuid.Nil() {
			id = identity.NewUUID()
		}
		w.Header().Set(HeaderRequestID, id.String())
		ctx := context.WithValue(r.Context(), ctxKeyRequestID{}, id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetRequestID(ctx context.Context) uuid.UUID {
	if id, ok := ctx.Value(ctxKeyRequestID{}).(uuid.UUID); ok {
		return id
	}

	return uuid.Nil()
}
