package utils

import (
	"net/http"
	"errors"
)

type InputValidation interface {
	Validate(r *http.Request) error
}

var (
	// ErrInvalidUUID - error when we have a UUID validation issue
	ErrInvalidUUID = errors.New("invalid uuid")
)
