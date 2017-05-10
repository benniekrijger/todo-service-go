package handlers

import (
	"todo-service-go/db"
	"github.com/nats-io/go-nats"
)

type CommonHandler struct {
	DbSession *db.Cassandra
	NatsSession *nats.Conn
}