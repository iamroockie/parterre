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

func TestGetVenue(t *testing.T) {
	var calls int
	venue := catalogtest.Venue(t)
	wantResp := rest.VenueResponseFromModel(venue)
	getVenue := func(_ context.Context, _ uuid.UUID) (*catalog.Venue, error) {
		calls++
		return venue, nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	r.SetPathValue(rest.PathVenueID, venue.ID.String())

	rest.GetVenue(getVenue).ServeHTTP(w, r)

	require.Equal(t, 1, calls)
	require.Equal(t, http.StatusOK, w.Code)
	var resp rest.VenueResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, wantResp, resp)
}

func TestGetVenue_InvalidPathVenueID(t *testing.T) {
	getVenue := func(_ context.Context, _ uuid.UUID) (*catalog.Venue, error) {
		t.Fatal()
		return nil, nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	r.SetPathValue(rest.PathVenueID, "11-22-33")

	rest.GetVenue(getVenue).ServeHTTP(w, r)

	require.Equal(t, http.StatusBadRequest, w.Code)
	resp := httpxtest.DecodeErrorResponse(t, w)
	require.Equal(t, httpx.CodeBadRequest, resp.Error.Code)
	require.Equal(t, "invalid path venue id", resp.Error.Message)
}

func TestGetVenue_VenueNotFound(t *testing.T) {
	getVenue := func(_ context.Context, _ uuid.UUID) (*catalog.Venue, error) {
		return nil, catalog.ErrVenueNotFound
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(t.Context(), http.MethodGet, "/", nil)
	r.SetPathValue(rest.PathVenueID, identity.NewUUID().String())

	rest.GetVenue(getVenue).ServeHTTP(w, r)

	require.Equal(t, http.StatusNotFound, w.Code)
	resp := httpxtest.DecodeErrorResponse(t, w)
	require.Equal(t, rest.CodeVenueNotFound, resp.Error.Code)
}
