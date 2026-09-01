package catalog

import (
	"fmt"
	"strings"
	"time"
	"uuid"

	v "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/iamroockie/parterre/internal/platform/identity"
)

type Hall struct {
	ID        uuid.UUID
	VenueID   uuid.UUID
	Name      string
	Seats     []Seat
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Seat struct {
	ID        uuid.UUID
	Section   string
	RowNumber int
	Number    int
}

type HallCreateParams struct {
	VenueID  string
	Name     string
	Sections []SectionParams
}

type SectionParams struct {
	Name        string
	Rows        int
	SeatsPerRow int
}

func NewHall(p HallCreateParams) (*Hall, error) {
	venueID, vidErr := identity.ParseUUID(p.VenueID)
	hallName := strings.TrimSpace(p.Name)
	totalSeats := 0
	for _, section := range p.Sections {
		totalSeats += section.Rows * section.SeatsPerRow
	}

	err := v.Errors{
		"venue_id":    vidErr,
		"name":        v.Validate(hallName, v.Required, v.RuneLength(2, 100)),
		"total_seats": v.Validate(totalSeats, v.Max(250_000)),
		"sections": v.Validate(p.Sections, v.Required, v.Length(1, 50), v.By(func(_ any) error {
			sectionNames := make(map[string]int, len(p.Sections))
			for i, section := range p.Sections {
				name := strings.TrimSpace(section.Name)
				idx, ok := sectionNames[name]
				if ok {
					return fmt.Errorf("section %d and section %d have same names", idx, i)
				}
				sectionNames[name] = i
			}

			return nil
		})),
	}.Filter()
	if err != nil {
		return nil, err
	}

	var seats []Seat
	for _, section := range p.Sections {
		sectionName := strings.TrimSpace(section.Name)
		for row := 1; row <= section.Rows; row++ {
			for seat := 1; seat <= section.SeatsPerRow; seat++ {
				seats = append(seats, Seat{
					ID:        identity.NewUUID(),
					Section:   sectionName,
					RowNumber: row,
					Number:    seat,
				})
			}
		}
	}

	now := time.Now().UTC()

	return &Hall{
		ID:        identity.NewUUID(),
		VenueID:   venueID,
		Name:      hallName,
		Seats:     seats,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (p SectionParams) Validate() error {
	name := strings.TrimSpace(p.Name)

	return v.Errors{
		"name":          v.Validate(name, v.Required, v.RuneLength(2, 100)),
		"rows":          v.Validate(p.Rows, v.Min(1), v.Max(100)),
		"seats_per_row": v.Validate(p.SeatsPerRow, v.Min(1), v.Max(200)),
	}.Filter()
}
