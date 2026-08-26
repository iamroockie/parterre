package catalog

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"
	"uuid"

	"github.com/iamroockie/parterre/internal/platform/validation"
)

const (
	maxNameLen    = 200
	maxCityLen    = 100
	maxAddressLen = 300
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

func NewVenue(name, city, address, timezone string) (*Venue, error) {
	var verrs validation.Builder

	name = strings.TrimSpace(name)
	switch {
	case name == "":
		verrs.Add("name", "must not be empty")
	case utf8.RuneCountInString(name) > maxNameLen:
		verrs.Add("name", fmt.Sprintf("must be at most %d characters", maxNameLen))
	}

	city = strings.TrimSpace(city)
	switch {
	case city == "":
		verrs.Add("city", "must not be empty")
	case utf8.RuneCountInString(city) > maxCityLen:
		verrs.Add("city", fmt.Sprintf("must be at most %d characters", maxCityLen))
	}

	address = strings.TrimSpace(address)
	switch {
	case address == "":
		verrs.Add("address", "must not be empty")
	case utf8.RuneCountInString(address) > maxAddressLen:
		verrs.Add("address", fmt.Sprintf("must be at most %d characters", maxAddressLen))
	}

	tz, err := parseTimezone(timezone)
	if err != nil {
		verrs.Add("timezone", err.Error())
	}

	if err := verrs.Err(); err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	return &Venue{
		ID:        uuid.NewV7(),
		Name:      name,
		City:      city,
		Address:   address,
		Timezone:  tz,
		CreatedAt: now,
		UpdatedAt: now,
	}, nil
}

func parseTimezone(s string) (*time.Location, error) {
	timezone := strings.TrimSpace(s)
	if timezone == "" {
		return nil, errors.New("must not be empty")
	}
	tz, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, errors.New("invalid timezone")
	}
	if tz == time.Local {
		return nil, errors.New("must not be local")
	}

	return tz, nil
}
