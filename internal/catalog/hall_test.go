package catalog_test

import (
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/catalog"
	"github.com/iamroockie/parterre/internal/platform/identity"
)

func TestHall_New(t *testing.T) {
	t.Run("valid", func(t *testing.T) {
		params := catalog.HallCreateParams{
			VenueID: identity.NewUUID().String(),
			Name:    "Test name",
			Sections: []catalog.SectionParams{
				{Name: "Default", Rows: 20, SeatsPerRow: 50},
				{Name: "VIP", Rows: 5, SeatsPerRow: 10},
			},
		}

		got, err := catalog.NewHall(params)

		require.NoError(t, err)
		require.NotNil(t, got)
		require.NotEqual(t, uuid.Nil(), got.ID)
		require.Equal(t, params.VenueID, got.VenueID.String())
		require.Equal(t, params.Name, got.Name)
		require.Equal(t, 1050, len(got.Seats), "total seats not same")
		require.Equal(t, got.CreatedAt, got.UpdatedAt)
		require.WithinDuration(t, got.CreatedAt.UTC(), got.CreatedAt, time.Microsecond)
		require.WithinDuration(t, got.UpdatedAt.UTC(), got.UpdatedAt, time.Microsecond)
	})
}

// TODO: Validation failed tests
