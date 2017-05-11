package handlers

import (
	"github.com/nats-io/go-nats"
	"todo-service-go/events"
	"github.com/golang/protobuf/proto"
	"todo-service-go/models"
	"log"
	"github.com/gocql/gocql"
)

type TodoHandler struct {
	CommonHandler
}

func (h *TodoHandler) Init() error {
	_, err := h.NatsSession.Subscribe("todos.new", func(msg *nats.Msg) {
		h.addTodo(msg)
	})
	if err != nil {
		return err
	}

	_, err = h.NatsSession.Subscribe("todos.remove", func(msg *nats.Msg) {
		h.removeTodo(msg)
	})
	if err != nil {
		return err
	}

	return nil
}

func (h *TodoHandler) addTodo(m *nats.Msg) error {
	event := events.TodoAdded{}
	err := proto.Unmarshal(m.Data, &event)
	if err != nil {
		log.Println("Unable to unmarshal todo added event", err)
		return err
	}

	todo := models.Todo{
		Title: event.GetTitle(),
		Completed: event.GetCompleted(),
	}

	id, err := h.TodoRepository.AddTodo(&todo)
	if err != nil {
		log.Println("Unable to add todo", err)
		return err
	}

	log.Printf("Added todo with id: %s", id.String())

	return nil
}

func (h *TodoHandler) removeTodo(m *nats.Msg) error {
	event := events.TodoRemoved{}
	err := proto.Unmarshal(m.Data, &event)
	if err != nil {
		log.Println("Unable to unmarshal todo removed event", err)
		return err
	}

	id, err := gocql.ParseUUID(event.GetId())
	if err != nil {
		log.Println("Unable to unmarshal todo id", err)
		return err
	}

	err = h.TodoRepository.RemoveTodo(id)
	if err != nil {
		log.Println("Unable to remove todo", err)
		return err
	}

	log.Printf("Removed todo with id: %s", id.String())

	return nil
}