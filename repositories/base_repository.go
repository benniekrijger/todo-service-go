package repositories

import (
	"todo-service-go/cassandra"
)

type BaseRepository struct {
	db *cassandra.Cassandra
}