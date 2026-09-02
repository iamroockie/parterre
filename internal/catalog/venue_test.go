package catalog_test

import (
	"strings"
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/catalog"
)

func TestVenue_New(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		params := catalog.VenueCreateParams{
			Name:     "Test name",
			City:     "Test city",
			Address:  "Test address",
			Timezone: "UTC",
		}

		got, err := catalog.NewVenue(params)

		require.NoError(t, err)
		require.NotNil(t, got)
		require.NotEqual(t, uuid.Nil(), got.ID)
		require.Equal(t, params.Name, got.Name)
		require.Equal(t, params.City, got.City)
		require.Equal(t, params.Address, got.Address)
		require.Equal(t, params.Timezone, got.Timezone.String())
		require.Equal(t, got.CreatedAt, got.UpdatedAt)
		require.WithinDuration(t, got.CreatedAt.UTC(), got.CreatedAt, time.Microsecond)
		require.WithinDuration(t, got.UpdatedAt.UTC(), got.UpdatedAt, time.Microsecond)
	})

	t.Run("trims fields", func(t *testing.T) {
		params := venueParams("  Test name  ", "  Test city  ", "  Test address  ", "  UTC  ")

		got, err := catalog.NewVenue(params)

		require.NoError(t, err)
		require.Equal(t, "Test name", got.Name)
		require.Equal(t, "Test city", got.City)
		require.Equal(t, "Test address", got.Address)
		require.Equal(t, "UTC", got.Timezone.String())
	})

	t.Run("accepts boundary lengths", func(t *testing.T) {
		params := venueParams(
			strings.Repeat("a", 100),
			strings.Repeat("b", 150),
			strings.Repeat("c", 200),
			"Europe/Moscow",
		)

		got, err := catalog.NewVenue(params)

		require.NoError(t, err)
		require.Equal(t, "Europe/Moscow", got.Timezone.String())
	})

	t.Run("counts runes, not bytes", func(t *testing.T) {
		params := venueParams(
			strings.Repeat("я", 100),
			strings.Repeat("я", 150),
			strings.Repeat("я", 200),
			"UTC",
		)

		got, err := catalog.NewVenue(params)

		require.NoError(t, err)
		require.Equal(t, strings.Repeat("я", 100), got.Name)
	})
}

func TestVenue_New_InvalidParams(t *testing.T) {
	const (
		validName     = "Test name"
		validCity     = "Test city"
		validAddress  = "Test address"
		validTimezone = "UTC"
	)

	longName := strings.Repeat("a", 101)
	longCity := strings.Repeat("b", 151)
	longAddress := strings.Repeat("c", 201)

	tests := map[string]struct {
		params    catalog.VenueCreateParams
		wantField string
	}{
		"empty name": {
			params:    venueParams("   ", validCity, validAddress, validTimezone),
			wantField: "name",
		},
		"too short name": {
			params:    venueParams("A", validCity, validAddress, validTimezone),
			wantField: "name",
		},
		"too long name": {
			params:    venueParams(longName, validCity, validAddress, validTimezone),
			wantField: "name",
		},
		"empty city": {
			params:    venueParams(validName, "   ", validAddress, validTimezone),
			wantField: "city",
		},
		"too short city": {
			params:    venueParams(validName, "A", validAddress, validTimezone),
			wantField: "city",
		},
		"too long city": {
			params:    venueParams(validName, longCity, validAddress, validTimezone),
			wantField: "city",
		},
		"empty address": {
			params:    venueParams(validName, validCity, "   ", validTimezone),
			wantField: "address",
		},
		"too short address": {
			params:    venueParams(validName, validCity, "A", validTimezone),
			wantField: "address",
		},
		"too long address": {
			params:    venueParams(validName, validCity, longAddress, validTimezone),
			wantField: "address",
		},
		"empty timezone": {
			params:    venueParams(validName, validCity, validAddress, "   "),
			wantField: "timezone",
		},
		"unknown timezone": {
			params:    venueParams(validName, validCity, validAddress, "Mars/Olympus"),
			wantField: "timezone",
		},
		"local timezone": {
			params:    venueParams(validName, validCity, validAddress, "Local"),
			wantField: "timezone",
		},
		"every field invalid": {
			params:    venueParams("", "", "", ""),
			wantField: "address",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := catalog.NewVenue(tt.params)

			require.Nil(t, got)
			require.ErrorContains(t, err, tt.wantField)
		})
	}
}

func venueParams(name, city, address, timezone string) catalog.VenueCreateParams {
	return catalog.VenueCreateParams{
		Name:     name,
		City:     city,
		Address:  address,
		Timezone: timezone,
	}
}
