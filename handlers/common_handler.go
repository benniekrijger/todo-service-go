package handlers

import (
	"github.com/nats-io/go-nats"
	"todo-service-go/repositories"
)

type CommonHandler struct {
	todoRepository *repositories.TodoRepository
	natsSession    *nats.Conn
}