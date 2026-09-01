package postgres_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/catalog"
	"github.com/iamroockie/parterre/internal/catalog/catalogtest"
	"github.com/iamroockie/parterre/internal/catalog/infra/postgres"
	"github.com/iamroockie/parterre/internal/platform/pg/pgtest"
)

func TestVenueStore_CreateAndGet(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	pool := pgtest.NewTestPool(t)
	store := postgres.NewVenueStore(pool)
	venue := catalogtest.Venue(t)

	_, err := store.Get(t.Context(), venue.ID)
	require.ErrorIs(t, err, catalog.ErrVenueNotFound)

	err = store.Save(t.Context(), venue)
	require.NoError(t, err)

	got, err := store.Get(t.Context(), venue.ID)
	require.NoError(t, err)
	require.NotNil(t, got)
	require.Equal(t, venue.ID, got.ID)
	require.Equal(t, venue.Name, got.Name)
	require.Equal(t, venue.City, got.City)
	require.Equal(t, venue.Address, got.Address)
	require.Equal(t, venue.Timezone, got.Timezone)
	require.WithinDuration(t, venue.CreatedAt, got.CreatedAt, time.Microsecond)
	require.WithinDuration(t, venue.UpdatedAt, got.UpdatedAt, time.Microsecond)
}
