package catalog_test

import (
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/catalog"
	"github.com/iamroockie/parterre/internal/catalog/catalogtest"
	"github.com/iamroockie/parterre/internal/platform/postgres/postgrestest"
)

func TestVenueStore_CreateAndGet(t *testing.T) {
	pool := postgrestest.NewTestPool(t)
	store := catalog.NewVenueStore(pool)
	venue := catalogtest.NewVenue(t)

	err := store.Create(t.Context(), venue)
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

func TestVenueStore_NotFound(t *testing.T) {
	pool := postgrestest.NewTestPool(t)
	store := catalog.NewVenueStore(pool)

	got, err := store.Get(t.Context(), uuid.NewV7())

	require.ErrorIs(t, err, catalog.ErrVenueNotFound)
	require.Nil(t, got)
}
