package rest

import (
	"net/http"

	"github.com/iamroockie/parterre/internal/catalog"
	"github.com/iamroockie/parterre/internal/platform/httpx"
	"github.com/iamroockie/parterre/internal/platform/identity"
)

func GetHall(get catalog.GetHallFunc) http.HandlerFunc {
	return handle(func(w http.ResponseWriter, r *http.Request) error {
		hallHallID := r.PathValue(PathHallID)

		hallID, err := identity.ParseUUID(hallHallID)
		if err != nil {
			return httpx.BadRequestError("invalid path hall id", err)
		}

		hall, err := get(r.Context(), hallID)
		if err != nil {
			return err
		}

		resp := HallResponseFromModel(hall)
		httpx.WriteJSON(w, http.StatusOK, resp)

		return nil
	})
}
