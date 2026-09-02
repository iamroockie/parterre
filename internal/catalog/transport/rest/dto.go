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

type HallResponse struct {
	ID         uuid.UUID             `json:"id"`
	VenueID    uuid.UUID             `json:"venue_id"`
	Name       string                `json:"name"`
	Sections   []HallSectionResponse `json:"sections"`
	TotalSeats int                   `json:"total_seats"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
}

type HallSectionResponse struct {
	Name        string `json:"name"`
	Rows        int    `json:"rows"`
	SeatsPerRow int    `json:"seats_per_row"`
	TotalSeats  int    `json:"total_seats"`
}

func HallResponseFromModel(h *catalog.Hall) HallResponse {
	sections := make([]HallSectionResponse, 0, len(h.Sections))
	for _, s := range h.Sections {
		sections = append(sections, HallSectionResponse{
			Name:        s.Name,
			Rows:        s.Rows,
			SeatsPerRow: s.SeatsPerRow,
			TotalSeats:  s.TotalSeats(),
		})
	}

	return HallResponse{
		ID:         h.ID,
		VenueID:    h.VenueID,
		Name:       h.Name,
		Sections:   sections,
		TotalSeats: h.TotalSeats(),
		CreatedAt:  h.CreatedAt,
		UpdatedAt:  h.UpdatedAt,
	}
}
