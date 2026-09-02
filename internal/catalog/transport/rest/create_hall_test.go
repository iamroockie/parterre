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

func hallCreateRequestData(t testing.TB) []byte {
	t.Helper()

	p := catalogtest.HallCreateParams(t)
	sections := make([]rest.CreateHallSectionRequest, 0, len(p.Sections))
	for _, s := range p.Sections {
		sections = append(sections, rest.CreateHallSectionRequest(s))
	}
	req := rest.CreateHallRequest{
		VenueID:  p.VenueID,
		Name:     p.Name,
		Sections: sections,
	}
	data, err := json.Marshal(req)
	require.NoError(t, err)

	return data
}

func TestCreateHall(t *testing.T) {
	calls := 0
	hall := catalogtest.Hall(t)
	wantResp := rest.HallResponseFromModel(hall)
	req := hallCreateRequestData(t)
	create := func(_ context.Context, _ catalog.HallCreateParams) (*catalog.Hall, error) {
		calls++
		return hall, nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", bytes.NewReader(req))

	rest.CreateHall(create).ServeHTTP(w, r)
	require.Equal(t, 1, calls)
	require.Equal(t, http.StatusCreated, w.Code)
	var resp rest.HallResponse
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &resp))
	require.Equal(t, wantResp, resp)
	location := path.Join(r.URL.Path, hall.ID.String())
	require.Equal(t, location, w.Result().Header.Get("Location"))
}

func TestCreateHall_InvalidBody(t *testing.T) {
	create := func(_ context.Context, _ catalog.HallCreateParams) (*catalog.Hall, error) {
		t.Fatal()
		return nil, nil
	}
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", strings.NewReader(""))

	rest.CreateHall(create).ServeHTTP(w, r)

	require.Equal(t, http.StatusBadRequest, w.Code)
	resp := httpxtest.DecodeErrorResponse(t, w)
	require.Equal(t, httpx.CodeBadRequest, resp.Error.Code)
}

func TestCreateHall_InternalError(t *testing.T) {
	calls := 0
	err := errors.New("internal test error")
	create := func(_ context.Context, _ catalog.HallCreateParams) (*catalog.Hall, error) {
		calls++
		return nil, err
	}
	w := httptest.NewRecorder()
	req := hallCreateRequestData(t)
	r := httptest.NewRequestWithContext(t.Context(), http.MethodPost, "/", bytes.NewReader(req))

	rest.CreateHall(create).ServeHTTP(w, r)

	require.Equal(t, 1, calls)
	require.Equal(t, http.StatusInternalServerError, w.Code)
	resp := httpxtest.DecodeErrorResponse(t, w)
	require.Equal(t, httpx.CodeInternalError, resp.Error.Code)
}
