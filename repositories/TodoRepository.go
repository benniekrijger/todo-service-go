package repositories

import (
	"todo-service-go/models"
	"github.com/gocql/gocql"
	"log"
)

type TodoRepository struct {
	BaseRepository
}

func (c *TodoRepository) Init() error {
	return c.Db.CreateTable(`create table if not exists todos (
		id UUID,
		title text,
		completed boolean,
		PRIMARY KEY(id)
	)`);
}

func (c *TodoRepository) GetTodos() []models.Todo {
	var todos []models.Todo
	m := map[string]interface{}{}

	query := "SELECT id, title, completed FROM todos"
	iterable := c.Db.Connection.Query(query).Iter()
	for iterable.MapScan(m) {
		todos = append(todos, models.Todo{
			Id:		m["id"].(gocql.UUID),
			Title:	    	m["title"].(string),
			Completed:  	m["completed"].(bool),
		})

		m = map[string]interface{}{}
	}

	return todos
}

func (c *TodoRepository) AddTodo(todo *models.Todo) (gocql.UUID, error) {
	log.Println("Creating a new todo")

	// generate a unique UUID for this model
	newId := gocql.TimeUUID()

	// write data to Cassandra
	err := c.Db.Connection.Query("INSERT INTO todos (id, title, completed) VALUES (?, ?, ?)", newId, todo.Title, todo.Completed).Exec()
	if err != nil {
		return newId, err
	}

	return newId, nil
}

