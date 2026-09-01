package usecase

import (
	"context"
	"fmt"

	"github.com/iamroockie/parterre/internal/catalog"
)

type CreateVenue struct {
	saver venueSaver
}

func NewCreateVenue(saver venueSaver) CreateVenue {
	return CreateVenue{saver}
}

func (uc CreateVenue) Execute(
	ctx context.Context,
	p catalog.VenueCreateParams,
) (*catalog.Venue, error) {
	venue, err := catalog.NewVenue(p)
	if err != nil {
		return nil, fmt.Errorf("create venue: %w", err)
	}

	if err := uc.saver.Save(ctx, venue); err != nil {
		return nil, fmt.Errorf("save venue: %w", err)
	}

	return venue, nil
}
