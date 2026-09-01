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

func (v *venueRecord) toModel() *catalog.Venue {
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
