package repositories

import (
	"todo-service-go/models"
	"github.com/gocql/gocql"
	"log"
	"todo-service-go/cassandra"
)

type TodoRepository struct {
	BaseRepository
}

func NewTodoRepository(db *cassandra.Cassandra) (*TodoRepository, error) {
	err := db.CreateTable(`create table if not exists todos (
		id UUID,
		title text,
		completed boolean,
		PRIMARY KEY(id)
	)`)

	if err != nil {
		return nil, err
	}

	return &TodoRepository{BaseRepository{db}}, nil
}

func (c *TodoRepository) GetTodos() *[]models.Todo {
	var todos []models.Todo
	m := map[string]interface{}{}

	query := "SELECT id, title, completed FROM todos"
	iterable := c.db.Connection.Query(query).Iter()
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
	iterable := c.db.Connection.Query(query, id).Iter()
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

	// write data to Cassandra
	err := c.db.Connection.Query("INSERT INTO todos (id, title, completed) VALUES (?, ?, ?)", todo.Id, todo.Title, todo.Completed).Exec()
	if err != nil {
		return todo.Id, err
	}

	return todo.Id, nil
}

func (c *TodoRepository) RemoveTodo(id gocql.UUID) error {
	log.Printf("Removing todo with id: %s", id.String())

	err := c.db.Connection.Query("DELETE FROM todos WHERE id = ?", id).Exec()
	if err != nil {
		return err
	}

	return nil
}

