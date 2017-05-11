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

	todoRepository := repositories.TodoRepository{repositories.BaseRepository{dbConn}}
	err = todoRepository.Init()
	if err != nil {
		panic(err)
	}

	todoHandler := &handlers.TodoHandler{handlers.CommonHandler{&todoRepository, natsSession}}
	todoHandler.Init()

	todoController := &controllers.TodoController{
		CommonController: controllers.CommonController{natsSession},
		TodoRepository: &todoRepository,
	}

	router := mux.NewRouter().StrictSlash(true)
	todoRouter := router.PathPrefix("/api/v1/").Subrouter()

	todoRouter.HandleFunc("/todos", todoController.Index).Methods("GET")
	todoRouter.HandleFunc("/todos", todoController.AddTodo).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
