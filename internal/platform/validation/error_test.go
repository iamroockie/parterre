package validation_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/validation"
)

func TestBuilderEmpty(t *testing.T) {
	var b validation.Builder

	require.NoError(t, b.Err())
}

func TestBuilderCollectsEveryField(t *testing.T) {
	var b validation.Builder
	b.Add("name", "must not be empty")
	b.Add("timezone", "invalid timezone")

	err := b.Err()
	require.Error(t, err)

	verrs, ok := errors.AsType[validation.Errors](err)
	require.True(t, ok, "expected validation.Errors, got %T", err)
	require.Len(t, verrs, 2)
	require.Equal(t, map[string]string{
		"name":     "must not be empty",
		"timezone": "invalid timezone",
	}, verrs.Fields())
}

func TestErrorsSurviveWrapping(t *testing.T) {
	var b validation.Builder
	b.Add("name", "must not be empty")

	wrapped := fmt.Errorf("create venue: %w", b.Err())

	verrs, ok := errors.AsType[validation.Errors](wrapped)
	require.True(t, ok, "AsType should find Errors through %%w, got %T", wrapped)
	require.Equal(t, map[string]string{"name": "must not be empty"}, verrs.Fields())
}

func TestErrorsErrorNamesFields(t *testing.T) {
	var b validation.Builder
	b.Add("name", "must not be empty")
	b.Add("timezone", "invalid timezone")

	msg := b.Err().Error()
	require.Contains(t, msg, "name: must not be empty")
	require.Contains(t, msg, "timezone: invalid timezone")
}
