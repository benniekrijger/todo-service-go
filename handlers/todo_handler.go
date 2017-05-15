package handlers

import (
	"github.com/nats-io/go-nats"
	"todo-service-go/events"
	"github.com/golang/protobuf/proto"
	"todo-service-go/models"
	"github.com/gocql/gocql"
	"todo-service-go/repositories"
	"github.com/Sirupsen/logrus"
)

type TodoHandler struct {
	CommonHandler
}

func NewTodoHandler(todoRepository *repositories.TodoRepository, natsSession *nats.Conn) (*TodoHandler, error) {
	handler := TodoHandler{CommonHandler{todoRepository, natsSession}}

	_, err := natsSession.Subscribe("todos.new", func(msg *nats.Msg) {
		handler.addTodo(msg)
	})
	if err != nil {
		return nil, err
	}

	_, err = natsSession.Subscribe("todos.remove", func(msg *nats.Msg) {
		handler.removeTodo(msg)
	})
	if err != nil {
		return nil, err
	}

	return &handler, nil
}

func (h *TodoHandler) addTodo(m *nats.Msg) error {
	event := events.TodoAdded{}
	err := proto.Unmarshal(m.Data, &event)
	if err != nil {
		logrus.Info("Unable to unmarshal todo added event", err)
		return err
	}

	id, err := gocql.ParseUUID(event.GetId())
	if err != nil {
		logrus.Info("Unable to parse todo id", err)
		return err
	}

	todo := models.Todo{
		Id: id,
		Title: event.GetTitle(),
		Completed: event.GetCompleted(),
	}

	_, err = h.todoRepository.AddTodo(&todo)
	if err != nil {
		logrus.Info("Unable to add todo", err)
		return err
	}

	logrus.Printf("Added todo with id: %s", id.String())

	return nil
}

func (h *TodoHandler) removeTodo(m *nats.Msg) error {
	event := events.TodoRemoved{}
	err := proto.Unmarshal(m.Data, &event)
	if err != nil {
		logrus.Info("Unable to unmarshal todo removed event", err)
		return err
	}

	id, err := gocql.ParseUUID(event.GetId())
	if err != nil {
		logrus.Info("Unable to unmarshal todo id", err)
		return err
	}

	err = h.todoRepository.RemoveTodo(id)
	if err != nil {
		logrus.Info("Unable to remove todo", err)
		return err
	}

	logrus.Printf("Removed todo with id: %s", id.String())

	return nil
}