package rest

import (
	"encoding/json/v2"
	"net/http"
	"path"
	"uuid"

	"github.com/iamroockie/parterre/internal/catalog"
	"github.com/iamroockie/parterre/internal/platform/httpx"
)

type CreateHallRequest struct {
	VenueID  uuid.UUID                  `json:"venue_id"`
	Name     string                     `json:"name"`
	Sections []CreateHallSectionRequest `json:"sections"`
}

type CreateHallSectionRequest struct {
	Name        string `json:"name"`
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seats_per_row"`
}

func (r CreateHallRequest) ToCreateParams() catalog.HallCreateParams {
	sections := make([]catalog.Section, 0, len(r.Sections))
	for _, s := range r.Sections {
		sections = append(sections, catalog.Section(s))
	}

	return catalog.HallCreateParams{
		VenueID:  r.VenueID,
		Name:     r.Name,
		Sections: sections,
	}
}

func CreateHall(create catalog.CreateHallFunc) http.HandlerFunc {
	return handle(func(w http.ResponseWriter, r *http.Request) error {
		var req CreateHallRequest
		err := json.UnmarshalRead(r.Body, &req, json.RejectUnknownMembers(true))
		if err != nil {
			return httpx.BadRequestError("invalid request body", err)
		}

		got, err := create(r.Context(), req.ToCreateParams())
		if err != nil {
			return err
		}

		resp := HallResponseFromModel(got)
		w.Header().Set("Location", path.Join(r.URL.Path, resp.ID.String()))
		httpx.WriteJSON(w, http.StatusCreated, resp)

		return nil
	})
}
