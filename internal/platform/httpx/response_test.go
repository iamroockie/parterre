package httpx_test

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/stretchr/testify/require"

	"github.com/iamroockie/parterre/internal/platform/httpx"
	"github.com/iamroockie/parterre/internal/platform/httpx/httpxtest"
	"github.com/iamroockie/parterre/internal/platform/logger"
	"github.com/iamroockie/parterre/internal/platform/logger/loggertest"
)

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()

	httpx.WriteJSON(w, http.StatusCreated, map[string]string{"id": "1"})

	require.Equal(t, http.StatusCreated, w.Code)
	require.Equal(t, "application/json", w.Header().Get("Content-Type"))
	require.JSONEq(t, `{"id":"1"}`, w.Body.String())
}

func TestWriteJSON_NotSerializable(t *testing.T) {
	w := httptest.NewRecorder()

	httpx.WriteJSON(w, http.StatusCreated, make(chan int))

	require.Equal(t, http.StatusInternalServerError, w.Code)
	resp := httpxtest.DecodeErrorResponse(t, w)
	require.Equal(t, httpx.CodeInternalError, resp.Error.Code)
}

func TestWriteError_WithLog(t *testing.T) {
	tests := map[string]struct {
		err      error
		wantCode httpx.Code
	}{
		"internal": {
			err:      httpx.InternalError(errors.New("test")),
			wantCode: httpx.CodeInternalError,
		},
		"unknnown": {
			err:      errors.New("unknnown error"),
			wantCode: httpx.CodeInternalError,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			log, buf := loggertest.NewLogger(t, slog.LevelDebug)
			ctx := logger.NewContext(t.Context(), log)
			w := httptest.NewRecorder()
			r := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)

			httpx.WriteError(w, r, tt.err)

			logs := loggertest.Logs(t, buf)
			require.Len(t, logs, 1)
			msg := logs[0]
			require.Equal(t, slog.LevelError.String(), msg["level"])
			require.Equal(t, string(tt.wantCode), msg["code"])
			resp := httpxtest.DecodeErrorResponse(t, w)
			require.Equal(t, httpx.CodeInternalError, resp.Error.Code)
		})
	}
}

func TestWriteError_WithoutLog(t *testing.T) {
	log, buf := loggertest.NewLogger(t, slog.LevelDebug)
	ctx := logger.NewContext(t.Context(), log)
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)

	httpx.WriteError(w, r, httpx.InternalError(nil))

	logs := loggertest.Logs(t, buf)
	require.Len(t, logs, 0)
	resp := httpxtest.DecodeErrorResponse(t, w)
	require.Equal(t, httpx.CodeInternalError, resp.Error.Code)
}

func TestWriteError_ValidationError(t *testing.T) {
	log, buf := loggertest.NewLogger(t, slog.LevelDebug)
	ctx := logger.NewContext(t.Context(), log)
	w := httptest.NewRecorder()
	r := httptest.NewRequestWithContext(ctx, http.MethodGet, "/", nil)
	wantErrs := map[string]string{
		"first":  "first error",
		"second": "second error",
	}
	errs := validation.Errors{
		"first":  errors.New("first error"),
		"second": errors.New("second error"),
	}.Filter()

	httpx.WriteError(w, r, errs)

	logs := loggertest.Logs(t, buf)
	require.Len(t, logs, 0)
	resp := httpxtest.DecodeErrorResponse(t, w)
	require.Equal(t, httpx.CodeValidationError, resp.Error.Code)
	require.Equal(t, wantErrs, resp.Error.Errors)
}
