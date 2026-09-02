package postgres

import (
	"context"
	"errors"
	"uuid"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/iamroockie/parterre/internal/catalog"
	"github.com/iamroockie/parterre/internal/platform/pg"
)

type HallStore struct {
	pool *pgxpool.Pool
}

func NewHallStore(pool *pgxpool.Pool) *HallStore {
	return &HallStore{pool}
}

func (s *HallStore) Save(ctx context.Context, hall *catalog.Hall) error {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	record := hallRecordFromModel(hall)

	sql := `
		INSERT INTO halls (id, venue_id, name, created_at, updated_at)
		VALUES (@id, @venue_id, @name, @created_at, @updated_at)
	`
	_, err = tx.Exec(ctx, sql, pgx.StrictStructArgs(record))
	if err != nil {
		return hallSaveError(err)
	}

	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"hall_sections"},
		[]string{"hall_id", "name", "rows_count", "seats_per_row"},
		pgx.CopyFromSlice(len(hall.Sections), func(i int) ([]any, error) {
			section := hall.Sections[i]
			return []any{hall.ID, section.Name, section.Rows, section.SeatsPerRow}, nil
		}),
	)
	if err != nil {
		return hallSaveError(err)
	}

	seats := hall.GenerateSeats()
	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"seats"},
		[]string{"id", "hall_id", "section", "row_number", "number"},
		pgx.CopyFromSlice(len(seats), func(i int) ([]any, error) {
			seat := seats[i]
			return []any{seat.ID, hall.ID, seat.Section, seat.RowNumber, seat.Number}, nil
		}),
	)
	if err != nil {
		return hallSaveError(err)
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func hallSaveError(err error) error {
	e, ok := errors.AsType[*pgconn.PgError](err)
	if !ok {
		return err
	}

	switch {
	case e.Code == pg.CodeForeignKeyViolation:
		return catalog.ErrVenueNotFound
	case e.Code == pg.CodeUniqueViolation && e.ConstraintName == "halls_venue_id_name_unique":
		return catalog.ErrHallNameTaken
	case e.Code == pg.CodeUniqueViolation && e.ConstraintName == "seats_position_unique":
		return catalog.ErrHallSeatTaken
	default:
		return err
	}
}

func (s *HallStore) Get(ctx context.Context, id uuid.UUID) (*catalog.Hall, error) {
	sql := `
		SELECT id, venue_id, name, created_at, updated_at
		FROM halls
		WHERE id = $1
	`
	rows, err := s.pool.Query(ctx, sql, id)
	if err != nil {
		return nil, err
	}
	record, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[hallRecord])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, catalog.ErrHallNotFound
		}
		return nil, err
	}

	sql = `
		SELECT name, rows_count, seats_per_row
		FROM hall_sections
		WHERE hall_id = $1
		ORDER BY name
	`
	rows, err = s.pool.Query(ctx, sql, id)
	if err != nil {
		return nil, err
	}
	sections, err := pgx.CollectRows(rows, pgx.RowToStructByName[hallSectionRecord])
	if err != nil {
		return nil, err
	}

	return record.toModel(sections), nil
}
