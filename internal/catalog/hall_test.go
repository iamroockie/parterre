package catalog_test

import (
	"strings"
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/catalog"
	"github.com/iamroockie/parterre/internal/platform/identity"
)

func TestHall_New(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		params := catalog.HallCreateParams{
			VenueID: identity.NewUUID(),
			Name:    "Test name",
			Sections: []catalog.Section{
				{Name: "Default", Rows: 20, SeatsPerRow: 50},
				{Name: "VIP", Rows: 5, SeatsPerRow: 10},
			},
		}

		got, err := catalog.NewHall(params)

		require.NoError(t, err)
		require.NotNil(t, got)
		require.NotEqual(t, uuid.Nil(), got.ID)
		require.Equal(t, params.VenueID, got.VenueID)
		require.Equal(t, params.Name, got.Name)
		require.Equal(t, params.Sections, got.Sections)
		require.Equal(t, 1050, got.TotalSeats())
		require.Equal(t, got.CreatedAt, got.UpdatedAt)
		require.WithinDuration(t, got.CreatedAt.UTC(), got.CreatedAt, time.Microsecond)
		require.WithinDuration(t, got.UpdatedAt.UTC(), got.UpdatedAt, time.Microsecond)
	})

	t.Run("trims name and section names", func(t *testing.T) {
		params := catalog.HallCreateParams{
			VenueID:  identity.NewUUID(),
			Name:     "  Test name  ",
			Sections: []catalog.Section{{Name: "  VIP  ", Rows: 1, SeatsPerRow: 1}},
		}

		got, err := catalog.NewHall(params)

		require.NoError(t, err)
		require.Equal(t, "Test name", got.Name)
		require.Equal(t, "VIP", got.Sections[0].Name)
	})

	t.Run("orders sections by name", func(t *testing.T) {
		params := catalog.HallCreateParams{
			VenueID: identity.NewUUID(),
			Name:    "Test name",
			Sections: []catalog.Section{
				{Name: "Партер", Rows: 10, SeatsPerRow: 20},
				{Name: "Амфитеатр", Rows: 5, SeatsPerRow: 20},
				{Name: "Балкон", Rows: 3, SeatsPerRow: 20},
			},
		}

		got, err := catalog.NewHall(params)

		require.NoError(t, err)
		require.Equal(t, []string{"Амфитеатр", "Балкон", "Партер"}, sectionNames(got.Sections))
	})
}

func TestHall_New_InvalidParams(t *testing.T) {
	validSections := []catalog.Section{{Name: "Default", Rows: 20, SeatsPerRow: 50}}

	tests := map[string]struct {
		params    catalog.HallCreateParams
		wantField string
	}{
		"empty name": {
			params:    hallParams("   ", validSections),
			wantField: "name",
		},
		"too short name": {
			params:    hallParams("A", validSections),
			wantField: "name",
		},
		"too long name": {
			params:    hallParams(strings.Repeat("a", 101), validSections),
			wantField: "name",
		},
		"no sections": {
			params:    hallParams("Test name", nil),
			wantField: "sections",
		},
		"too many sections": {
			params:    hallParams("Test name", repeatSections(51, 1, 1)),
			wantField: "sections",
		},
		"duplicate section names": {
			params: hallParams("Test name", []catalog.Section{
				{Name: "VIP", Rows: 1, SeatsPerRow: 1},
				{Name: "  VIP  ", Rows: 2, SeatsPerRow: 2},
			}),
			wantField: "sections",
		},
		"section without rows": {
			params: hallParams("Test name", []catalog.Section{
				{Name: "VIP", Rows: 0, SeatsPerRow: 10},
			}),
			wantField: "sections",
		},
		"section with too many seats per row": {
			params: hallParams("Test name", []catalog.Section{
				{Name: "VIP", Rows: 1, SeatsPerRow: 201},
			}),
			wantField: "sections",
		},
		"too many seats in total": {
			params:    hallParams("Test name", repeatSections(13, 100, 200)),
			wantField: "total_seats",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := catalog.NewHall(tt.params)

			require.Nil(t, got)
			require.ErrorContains(t, err, tt.wantField)
		})
	}
}

func TestHall_GenerateSeats(t *testing.T) {
	hall, err := catalog.NewHall(catalog.HallCreateParams{
		VenueID: identity.NewUUID(),
		Name:    "Test name",
		Sections: []catalog.Section{
			{Name: "Партер", Rows: 2, SeatsPerRow: 3},
			{Name: "Балкон", Rows: 1, SeatsPerRow: 2},
		},
	})
	require.NoError(t, err)

	seats := hall.GenerateSeats()

	require.Len(t, seats, hall.TotalSeats())
	require.Equal(t, catalog.Seat{
		ID:        seats[0].ID,
		Section:   "Балкон",
		RowNumber: 1,
		Number:    1,
	}, seats[0])
	require.Equal(t, "Партер", seats[len(seats)-1].Section)

	type position struct {
		section   string
		rowNumber int
		number    int
	}

	positions := make(map[position]struct{}, len(seats))
	ids := make(map[uuid.UUID]struct{}, len(seats))
	for _, seat := range seats {
		require.NotEqual(t, uuid.Nil(), seat.ID)
		ids[seat.ID] = struct{}{}
		positions[position{seat.Section, seat.RowNumber, seat.Number}] = struct{}{}
	}
	require.Len(t, ids, len(seats), "seat ids must be unique")
	require.Len(t, positions, len(seats), "seat positions must be unique")
}

func sectionNames(sections []catalog.Section) []string {
	names := make([]string, 0, len(sections))
	for _, s := range sections {
		names = append(names, s.Name)
	}

	return names
}

func hallParams(name string, sections []catalog.Section) catalog.HallCreateParams {
	return catalog.HallCreateParams{
		VenueID:  identity.NewUUID(),
		Name:     name,
		Sections: sections,
	}
}

func repeatSections(n, rows, seatsPerRow int) []catalog.Section {
	sections := make([]catalog.Section, 0, n)
	for i := range n {
		sections = append(sections, catalog.Section{
			Name:        "Section " + strings.Repeat("a", i+1),
			Rows:        rows,
			SeatsPerRow: seatsPerRow,
		})
	}

	return sections
}
