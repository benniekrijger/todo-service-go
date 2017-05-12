package utils

import (
	"net/http"
)

type InputValidation interface {
	Validate(r *http.Request) error
}