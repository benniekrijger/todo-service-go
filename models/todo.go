package models

import (
	"github.com/gocql/gocql"
)

type Todo struct {
	Id        gocql.UUID	`json:"id"`
	Title     string	`json:"title"`
	Completed bool		`json:"completed"`
}
