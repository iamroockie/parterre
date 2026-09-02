package postgres_test

import (
	"testing"
	"time"
	"uuid"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/catalog"
	"github.com/iamroockie/parterre/internal/catalog/catalogtest"
	"github.com/iamroockie/parterre/internal/catalog/infra/postgres"
	"github.com/iamroockie/parterre/internal/platform/identity"
	"github.com/iamroockie/parterre/internal/platform/pg/pgtest"
)

func saveVenue(t *testing.T, pool *pgxpool.Pool) uuid.UUID {
	t.Helper()

	store := postgres.NewVenueStore(pool)
	venue := catalogtest.Venue(t)
	err := store.Save(t.Context(), venue)
	require.NoError(t, err)

	return venue.ID
}

func TestHallStore_CreateAndGet(t *testing.T) {
	pool := pgtest.NewTestPool(t)
	store := postgres.NewHallStore(pool)
	hall := catalogtest.Hall(t)
	hall.VenueID = saveVenue(t, pool)

	_, err := store.Get(t.Context(), hall.ID)
	require.ErrorIs(t, err, catalog.ErrHallNotFound)

	err = store.Save(t.Context(), hall)
	require.NoError(t, err)

	got, err := store.Get(t.Context(), hall.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, hall.ID, got.ID)
	require.Equal(t, hall.VenueID, got.VenueID)
	require.Equal(t, hall.Name, got.Name)
	require.Equal(t, hall.Sections, got.Sections, "sections must survive the round trip")
	require.Equal(t, hall.TotalSeats(), got.TotalSeats())
	require.WithinDuration(t, hall.CreatedAt, got.CreatedAt, time.Microsecond)
	require.WithinDuration(t, hall.UpdatedAt, got.UpdatedAt, time.Microsecond)
}

func TestHallStore_Save_WritesSeats(t *testing.T) {
	pool := pgtest.NewTestPool(t)
	store := postgres.NewHallStore(pool)
	hall := catalogtest.Hall(t)
	hall.VenueID = saveVenue(t, pool)

	err := store.Save(t.Context(), hall)
	require.NoError(t, err)

	var seats int
	err = pool.QueryRow(t.Context(),
		`SELECT COUNT(*) FROM seats WHERE hall_id = $1`, hall.ID).Scan(&seats)
	require.NoError(t, err)
	require.Equal(t, hall.TotalSeats(), seats)

	var vipSeats int
	sql := `SELECT COUNT(*) FROM seats WHERE hall_id = $1 AND section = 'VIP'`
	err = pool.QueryRow(t.Context(), sql, hall.ID).Scan(&vipSeats)
	require.NoError(t, err)
	require.Equal(t, 50, vipSeats)
}

func TestHallStore_Save_VenueNotFound(t *testing.T) {
	pool := pgtest.NewTestPool(t)
	store := postgres.NewHallStore(pool)
	hall := catalogtest.Hall(t)

	err := store.Save(t.Context(), hall)

	require.ErrorIs(t, err, catalog.ErrVenueNotFound)
}

func TestHallStore_Save_HallNameTaken(t *testing.T) {
	pool := pgtest.NewTestPool(t)
	store := postgres.NewHallStore(pool)
	venueID := saveVenue(t, pool)

	first := catalogtest.Hall(t)
	first.VenueID = venueID
	require.NoError(t, store.Save(t.Context(), first))

	second := catalogtest.Hall(t)
	second.VenueID = venueID
	second.Name = first.Name

	err := store.Save(t.Context(), second)

	require.ErrorIs(t, err, catalog.ErrHallNameTaken)
}

func TestHallStore_Save_RollsBackOnFailure(t *testing.T) {
	pool := pgtest.NewTestPool(t)
	store := postgres.NewHallStore(pool)
	hall := catalogtest.Hall(t)

	err := store.Save(t.Context(), hall)
	require.ErrorIs(t, err, catalog.ErrVenueNotFound)

	var seats int
	err = pool.QueryRow(t.Context(),
		`SELECT COUNT(*) FROM seats WHERE hall_id = $1`, hall.ID).Scan(&seats)
	require.NoError(t, err)
	require.Zero(t, seats)
}

func TestHallStore_Get_NotFound(t *testing.T) {
	pool := pgtest.NewTestPool(t)
	store := postgres.NewHallStore(pool)

	_, err := store.Get(t.Context(), identity.NewUUID())

	require.ErrorIs(t, err, catalog.ErrHallNotFound)
}
