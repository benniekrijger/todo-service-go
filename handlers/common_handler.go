package handlers

import (
	"todo-service-go/repositories"
	"github.com/nats-io/go-nats-streaming"
)

type CommonHandler struct {
	todoRepository *repositories.TodoRepository
	natsSession    stan.Conn
}