package catalogtest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/catalog"
	"github.com/iamroockie/parterre/internal/platform/identity"
)

func NewHall(t testing.TB) *catalog.Hall {
	t.Helper()

	hall, err := catalog.NewHall(catalog.HallCreateParams{
		VenueID: identity.NewUUID().String(),
		Name:    "Test hall",
		Sections: []catalog.SectionParams{
			{Name: "Default", Rows: 20, SeatsPerRow: 50},
			{Name: "VIP", Rows: 5, SeatsPerRow: 10},
		},
	})
	require.NoError(t, err)

	return hall
}
