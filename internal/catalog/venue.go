package catalog

import (
	"errors"
	"strings"
	"time"
	"uuid"

	v "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/iamroockie/parterre/internal/platform/identity"
)

type Venue struct {
	ID        uuid.UUID
	Name      string
	City      string
	Address   string
	Timezone  *time.Location
	CreatedAt time.Time
	UpdatedAt time.Time
}

type VenueCreateParams struct {
	Name     string
	City     string
	Address  string
	Timezone string
}

func NewVenue(p VenueCreateParams) (*Venue, error) {
	name := strings.TrimSpace(p.Name)
	city := strings.TrimSpace(p.City)
	address := strings.TrimSpace(p.Address)
	timezone, tzErr := parseTimezone(strings.TrimSpace(p.Timezone))

	err := v.Errors{
		"name":     v.Validate(name, v.Required, v.RuneLength(2, 100)),
		"city":     v.Validate(city, v.Required, v.RuneLength(2, 150)),
		"address":  v.Validate(address, v.Required, v.RuneLength(2, 200)),
		"timezone": tzErr,
	}.Filter()
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()

	return &Venue{
		ID:        identity.NewUUID(),
		Name:      name,
		City:      city,
		Address:   address,
		Timezone:  timezone,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func parseTimezone(raw string) (*time.Location, error) {
	if raw == "" {
		return nil, errors.New("cannot be empty")
	}
	tz, err := time.LoadLocation(raw)
	if err != nil {
		return nil, errors.New("invalid format")
	}
	if tz == time.Local {
		return nil, errors.New("cannot be local")
	}

	return tz, nil
}
