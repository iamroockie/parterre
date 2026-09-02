package catalog

import (
	"fmt"
	"slices"
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
	Sections  []Section
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Section struct {
	Name        string
	Rows        int
	SeatsPerRow int
}

type Seat struct {
	ID        uuid.UUID
	Section   string
	RowNumber int
	Number    int
}

type HallCreateParams struct {
	VenueID  uuid.UUID
	Name     string
	Sections []Section
}

func NewHall(p HallCreateParams) (*Hall, error) {
	hallName := strings.TrimSpace(p.Name)

	sections := make([]Section, 0, len(p.Sections))
	totalSeats := 0
	for _, section := range p.Sections {
		section.Name = strings.TrimSpace(section.Name)
		sections = append(sections, section)
		totalSeats += section.TotalSeats()
	}
	slices.SortFunc(sections, func(a, b Section) int {
		return strings.Compare(a.Name, b.Name)
	})

	err := v.Errors{
		"name":        v.Validate(hallName, v.Required, v.RuneLength(2, 100)),
		"total_seats": v.Validate(totalSeats, v.Max(250_000)),
		"sections": v.Validate(sections, v.Required, v.Length(1, 50), v.By(func(_ any) error {
			seen := make(map[string]struct{}, len(sections))
			for _, section := range sections {
				if _, ok := seen[section.Name]; ok {
					return fmt.Errorf("duplicate section name %q", section.Name)
				}
				seen[section.Name] = struct{}{}
			}

			return nil
		})),
	}.Filter()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	return &Hall{
		ID:        identity.NewUUID(),
		VenueID:   p.VenueID,
		Name:      hallName,
		Sections:  sections,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func (s Section) Validate() error {
	name := strings.TrimSpace(s.Name)

	return v.Errors{
		"name":          v.Validate(name, v.Required, v.RuneLength(2, 100)),
		"rows":          v.Validate(s.Rows, v.Min(1), v.Max(100)),
		"seats_per_row": v.Validate(s.SeatsPerRow, v.Min(1), v.Max(200)),
	}.Filter()
}

func (s Section) TotalSeats() int {
	return s.Rows * s.SeatsPerRow
}

func (h Hall) TotalSeats() int {
	seats := 0
	for _, s := range h.Sections {
		seats += s.TotalSeats()
	}

	return seats
}

func (h Hall) GenerateSeats() []Seat {
	seats := make([]Seat, 0, h.TotalSeats())
	for _, sec := range h.Sections {
		for row := 1; row <= sec.Rows; row++ {
			for n := 1; n <= sec.SeatsPerRow; n++ {
				seats = append(seats, Seat{
					ID:        identity.NewUUID(),
					Section:   sec.Name,
					RowNumber: row,
					Number:    n,
				})
			}
		}
	}

	return seats
}
