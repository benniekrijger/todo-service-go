package controllers

import (
	"net/http"
	"gosite/db"
	"encoding/json"
	"gosite/models"
	"log"
)

type TodoController struct {
	CommonController
	DbSession *db.Cassandra
}

func (c *TodoController) Index(w http.ResponseWriter, req *http.Request) {
	todos := c.DbSession.GetTodos()

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

	id, err := c.DbSession.AddTodo(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	todo.Id = id

	log.Printf("Inserted todo with id: %v", id)

	c.SendJSON(
		w,
		req,
		todo,
		http.StatusOK,
	)
}
