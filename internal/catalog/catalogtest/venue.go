package catalogtest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/catalog"
)

func NewVenue(t testing.TB) *catalog.Venue {
	t.Helper()

	venue, err := catalog.NewVenue("Test name", "Test city", "Test address", "UTC")
	require.NoError(t, err)

	return venue
}
