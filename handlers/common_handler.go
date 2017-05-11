package handlers

import (
	"github.com/nats-io/go-nats"
	"todo-service-go/repositories"
)

type CommonHandler struct {
	TodoRepository *repositories.TodoRepository
	NatsSession    *nats.Conn
}