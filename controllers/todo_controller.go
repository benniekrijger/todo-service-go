package controllers

import (
	"net/http"
	"encoding/json"
	"github.com/golang/protobuf/proto"
	"todo-service-go/models"
	"todo-service-go/events"
	"todo-service-go/repositories"
)

type TodoController struct {
	CommonController
	TodoRepository *repositories.TodoRepository
}

func (c *TodoController) Index(w http.ResponseWriter, req *http.Request) {
	todos := c.TodoRepository.GetTodos()

	c.SendJSON(
		w,
		req,
		todos,
		http.StatusOK,
	)
}

func (c *TodoController) AddTodo(w http.ResponseWriter, req *http.Request)  {
	defer req.Body.Close()

	decoder := json.NewDecoder(req.Body)
	todo := models.Todo{}
	err := decoder.Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	event := events.TodoAdded{
		Title: proto.String(todo.Title),
		Completed: proto.Bool(todo.Completed),
	}

	data, err := proto.Marshal(&event)
	if err != nil {
		http.Error(w, "Unable to add todo", http.StatusBadRequest)
		return
	}

	c.NatsSession.Publish("todos.new", data)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"ok": true}`))
}
