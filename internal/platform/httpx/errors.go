package httpx

import (
	"fmt"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type errorBody struct {
	Error Error `json:"error"`
}

type Error struct {
	status int
	cause  error

	Code    Code              `json:"code"`
	Message string            `json:"message"`
	Errors  validation.Errors `json:"errors,omitempty"`
}

func NewError(status int, code Code, msg string, err error) Error {
	return Error{
		status: status,
		cause:  err,

		Code:    code,
		Message: msg,
		Errors:  nil,
	}
}

func (e Error) Status() int { return e.status }

func (e Error) Error() string { return e.Message }

func (e Error) Unwrap() error { return e.cause }

func BadRequestError(msg string, err error) Error {
	return NewError(http.StatusBadRequest, CodeBadRequest, msg, err)
}

func InternalError(err error) Error {
	return NewError(http.StatusInternalServerError, CodeInternalError, "internal error", err)
}

func RouteNotFoundError(route string) Error {
	msg := fmt.Sprintf("route %q not found", route)

	return NewError(http.StatusNotFound, CodeRouteNotFound, msg, nil)
}
