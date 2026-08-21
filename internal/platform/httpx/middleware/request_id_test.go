package middleware_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/httpx/middleware"
	"github.com/iamroockie/parterre/internal/platform/httpx/middleware/mwtest"
)

func TestRequestID(t *testing.T) {
	tests := map[string]struct {
		id            string
		setHeader     bool
		wantInherited bool
	}{
		"without request_id": {
			id:            "",
			setHeader:     false,
			wantInherited: false,
		},
		"empty in header": {
			id:            "",
			setHeader:     true,
			wantInherited: false,
		},
		"invalid request_id": {
			id:            "invalid ID",
			setHeader:     true,
			wantInherited: false,
		},
		"zero uuid in header": {
			id:            uuid.Nil().String(),
			setHeader:     true,
			wantInherited: false,
		},
		"valid request_id": {
			id:            uuid.NewV4().String(),
			setHeader:     true,
			wantInherited: true,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			var fromCtx uuid.UUID
			h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fromCtx = middleware.RequestIDFromContext(r.Context())
				w.WriteHeader(http.StatusOK)
			})
			log, buf := mwtest.NewTestLogger(t)
			r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
			if tt.setHeader {
				r.Header.Set(middleware.HeaderRequestID, tt.id)
			}
			w := httptest.NewRecorder()
			handler := middleware.InjectLogger(log)(middleware.RequestID(middleware.Logger(h)))

			handler.ServeHTTP(w, r)
			reqID := w.Result().Header.Get(middleware.HeaderRequestID)
			id, err := uuid.Parse(reqID)
			lines := mwtest.LogLines(t, buf)

			require.NoError(t, err)
			require.NotEqual(t, uuid.Nil(), id)
			require.Equal(t, id, fromCtx)
			require.Len(t, lines, 1)
			require.Equal(t, id.String(), lines[0]["request_id"])
			if tt.wantInherited {
				require.Equal(t, tt.id, id.String())
				require.Equal(t, tt.id, lines[0]["request_id"])
			} else {
				require.NotEqual(t, tt.id, id.String())
			}
		})
	}
}

func TestContextWithoutRequestID(t *testing.T) {
	ctx := context.Background()

	reqID := middleware.RequestIDFromContext(ctx)

	require.Equal(t, uuid.Nil(), reqID)
}
