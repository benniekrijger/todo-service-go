package main

import (
	"github.com/gorilla/mux"
	"todo-service-go/controllers"
	"todo-service-go/repositories"
	"todo-service-go/handlers"
	"log"
	"net/http"
	"github.com/nats-io/go-nats"
	"todo-service-go/cassandra"
	"os"
)

func main() {
	cassandraHost := os.Getenv("CASSANDRA_URL")
	if cassandraHost == "" {
		cassandraHost = cassandra.DefaultURL
	}
	dbSession, err := cassandra.Connect(cassandraHost, "todos")
	if err != nil {
		panic(err)
	}
	defer dbSession.Close()

	natsHost := os.Getenv("NATS_URL")
	if natsHost == "" {
		natsHost = nats.DefaultURL
	}
	natsSession, err := nats.Connect(natsHost)
	if err != nil {
		panic(err)
	}
	defer natsSession.Close()

	todoRepository, err := repositories.NewTodoRepository(dbSession)
	if err != nil {
		panic(err)
	}

	_, err = handlers.NewTodoHandler(todoRepository, natsSession)
	if err != nil {
		panic(err)
	}

	todoController := controllers.NewTodoController(todoRepository, natsSession)

	router := mux.NewRouter().StrictSlash(true)
	todoRouter := router.PathPrefix("/api/v1/").Subrouter()

	todoRouter.HandleFunc("/todos", todoController.Index).Methods(http.MethodGet)
	todoRouter.HandleFunc("/todos", todoController.AddTodo).Methods(http.MethodPost)
	todoRouter.HandleFunc("/todos/{todo_id}", todoController.GetTodo).Methods(http.MethodGet)
	todoRouter.HandleFunc("/todos/{todo_id}", todoController.RemoveTodo).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8080", router))
}
