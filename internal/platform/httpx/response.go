package httpx

import (
	"encoding/json/v2"
	"errors"
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"github.com/iamroockie/parterre/internal/platform/logger"
)

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	resp := classify(err)

	if resp.cause != nil && resp.status >= 500 {
		log := logger.FromContext(r.Context())
		log.Error("request failed", "code", resp.Code, "error", resp.cause)
	}

	WriteJSON(w, resp.status, errorBody{resp})
}

func WriteJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(payload)
	if err != nil {
		errResp := InternalError(err)
		status = errResp.status
		data, _ = json.Marshal(errorBody{errResp})
	}

	w.WriteHeader(status)
	_, _ = w.Write(data)
}

func classify(err error) Error {
	if mbe, ok := errors.AsType[*http.MaxBytesError](err); ok {
		return NewError(
			http.StatusRequestEntityTooLarge,
			CodeEntityTooLarge,
			fmt.Sprintf("request too large, max %d bytes", mbe.Limit),
			nil,
		)
	}

	if e, ok := errors.AsType[Error](err); ok {
		return e
	}

	if ve, ok := errors.AsType[validation.Errors](err); ok {
		return Error{
			status:  http.StatusUnprocessableEntity,
			cause:   nil,
			Code:    CodeValidationError,
			Message: "validation failed",
			Errors:  ve,
		}
	}

	return InternalError(err)
}
