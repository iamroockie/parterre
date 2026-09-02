package postgres

import (
	"time"
	"uuid"

	"github.com/iamroockie/parterre/internal/catalog"
)

type venueRecord struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	City      string    `db:"city"`
	Address   string    `db:"address"`
	Timezone  string    `db:"timezone"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type hallRecord struct {
	ID        uuid.UUID `db:"id"`
	VenueID   uuid.UUID `db:"venue_id"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type hallSectionRecord struct {
	Name        string `db:"name"`
	Rows        int    `db:"rows_count"`
	SeatsPerRow int    `db:"seats_per_row"`
}

func venueRecordFromModel(v *catalog.Venue) *venueRecord {
	return &venueRecord{
		ID:        v.ID,
		Name:      v.Name,
		City:      v.City,
		Address:   v.Address,
		Timezone:  v.Timezone.String(),
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}

func hallRecordFromModel(h *catalog.Hall) *hallRecord {
	return &hallRecord{
		ID:        h.ID,
		VenueID:   h.VenueID,
		Name:      h.Name,
		CreatedAt: h.CreatedAt,
		UpdatedAt: h.UpdatedAt,
	}
}

func (h hallRecord) toModel(sectionRecords []hallSectionRecord) *catalog.Hall {
	sections := make([]catalog.Section, 0, len(sectionRecords))
	for _, s := range sectionRecords {
		sections = append(sections, catalog.Section(s))
	}

	return &catalog.Hall{
		ID:        h.ID,
		VenueID:   h.VenueID,
		Name:      h.Name,
		Sections:  sections,
		CreatedAt: h.CreatedAt,
		UpdatedAt: h.UpdatedAt,
	}
}

func (v venueRecord) toModel() *catalog.Venue {
	tz, _ := time.LoadLocation(v.Timezone)

	return &catalog.Venue{
		ID:        v.ID,
		Name:      v.Name,
		City:      v.City,
		Address:   v.Address,
		Timezone:  tz,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}
}
