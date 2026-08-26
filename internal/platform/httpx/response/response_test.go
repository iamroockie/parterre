package response_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/httpx/response"
)

func TestJSON(t *testing.T) {
	w := httptest.NewRecorder()

	response.JSON(w, http.StatusCreated, map[string]string{"id": "1"})

	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	require.JSONEq(t, `{"id":"1"}`, w.Body.String())
}

func TestJSONNotSerializable(t *testing.T) {
	w := httptest.NewRecorder()

	response.JSON(w, http.StatusCreated, make(chan int))

	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	require.JSONEq(t, `{"error":"`+response.InternalErrorMsg+`"}`, w.Body.String())
}

func TestError(t *testing.T) {
	w := httptest.NewRecorder()

	response.Error(w, http.StatusNotFound, "Venue not found")

	require.Equal(t, http.StatusNotFound, w.Code)
	require.JSONEq(t, `{"error":"Venue not found"}`, w.Body.String())
	require.NotContains(t, w.Body.String(), "fields")
}

func TestErrorWithFields(t *testing.T) {
	w := httptest.NewRecorder()

	fields := map[string]string{"name": "must not be empty"}

	response.ErrorWithFields(w, http.StatusUnprocessableEntity, "Validation error", fields)

	require.Equal(t, http.StatusUnprocessableEntity, w.Code)
	require.JSONEq(
		t,
		`{"error":"Validation error","fields":{"name":"must not be empty"}}`,
		w.Body.String(),
	)
}

func TestErrorInternal(t *testing.T) {
	w := httptest.NewRecorder()

	response.ErrorInternal(w)

	require.Equal(t, http.StatusInternalServerError, w.Code)
	require.JSONEq(t, `{"error":"`+response.InternalErrorMsg+`"}`, w.Body.String())
}
