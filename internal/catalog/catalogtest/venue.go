package catalogtest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/catalog"
)

func Venue(t testing.TB) *catalog.Venue {
	t.Helper()

	venue, err := catalog.NewVenue(VenueCreateParams(t))
	require.NoError(t, err)

	return venue
}

func VenueCreateParams(t testing.TB) catalog.VenueCreateParams {
	t.Helper()

	return catalog.VenueCreateParams{
		Name:     "Test name",
		City:     "Test city",
		Address:  "Test address",
		Timezone: "UTC",
	}
}
