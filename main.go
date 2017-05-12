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
)

func main() {
	dbConn, err := cassandra.Connect(cassandra.DefaultURL, "todos")
	if err != nil {
		panic(err)
	}

	defer dbConn.Close()

	natsSession, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	defer natsSession.Close()

	todoRepository, err := repositories.NewTodoRepository(dbConn)
	if err != nil {
		panic(err)
	}

	_, err = handlers.NewTodoHandler(todoRepository, natsSession)
	if err != nil {
		panic(err)
	}

	todoController := controllers.NewTodoController(natsSession, todoRepository)

	router := mux.NewRouter().StrictSlash(true)
	todoRouter := router.PathPrefix("/api/v1/").Subrouter()

	todoRouter.HandleFunc("/todos", todoController.Index).Methods(http.MethodGet)
	todoRouter.HandleFunc("/todos", todoController.AddTodo).Methods(http.MethodPost)
	todoRouter.HandleFunc("/todos/{todo_id}", todoController.GetTodo).Methods(http.MethodGet)
	todoRouter.HandleFunc("/todos/{todo_id}", todoController.RemoveTodo).Methods(http.MethodDelete)

	log.Fatal(http.ListenAndServe(":8080", router))
}
