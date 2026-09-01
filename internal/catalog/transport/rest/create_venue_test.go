package rest_test

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"errors"
	"net/http"
	"net/http/httptest"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/catalog"
	"github.com/iamroockie/parterre/internal/catalog/catalogtest"
	"github.com/iamroockie/parterre/internal/catalog/transport/rest"
	"github.com/iamroockie/parterre/internal/platform/httpx"
	"github.com/iamroockie/parterre/internal/platform/httpx/httpxtest"
)

func venueCreateRequestData(t testing.TB) []byte {
	t.Helper()

	req := rest.CreateVenueRequest(catalogtest.VenueCreateParams(t))
	data, err := json.Marshal(req)
	require.NoError(t, err)

	return data
}

func TestCreateVenue(t *testing.T) {
	var calls int
	venue := catalogtest.Venue(t)
	wantResp := rest.VenueResponseFromModel(venue)
	req := venueCreateRequestData(t)
	createVenue := func(_ context.Context, _ catalog.VenueCreateParams) (*catalog.Venue, error) {
		calls++
		return venue, nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", bytes.NewReader(req))

	rest.CreateVenue(createVenue).ServeHTTP(w, r)

	require.Equal(t, 1, calls)
	require.Equal(t, http.StatusCreated, w.Code)
	var resp rest.VenueResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, wantResp, resp)
	location := path.Join(r.URL.Path, venue.ID.String())
	require.Equal(t, location, w.Result().Header.Get("Location"))
}

func TestCreateVenue_InvalidBody(t *testing.T) {
	createVenue := func(_ context.Context, _ catalog.VenueCreateParams) (*catalog.Venue, error) {
		t.Fatal()
		return nil, nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", strings.NewReader(""))

	rest.CreateVenue(createVenue).ServeHTTP(w, r)

	require.Equal(t, http.StatusBadRequest, w.Code)
	resp := httpxtest.DecodeErrorResponse(t, w)
	require.Equal(t, httpx.CodeBadRequest, resp.Error.Code)
}

func TestCreateVenue_InternalError(t *testing.T) {
	var calls int
	err := errors.New("internal test error")
	createVenue := func(_ context.Context, _ catalog.VenueCreateParams) (*catalog.Venue, error) {
		calls++
		return nil, err
	}
	w := httptest.NewRecorder()
	req := venueCreateRequestData(t)
	r := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", bytes.NewReader(req))

	rest.CreateVenue(createVenue).ServeHTTP(w, r)

	require.Equal(t, 1, calls)
	require.Equal(t, http.StatusInternalServerError, w.Code)
	resp := httpxtest.DecodeErrorResponse(t, w)
	require.Equal(t, httpx.CodeInternalError, resp.Error.Code)
}
