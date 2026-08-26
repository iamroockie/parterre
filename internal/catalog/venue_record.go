package catalog

import (
	"fmt"
	"time"
	"uuid"
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

func venueRecordFromModel(v *Venue) *venueRecord {
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

func (v *venueRecord) toModel() (*Venue, error) {
	tz, err := time.LoadLocation(v.Timezone)
	if err != nil {
		return nil, fmt.Errorf("invalid timezone: %w", err)
	}

	return &Venue{
		ID:        v.ID,
		Name:      v.Name,
		City:      v.City,
		Address:   v.Address,
		Timezone:  tz,
		CreatedAt: v.CreatedAt,
		UpdatedAt: v.UpdatedAt,
	}, nil
}
