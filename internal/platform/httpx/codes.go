package httpx

type Code string

const (
	CodeBadRequest      Code = "BAD_REQUEST"
	CodeEntityTooLarge  Code = "ENTITY_TOO_LARGE"
	CodeInternalError   Code = "INTERNAL_ERROR"
	CodeRouteNotFound   Code = "ROUTE_NOT_FOUND"
	CodeValidationError Code = "VALIDATION_ERROR"
)
