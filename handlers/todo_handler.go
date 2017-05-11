package handlers

import (
	"github.com/nats-io/go-nats"
	"todo-service-go/events"
	"github.com/golang/protobuf/proto"
	"todo-service-go/models"
	"log"
)

type TodoHandler struct {
	CommonHandler
}

func (h *TodoHandler) Init() error {
	_, err := h.NatsSession.Subscribe("todos.new", func(m *nats.Msg) {
		event := events.TodoAdded{}
		err := proto.Unmarshal(m.Data, &event)
		if err != nil {
			log.Fatal("Unable to unmarshal todo added event", err)
		}

		todo := models.Todo{
			Title: event.GetTitle(),
			Completed: event.GetCompleted(),
		}

		id, err := h.TodoRepository.AddTodo(&todo)
		if err != nil {
			log.Fatal("Unable to add todo", err)
		}

		log.Printf("Added todo with id: %s", id.String())
	})
	if err != nil {
		return err
	}

	return nil
}