package controllers

import (
	"net/http"
	"github.com/golang/protobuf/proto"
	"todo-service-go/events"
	"todo-service-go/repositories"
	"todo-service-go/commands"
	"github.com/gorilla/mux"
	"github.com/gocql/gocql"
	"github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats-streaming"
)

type TodoController struct {
	CommonController
	todoRepository *repositories.TodoRepository
}

func NewTodoController(todoRepository *repositories.TodoRepository, natsSession stan.Conn) *TodoController {
	return &TodoController{
		CommonController{natsSession},
		todoRepository,
	}
}

func (c *TodoController) Index(w http.ResponseWriter, req *http.Request) {
	todos := c.todoRepository.GetTodos()

	c.sendJSON(
		w,
		req,
		todos,
		http.StatusOK,
	)
}

func (c *TodoController) RemoveTodo(w http.ResponseWriter, req *http.Request)  {
	vars := mux.Vars(req)
	id := vars["todo_id"]

	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		logrus.Info(err)
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	event := events.TodoRemoved{
		Id: uuid.String(),
	}

	data, err := proto.Marshal(&event)
	if err != nil {
		http.Error(w, "Unable to add todo", http.StatusInternalServerError)
		return
	}

	c.natsSession.Publish("todos.remove", data)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"ok": true}`))
}

func (c *TodoController) GetTodo(w http.ResponseWriter, req *http.Request)  {
	vars := mux.Vars(req)
	id := vars["todo_id"]

	uuid, err := gocql.ParseUUID(id)
	if err != nil {
		logrus.Info(err)
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	todo := c.todoRepository.GetTodo(uuid)
	if todo != nil {
		c.sendJSON(
			w,
			req,
			todo,
			http.StatusOK,
		)
	} else {
		http.Error(w, "Todo not found", http.StatusNotFound)
	}
}

func (c *TodoController) AddTodo(w http.ResponseWriter, req *http.Request)  {
	defer req.Body.Close()

	command := &commands.AddTodo{}
	err := c.decodeAndValidate(req, command)
	if err != nil {
		logrus.Info(err)
		http.Error(w, "Invalid body", http.StatusBadRequest)
		return
	}

	newId, err := gocql.RandomUUID()
	if err != nil {
		logrus.Info(err)
		http.Error(w, "Unknown error", http.StatusInternalServerError)
		return
	}

	event := events.TodoAdded{
		Id: newId.String(),
		Title: command.Title,
		Completed: command.Completed,
	}

	data, err := proto.Marshal(&event)
	if err != nil {
		http.Error(w, "Unable to add todo", http.StatusBadRequest)
		return
	}

	c.natsSession.Publish("todos.new", data)

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"ok": true}`))
}
