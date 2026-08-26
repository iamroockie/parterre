package httpx_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/httpx"
	"github.com/iamroockie/parterre/internal/platform/httpx/middleware"
	"github.com/iamroockie/parterre/internal/platform/httpx/response"
	"github.com/iamroockie/parterre/internal/platform/logger/loggertest"
	"github.com/iamroockie/parterre/internal/platform/validation"
)

func serve(t *testing.T, h http.HandlerFunc) (*httptest.ResponseRecorder, []map[string]any) {
	t.Helper()

	log, buf := loggertest.NewLogger(t, slog.LevelDebug)
	r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	middleware.Logger(log)(h).ServeHTTP(w, r)

	return w, loggertest.Logs(t, buf)
}

func decodeError(t *testing.T, w *httptest.ResponseRecorder) (string, map[string]string) {
	t.Helper()

	var body struct {
		Error  string            `json:"error"`
		Fields map[string]string `json:"fields"`
	}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &body))

	return body.Error, body.Fields
}

func validationError() error {
	var b validation.Builder
	b.Add("name", "must not be empty")
	b.Add("timezone", "invalid timezone")

	return b.Err()
}

func TestWriteErrorValidation(t *testing.T) {
	w, logs := serve(t, func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteError(w, r, validationError())
	})

	require.Equal(t, http.StatusUnprocessableEntity, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))

	msg, fields := decodeError(t, w)
	require.Equal(t, "Validation error", msg)
	require.Equal(t, map[string]string{
		"name":     "must not be empty",
		"timezone": "invalid timezone",
	}, fields)

	require.Len(t, logs, 1)
	require.Equal(t, "http request", logs[0]["msg"])
}

func TestWriteErrorValidationWrapped(t *testing.T) {
	w, _ := serve(t, func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteError(w, r, fmt.Errorf("create venue: %w", validationError()))
	})

	require.Equal(t, http.StatusUnprocessableEntity, w.Code)

	_, fields := decodeError(t, w)
	require.Len(t, fields, 2)
}

func TestWriteErrorInternal(t *testing.T) {
	w, logs := serve(t, func(w http.ResponseWriter, r *http.Request) {
		httpx.WriteError(w, r, fmt.Errorf("create venue: %w", errors.New("boom")))
	})

	require.Equal(t, http.StatusInternalServerError, w.Code)

	require.JSONEq(t, `{"error":"`+response.InternalErrorMsg+`"}`, w.Body.String())
	require.NotContains(t, w.Body.String(), "boom")
	require.NotContains(t, w.Body.String(), "create venue")

	require.Len(t, logs, 2)
	require.Equal(t, "request failed", logs[0]["msg"])
	require.Contains(t, logs[0]["error"], "create venue")
	require.Contains(t, logs[0]["error"], "boom")
}
