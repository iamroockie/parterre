package catalog

import (
	"context"
	"uuid"
)

type (
	GetVenueFunc    func(c context.Context, venueID uuid.UUID) (*Venue, error)
	CreateVenueFunc func(c context.Context, p VenueCreateParams) (*Venue, error)
)

type Module struct {
	GetVenue    GetVenueFunc
	CreateVenue CreateVenueFunc
}
