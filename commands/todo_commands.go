package commands

import (
	"net/http"
	"strings"
	"errors"
)

type AddTodo struct {
	Title     string	`json:"title"`
	Completed bool		`json:"completed"`
}

func (t AddTodo) Validate(r *http.Request) error {
	if strings.Trim(t.Title, " ") == "" {
		return errors.New("No title defined")
	}

	return nil
}
