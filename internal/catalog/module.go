package catalog

import (
	"context"
	"uuid"
)

type (
	GetVenueFunc    func(c context.Context, venueID uuid.UUID) (*Venue, error)
	CreateVenueFunc func(c context.Context, p VenueCreateParams) (*Venue, error)

	GetHallFunc    func(c context.Context, hallID uuid.UUID) (*Hall, error)
	CreateHallFunc func(c context.Context, p HallCreateParams) (*Hall, error)
)

type Module struct {
	GetVenue    GetVenueFunc
	CreateVenue CreateVenueFunc

	GetHall    GetHallFunc
	CreateHall CreateHallFunc
}
