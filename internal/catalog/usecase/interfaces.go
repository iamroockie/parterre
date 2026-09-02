package usecase

import (
	"context"
	"uuid"

	"github.com/iamroockie/parterre/internal/catalog"
)

type venueSaver interface {
	Save(context.Context, *catalog.Venue) error
}

type venueGetter interface {
	Get(context.Context, uuid.UUID) (*catalog.Venue, error)
}

type hallSaver interface {
	Save(context.Context, *catalog.Hall) error
}

type hallGetter interface {
	Get(context.Context, uuid.UUID) (*catalog.Hall, error)
}
