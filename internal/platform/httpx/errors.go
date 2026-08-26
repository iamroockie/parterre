package httpx

import (
	"errors"
	"net/http"

	"github.com/iamroockie/parterre/internal/platform/httpx/response"
	"github.com/iamroockie/parterre/internal/platform/logger"
	"github.com/iamroockie/parterre/internal/platform/validation"
)

func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	if verrs, ok := errors.AsType[validation.Errors](err); ok {
		status := http.StatusUnprocessableEntity
		response.ErrorWithFields(w, status, "Validation error", verrs.Fields())
		return
	}

	logger.FromContext(r.Context()).Error("request failed", "error", err)
	response.ErrorInternal(w)
}
