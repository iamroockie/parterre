package httpx_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/httpx"
)

func TestError(t *testing.T) {
	t.Run("new", func(t *testing.T) {
		e := errors.New("test conflict")
		code := httpx.Code("CONFLICT")
		status := http.StatusConflict
		msg := "test message"

		err := httpx.NewError(status, code, msg, e)

		require.ErrorIs(t, err, e)
		require.Equal(t, msg, err.Error())
		require.Equal(t, status, err.Status())
	})

	t.Run("bad request", func(t *testing.T) {
		msg := "bad request"
		e := errors.New("test bad request")

		err := httpx.BadRequestError(msg, e)

		require.ErrorIs(t, err, e)
		require.Equal(t, msg, err.Error())
		require.Equal(t, http.StatusBadRequest, err.Status())
		require.Equal(t, msg, err.Message)
	})

	t.Run("internal server error", func(t *testing.T) {
		e := errors.New("test internal error")

		err := httpx.InternalError(e)

		require.ErrorIs(t, err, e)
		require.Equal(t, "internal error", err.Error())
		require.Equal(t, http.StatusInternalServerError, err.Status())
	})

	t.Run("route not found", func(t *testing.T) {
		err := httpx.RouteNotFoundError("api/unknwown")

		require.Equal(t, `route "api/unknwown" not found`, err.Error())
		require.Equal(t, http.StatusNotFound, err.Status())
	})
}
