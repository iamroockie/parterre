package catalog_test

import (
	"errors"
	"strings"
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/catalog"
	"github.com/iamroockie/parterre/internal/platform/validation"
)

func TestVenue_New(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		got, err := catalog.NewVenue(catalog.VenueCreateParams{
			Name:     " name\n",
			City:     "  city",
			Address:  "address  ",
			Timezone: "\tEurope/Moscow",
		})

		require.NoError(t, err)
		require.NotEqual(t, uuid.Nil(), got.ID)
		require.Equal(t, "name", got.Name)
		require.Equal(t, "city", got.City)
		require.Equal(t, "address", got.Address)
		require.Equal(t, "Europe/Moscow", got.Timezone.String())
		require.NotEqual(t, time.Time{}, got.CreatedAt)
		require.Equal(t, got.CreatedAt, got.UpdatedAt)
	})

	t.Run("fail", func(t *testing.T) {
		const (
			validName    = "Площадка"
			validCity    = "Москва"
			validAddress = "Ленина 1"
			validTZ      = "Europe/Moscow"
		)

		tests := map[string]struct {
			name, city, address, tz string
			wantFields              map[string]string
		}{
			"empty name": {
				name: "\n   \t", city: validCity, address: validAddress, tz: validTZ,
				wantFields: map[string]string{"name": "must not be empty"},
			},
			"name too long": {
				name: strings.Repeat("И", 201), city: validCity, address: validAddress, tz: validTZ,
				wantFields: map[string]string{"name": "must be at most 200 characters"},
			},
			"empty city": {
				name: validName, city: "   \t", address: validAddress, tz: validTZ,
				wantFields: map[string]string{"city": "must not be empty"},
			},
			"city too long": {
				name: validName, city: strings.Repeat("И", 101), address: validAddress, tz: validTZ,
				wantFields: map[string]string{"city": "must be at most 100 characters"},
			},
			"empty address": {
				name: validName, city: validCity, address: "\n ", tz: validTZ,
				wantFields: map[string]string{"address": "must not be empty"},
			},
			"address too long": {
				name: validName, city: validCity, address: strings.Repeat("И", 301), tz: validTZ,
				wantFields: map[string]string{"address": "must be at most 300 characters"},
			},
			"empty timezone": {
				name: validName, city: validCity, address: validAddress, tz: "\t ",
				wantFields: map[string]string{"timezone": "must not be empty"},
			},
			"local timezone": {
				name: validName, city: validCity, address: validAddress, tz: "Local",
				wantFields: map[string]string{"timezone": "must not be local"},
			},
			"unknown timezone": {
				name: validName, city: validCity, address: validAddress, tz: "Mars/Olympus",
				wantFields: map[string]string{"timezone": "invalid timezone"},
			},
			"every field": {
				name: " ", city: "", address: "\t", tz: "Mars/Olympus",
				wantFields: map[string]string{
					"name":     "must not be empty",
					"city":     "must not be empty",
					"address":  "must not be empty",
					"timezone": "invalid timezone",
				},
			},
		}

		for name, tt := range tests {
			t.Run(name, func(t *testing.T) {
				got, err := catalog.NewVenue(catalog.VenueCreateParams{
					Name:     tt.name,
					City:     tt.city,
					Address:  tt.address,
					Timezone: tt.tz,
				})

				require.Nil(t, got)

				verrs, ok := errors.AsType[validation.Errors](err)
				require.True(t, ok, "ожидался validation.Errors, получено %T", err)
				require.Equal(t, tt.wantFields, verrs.Fields())
			})
		}
	})
}
