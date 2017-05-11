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

func (c *TodoRepository) GetTodos() *[]models.Todo {
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

	return &todos
}

func (c *TodoRepository) GetTodo(id gocql.UUID) *models.Todo {
	m := map[string]interface{}{}

	query := "SELECT id, title, completed FROM todos WHERE id = ? LIMIT 1"
	iterable := c.Db.Connection.Query(query, id).Iter()
	for iterable.MapScan(m) {
		return &models.Todo{
			Id:		m["id"].(gocql.UUID),
			Title:	    	m["title"].(string),
			Completed:  	m["completed"].(bool),
		}
	}

	return nil
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

func (c *TodoRepository) RemoveTodo(id gocql.UUID) error {
	log.Printf("Removing todo with id: %s", id.String())

	err := c.Db.Connection.Query("DELETE FROM todos WHERE id = ?", id).Exec()
	if err != nil {
		return err
	}

	return nil
}

