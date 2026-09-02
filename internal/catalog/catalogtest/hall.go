package catalogtest

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/catalog"
	"github.com/iamroockie/parterre/internal/platform/identity"
)

func Hall(t testing.TB) *catalog.Hall {
	t.Helper()

	hall, err := catalog.NewHall(HallCreateParams(t))
	require.NoError(t, err)

	return hall
}

func HallCreateParams(t testing.TB) catalog.HallCreateParams {
	t.Helper()

	return catalog.HallCreateParams{
		VenueID: identity.NewUUID(),
		Name:    "Test hall",
		Sections: []catalog.Section{
			{Name: "Default", Rows: 20, SeatsPerRow: 50},
			{Name: "VIP", Rows: 5, SeatsPerRow: 10},
		},
	}
}
