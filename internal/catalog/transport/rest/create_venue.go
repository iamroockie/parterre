package rest

import (
	"encoding/json/v2"
	"net/http"
	"path"

	"github.com/iamroockie/parterre/internal/catalog"
	"github.com/iamroockie/parterre/internal/platform/httpx"
)

type CreateVenueRequest struct {
	Name     string `json:"name"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Timezone string `json:"timezone"`
}

func CreateVenue(create catalog.CreateVenueFunc) http.HandlerFunc {
	return handle(func(w http.ResponseWriter, r *http.Request) error {
		var req CreateVenueRequest
		err := json.UnmarshalRead(r.Body, &req, json.RejectUnknownMembers(true))
		if err != nil {
			return httpx.BadRequestError("invalid request body", err)
		}

		venue, err := create(r.Context(), catalog.VenueCreateParams(req))
		if err != nil {
			return err
		}

		resp := VenueResponseFromModel(venue)
		w.Header().Set("Location", path.Join(r.URL.Path, resp.ID.String()))
		httpx.WriteJSON(w, http.StatusCreated, resp)

		return nil
	})
}
