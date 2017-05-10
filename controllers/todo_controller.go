package controllers

import (
	"net/http"
	"encoding/json"
	"todo-service-go/models"
)

type TodoController struct {
	CommonController
}

func (c *TodoController) Index(w http.ResponseWriter, req *http.Request) {
	todos := [0]string{}

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

	c.NatsSession.Publish("todos.new", []byte("Hello!"))

	//id, err := c.DbSession.AddTodo(&todo)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//
	//todo.Id = id
	//
	//log.Printf("Inserted todo with id: %v", id)
	//
	c.SendJSON(
		w,
		req,
		todo,
		http.StatusOK,
	)
}
