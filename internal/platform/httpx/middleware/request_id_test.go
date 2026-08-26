package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/httpx/middleware"
)

func TestRequestID(t *testing.T) {
	tests := map[string]struct {
		id            string
		setHeader     bool
		wantInherited bool
	}{
		"valid id": {
			id:            uuid.NewV7().String(),
			setHeader:     true,
			wantInherited: true,
		},
		"without id": {
			id:            "",
			setHeader:     false,
			wantInherited: false,
		},
		"empty id": {
			id:            "",
			setHeader:     true,
			wantInherited: false,
		},
		"invalid id": {
			id:            "invalid id",
			setHeader:     true,
			wantInherited: false,
		},
		"zero id": {
			id:            uuid.Nil().String(),
			setHeader:     true,
			wantInherited: false,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var idFromCtx uuid.UUID
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				idFromCtx = middleware.GetRequestID(r.Context())
				w.WriteHeader(http.StatusOK)
			})
			r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
			if tt.setHeader {
				r.Header.Set(middleware.HeaderRequestID, tt.id)
			}
			w := httptest.NewRecorder()
			handler := middleware.RequestID(h)

			handler.ServeHTTP(w, r)

			id, err := uuid.Parse(w.Result().Header.Get(middleware.HeaderRequestID))
			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil(), id)
			require.Equal(t, idFromCtx, id)
			if tt.wantInherited {
				require.Equal(t, tt.id, id.String())
			} else {
				require.NotEqual(t, tt.id, id.String())
			}
		})
	}
}

func TestContextWithoutRequestID(t *testing.T) {
	ctx := context.Background()

	reqID := middleware.GetRequestID(ctx)

	require.Equal(t, uuid.Nil(), reqID)
}
