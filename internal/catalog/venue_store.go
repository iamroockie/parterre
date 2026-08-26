package catalog

import (
	"context"
	"errors"
	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type VenueStore struct {
	pool *pgxpool.Pool
}

func NewVenueStore(pool *pgxpool.Pool) *VenueStore {
	return &VenueStore{
		pool: pool,
	}
}

func (s *VenueStore) Create(ctx context.Context, venue *Venue) error {
	sql := `
		INSERT INTO venues (id, name, city, address, timezone, created_at, updated_at)
		VALUES (@id, @name, @city, @address, @timezone, @created_at, @updated_at)
	`

	record := venueRecordFromModel(venue)
	_, err := s.pool.Exec(ctx, sql, pgx.StrictStructArgs(record))
	if err != nil {
		return err
	}

	return nil
}

func (s *VenueStore) Get(ctx context.Context, id uuid.UUID) (*Venue, error) {
	sql := `
		SELECT id, name, city, address, timezone, created_at, updated_at
		FROM venues
		WHERE id = $1
	`
	rows, err := s.pool.Query(ctx, sql, id)
	if err != nil {
		return nil, err
	}

	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[venueRecord])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ErrVenueNotFound
		}
		return nil, err
	}

	return record.toModel()
}
