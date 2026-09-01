package rest

import (
	"time"
	"uuid"

	"github.com/iamroockie/parterre/internal/catalog"
)

type VenueResponse struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	City      string    `json:"city"`
	Address   string    `json:"address"`
	Timezone  string    `json:"timezone"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func VenueResponseFromModel(v *catalog.Venue) VenueResponse {
	return VenueResponse{
		ID:        v.ID,
		Name:      v.Name,
		City:      v.City,
		Address:   v.Address,
		Timezone:  v.Timezone.String(),
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}
