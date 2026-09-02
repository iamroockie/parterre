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

type getHall func(ctx context.Context, id uuid.UUID) (*catalog.Hall, error)

func (g getHall) Get(ctx context.Context, id uuid.UUID) (*catalog.Hall, error) {
	return g(ctx, id)
}

func TestGetHall(t *testing.T) {
	calls := 0
	hall := catalogtest.Hall(t)
	get := getHall(func(_ context.Context, _ uuid.UUID) (*catalog.Hall, error) {
		calls++
		return hall, nil
	})
	uc := usecase.NewGetHall(get)

	got, err := uc.Execute(t.Context(), hall.ID)

	require.Equal(t, 1, calls)
	require.NoError(t, err)
	require.Same(t, hall, got)
}
