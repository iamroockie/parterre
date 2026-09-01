package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/catalog"
	"github.com/iamroockie/parterre/internal/catalog/catalogtest"
	"github.com/iamroockie/parterre/internal/catalog/usecase"
)

type saveVenue func(context.Context, *catalog.Venue) error

func (f saveVenue) Save(ctx context.Context, venue *catalog.Venue) error {
	return f(ctx, venue)
}

func TestCreateVenue(t *testing.T) {
	var calls int
	save := saveVenue(func(_ context.Context, _ *catalog.Venue) error {
		calls++
		return nil
	})
	uc := usecase.NewCreateVenue(save)
	p := catalogtest.VenueCreateParams(t)

	got, err := uc.Execute(t.Context(), p)

	require.Equal(t, 1, calls)
	require.NoError(t, err)
	require.NotNil(t, got)
}

func TestCreateVenue_InvalidVenueCreateParams(t *testing.T) {
	save := saveVenue(func(_ context.Context, _ *catalog.Venue) error {
		t.Fatal()
		return nil
	})
	uc := usecase.NewCreateVenue(save)
	// nolint:exhaustruct_v5
	p := catalog.VenueCreateParams{}

	got, err := uc.Execute(t.Context(), p)

	require.Error(t, err)
	require.Nil(t, got)
}

func TestCreateVenue_InternalError(t *testing.T) {
	var calls int
	throwErr := errors.New("error")
	save := saveVenue(func(_ context.Context, _ *catalog.Venue) error {
		calls++
		return throwErr
	})
	uc := usecase.NewCreateVenue(save)
	p := catalogtest.VenueCreateParams(t)

	got, err := uc.Execute(t.Context(), p)

	require.Equal(t, 1, calls)
	require.ErrorIs(t, err, throwErr)
	require.Nil(t, got)
}
