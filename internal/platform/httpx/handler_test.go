package httpx_test

import (
	"bytes"
	"encoding/json/v2"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/httpx"
	"github.com/iamroockie/parterre/internal/platform/httpx/httpxtest"
)

func TestHandler_TooLargeBody(t *testing.T) {
	calls := 0
	handler := func(_ http.ResponseWriter, r *http.Request) error {
		calls++
		var a any
		return json.UnmarshalRead(r.Body, &a)
	}
	body := strings.Repeat("a", httpx.MaxBodyBytes+1)
	data, err := json.Marshal(body)
	require.NoError(t, err)
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", bytes.NewReader(data))
	fn := httpx.Handle(handler)

	fn(w, r)

	require.Equal(t, 1, calls)
	require.Equal(t, http.StatusRequestEntityTooLarge, w.Code)
	resp := httpxtest.DecodeErrorResponse(t, w)
	require.Equal(t, httpx.CodeEntityTooLarge, resp.Error.Code)
}
