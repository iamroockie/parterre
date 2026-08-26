package catalogtest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/catalog"
)

func NewVenue(t testing.TB) *catalog.Venue {
	t.Helper()

	venue, err := catalog.NewVenue(catalog.VenueCreateParams{
		Name:     "Test name",
		City:     "Test city",
		Address:  "Test address",
		Timezone: "UTC",
	})
	require.NoError(t, err)

	return venue
}
