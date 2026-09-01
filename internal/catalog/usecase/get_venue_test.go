package usecase_test

import (
	"context"
	"testing"
	"uuid"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/catalog"
	"github.com/iamroockie/parterre/internal/catalog/catalogtest"
	"github.com/iamroockie/parterre/internal/catalog/usecase"
)

type getVenue func(context.Context, uuid.UUID) (*catalog.Venue, error)

func (f getVenue) Get(ctx context.Context, id uuid.UUID) (*catalog.Venue, error) {
	return f(ctx, id)
}

func TestGetVenue(t *testing.T) {
	var calls int
	venue := catalogtest.Venue(t)
	get := getVenue(func(_ context.Context, _ uuid.UUID) (*catalog.Venue, error) {
		calls++
		return venue, nil
	})
	uc := usecase.NewGetVenue(get)

	got, err := uc.Execute(t.Context(), venue.ID)

	require.Equal(t, 1, calls)
	require.NoError(t, err)
	require.Same(t, venue, got)
}
