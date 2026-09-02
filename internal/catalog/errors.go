package catalog

import "errors"

var (
	ErrHallNameTaken = errors.New("hall name taken")
	ErrHallNotFound  = errors.New("hall not found")
	ErrHallSeatTaken = errors.New("hall seat taken")
	ErrVenueNotFound = errors.New("venue not found")
)
