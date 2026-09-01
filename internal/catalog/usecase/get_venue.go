package usecase

import (
	"context"
	"uuid"

	"github.com/iamroockie/parterre/internal/catalog"
)

type GetVenue struct {
	getter venueGetter
}

func NewGetVenue(getter venueGetter) GetVenue {
	return GetVenue{getter}
}

func (g GetVenue) Execute(ctx context.Context, id uuid.UUID) (*catalog.Venue, error) {
	return g.getter.Get(ctx, id)
}
