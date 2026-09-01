package identity_test

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/identity"
)

func TestUUID(t *testing.T) {
	zero := uuid.Nil()
	t.Run("valid", func(t *testing.T) {
		id := identity.NewUUID()

		got, err := identity.ParseUUID(id.String())

		require.NoError(t, err)
		require.Equal(t, id, got)
	})

	t.Run("invalid", func(t *testing.T) {
		got, err := identity.ParseUUID("00-11-22")

		require.Error(t, err)
		require.ErrorContains(t, err, "invalid format")
		require.Equal(t, zero, got)
	})

	t.Run("nil", func(t *testing.T) {
		got, err := identity.ParseUUID(zero.String())

		require.Error(t, err)
		require.ErrorContains(t, err, "cannot be zero")
		require.Equal(t, zero, got)
	})

	t.Run("v4", func(t *testing.T) {
		got, err := identity.ParseUUID(uuid.NewV4().String())

		require.Error(t, err)
		require.ErrorContains(t, err, "expected version 7")
		require.Equal(t, zero, got)
	})
}
