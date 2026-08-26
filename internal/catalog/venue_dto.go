package catalog

import (
	"time"
	"uuid"
)

type createVenueRequest struct {
	Name     string `json:"name"`
	City     string `json:"city"`
	Address  string `json:"address"`
	Timezone string `json:"timezone"`
}

type venueResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	City      string    `json:"city"`
	Address   string    `json:"address"`
	Timezone  string    `json:"timezone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func venueResponseFromModel(venue *Venue) *venueResponse {
	return &venueResponse{
		ID:        venue.ID,
		Name:      venue.Name,
		City:      venue.City,
		Address:   venue.Address,
		Timezone:  venue.Timezone.String(),
		CreatedAt: venue.CreatedAt,
		UpdatedAt: venue.UpdatedAt,
	}
}
