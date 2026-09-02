package rest_test

import (
	"context"
	"encoding/json/v2"
	"net/http"
	"net/http/httptest"
	"testing"
	"uuid"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/catalog"
	"github.com/iamroockie/parterre/internal/catalog/catalogtest"
	"github.com/iamroockie/parterre/internal/catalog/transport/rest"
	"github.com/iamroockie/parterre/internal/platform/httpx"
	"github.com/iamroockie/parterre/internal/platform/httpx/httpxtest"
	"github.com/iamroockie/parterre/internal/platform/identity"
)

func TestGetHall(t *testing.T) {
	calls := 0
	hall := catalogtest.Hall(t)
	wantResp := rest.HallResponseFromModel(hall)
	get := func(_ context.Context, _ uuid.UUID) (*catalog.Hall, error) {
		calls++
		return hall, nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	r.SetPathValue(rest.PathHallID, hall.ID.String())

	rest.GetHall(get).ServeHTTP(w, r)

	require.Equal(t, 1, calls)
	require.Equal(t, http.StatusOK, w.Code)
	var resp rest.HallResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, wantResp, resp)
}

func TestGetHall_InvalidPathHallID(t *testing.T) {
	get := func(_ context.Context, _ uuid.UUID) (*catalog.Hall, error) {
		t.Fatal()
		return nil, nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	r.SetPathValue(rest.PathHallID, "11-22-33")

	rest.GetHall(get).ServeHTTP(w, r)

	require.Equal(t, http.StatusBadRequest, w.Code)
	resp := httpxtest.DecodeErrorResponse(t, w)
	require.Equal(t, httpx.CodeBadRequest, resp.Error.Code)
	require.Equal(t, "invalid path hall id", resp.Error.Message)
}

func TestGetHall_HallNotFound(t *testing.T) {
	get := func(_ context.Context, _ uuid.UUID) (*catalog.Hall, error) {
		return nil, catalog.ErrHallNotFound
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	r.SetPathValue(rest.PathHallID, identity.NewUUID().String())

	rest.GetHall(get).ServeHTTP(w, r)

	require.Equal(t, http.StatusNotFound, w.Code)
	resp := httpxtest.DecodeErrorResponse(t, w)
	require.Equal(t, rest.CodeHallNotFound, resp.Error.Code)
}
