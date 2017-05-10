package db

import (
	"github.com/gocql/gocql"
	"todo-service-go/models"
	"log"
)

type Cassandra struct {
	cluster *gocql.ClusterConfig
	session *gocql.Session
}

func (c *Cassandra) Init(keyspace string) error {
	var err error

	c.cluster = gocql.NewCluster("127.0.0.1")
	c.cluster.Keyspace = keyspace
	if err != nil {
		return err
	}

	c.session, err = c.cluster.CreateSession()
	if err != nil {
		return err
	}

	log.Println("cassandra init done")

	return nil
}

func (c *Cassandra) CreateTable(table string) error {
	if err := c.session.Query(table).RetryPolicy(nil).Exec(); err != nil {
		log.Printf("error creating table table=%q err=%v\n", table, err)
		return err
	}

	return nil
}

func (c *Cassandra) Close() {
	log.Println("cassandra connection closed")
	c.session.Close()
}

func (c *Cassandra) GetTodos() []models.Todo {
	var todos []models.Todo
	m := map[string]interface{}{}

	query := "SELECT id, title, completed FROM todos"
	iterable := c.session.Query(query).Iter()
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

func (c *Cassandra) AddTodo(todo *models.Todo) (gocql.UUID, error) {
	log.Println("creating a new todo")

	// generate a unique UUID for this model
	newId := gocql.TimeUUID()

	// write data to Cassandra
	if err := c.session.Query("INSERT INTO todos (id, title, completed) VALUES (?, ?, ?)", newId, todo.Title, todo.Completed).Exec(); err != nil {
		return newId, err
	}

	return newId, nil
}
