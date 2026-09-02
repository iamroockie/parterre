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

type saveHall func(ctx context.Context, hall *catalog.Hall) error

func (f saveHall) Save(ctx context.Context, hall *catalog.Hall) error {
	return f(ctx, hall)
}

func TestCreateHall(t *testing.T) {
	calls := 0
	save := saveHall(func(_ context.Context, _ *catalog.Hall) error {
		calls++
		return nil
	})
	uc := usecase.NewCreateHall(save)
	p := catalogtest.HallCreateParams(t)

	got, err := uc.Execute(t.Context(), p)

	require.Equal(t, 1, calls)
	require.NoError(t, err)
	require.NotNil(t, got)
}

func TestCreateHall_InvalidCreateParams(t *testing.T) {
	save := saveHall(func(_ context.Context, _ *catalog.Hall) error {
		t.Fatal()
		return nil
	})
	uc := usecase.NewCreateHall(save)
	// nolint:exhaustruct_v5
	p := catalog.HallCreateParams{}

	got, err := uc.Execute(t.Context(), p)

	require.Error(t, err)
	require.Nil(t, got)
}

func TestCreateHall_InternalError(t *testing.T) {
	calls := 0
	throwErr := errors.New("error")
	save := saveHall(func(_ context.Context, _ *catalog.Hall) error {
		calls++
		return throwErr
	})
	uc := usecase.NewCreateHall(save)
	p := catalogtest.HallCreateParams(t)

	got, err := uc.Execute(t.Context(), p)

	require.Equal(t, 1, calls)
	require.ErrorIs(t, err, throwErr)
	require.Nil(t, got)
}
