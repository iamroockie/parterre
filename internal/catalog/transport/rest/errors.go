package rest

import (
	"errors"
	"net/http"

	"github.com/iamroockie/parterre/internal/catalog"
	"github.com/iamroockie/parterre/internal/platform/httpx"
)

const CodeVenueNotFound httpx.Code = "VENUE_NOT_FOUND"

var errVenueNotFound = httpx.NewError(
	http.StatusNotFound,
	CodeVenueNotFound,
	catalog.ErrVenueNotFound.Error(),
	catalog.ErrVenueNotFound,
)

func handle(fn httpx.Handler) http.HandlerFunc {
	return httpx.Handle(func(w http.ResponseWriter, r *http.Request) error {
		return toHTTP(fn(w, r))
	})
}

func toHTTP(err error) error {
	switch {
	case err == nil:
		return nil
	case errors.Is(err, catalog.ErrVenueNotFound):
		return errVenueNotFound
	default:
		return err
	}
}
