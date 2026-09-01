package identity

import (
	"errors"
	"strings"
	"uuid"
)

func NewUUID() uuid.UUID {
	return uuid.NewV7()
}

func ParseUUID(raw string) (uuid.UUID, error) {
	zero := uuid.Nil()
	trimmed := strings.TrimSpace(raw)

	parsed, err := uuid.Parse(trimmed)
	if err != nil {
		return zero, errors.New("invalid format")
	}
	if parsed == zero {
		return zero, errors.New("cannot be zero")
	}
	if parsed[6]>>4 != 7 {
		return zero, errors.New("expected version 7")
	}

	return parsed, nil
}
