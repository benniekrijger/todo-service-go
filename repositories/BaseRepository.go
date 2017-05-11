package repositories

import (
	"todo-service-go/cassandra"
)

type BaseRepository struct {
	Db *cassandra.Cassandra
}