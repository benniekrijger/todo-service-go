package main

import (
	"github.com/gorilla/mux"
	"todo-service-go/controllers"
	"todo-service-go/db"
	"log"
	"net/http"
	"todo-service-go/handlers"
	"github.com/nats-io/go-nats"
)

func main() {
	dbSession := db.Cassandra{}
	if err := dbSession.Init("todos"); err != nil {
		panic(err)
	}

	defer dbSession.Close()

	if err := dbSession.CreateTable(`create table if not exists todos (
		id UUID,
		title text,
		completed boolean,
		PRIMARY KEY(id)
	)`); err != nil {
		panic(err)
	}

	natsSession, err := nats.Connect(nats.DefaultURL)
	if err != nil {
		panic(err)
	}

	defer natsSession.Close()

	todoHandler := &handlers.TodoHandler{handlers.CommonHandler{&dbSession, natsSession}}
	todoHandler.Init()

	todoController := &controllers.TodoController{controllers.CommonController{natsSession}}

	router := mux.NewRouter().StrictSlash(true)
	todoRouter := router.PathPrefix("/api/v1/").Subrouter()

	todoRouter.HandleFunc("/todos", todoController.Index).Methods("GET")
	todoRouter.HandleFunc("/todos", todoController.AddTodo).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
