package httpxtest

import (
	"encoding/json/v2"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/httpx"
)

type ErrorResponse struct {
	Error struct {
		Code    httpx.Code        `json:"code"`
		Message string            `json:"message"`
		Errors  map[string]string `json:"errors,omitempty"`
	} `json:"error"`
}

func DecodeErrorResponse(t testing.TB, w *httptest.ResponseRecorder) ErrorResponse {
	t.Helper()

	var resp ErrorResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))

	return resp
}
