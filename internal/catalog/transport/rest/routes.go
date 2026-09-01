package rest

import (
	"fmt"

	"github.com/go-chi/chi/v5"

	"github.com/iamroockie/parterre/internal/catalog"
)

const PathVenueID = "venue_id"

func RegisterRoutes(r chi.Router, m catalog.Module) {
	r.Route("/venues", func(r chi.Router) {
		r.Post("/", CreateVenue(m.CreateVenue))
		r.Get(fmt.Sprintf("/{%s}", PathVenueID), GetVenue(m.GetVenue))
	})
}
