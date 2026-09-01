package rest

import (
	"net/http"

	"github.com/iamroockie/parterre/internal/catalog"
	"github.com/iamroockie/parterre/internal/platform/httpx"
	"github.com/iamroockie/parterre/internal/platform/identity"
)

func GetVenue(get catalog.GetVenueFunc) http.HandlerFunc {
	return handle(func(w http.ResponseWriter, r *http.Request) error {
		pathVenueID := r.PathValue(PathVenueID)

		venueID, err := identity.ParseUUID(pathVenueID)
		if err != nil {
			return httpx.BadRequestError("invalid path venue id", err)
		}

		venue, err := get(r.Context(), venueID)
		if err != nil {
			return err
		}

		resp := VenueResponseFromModel(venue)
		httpx.WriteJSON(w, http.StatusOK, resp)

		return nil
	})
}
