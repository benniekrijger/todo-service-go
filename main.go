package main

import (
	"github.com/gorilla/mux"
	"todo-service-go/controllers"
	"todo-service-go/db"
	"log"
	"net/http"
)

func main() {
	session := db.Cassandra{}
	if err := session.Init("todos"); err != nil {
		panic(err)
	}

	defer session.Close()

	if err := session.CreateTable(`create table if not exists todos (
		id UUID,
		title text,
		completed boolean,
		PRIMARY KEY(id)
	)`); err != nil {
		panic(err)
	}

	personController := &controllers.TodoController{
		DbSession: &session,
	}

	router := mux.NewRouter().StrictSlash(true)
	subRouter := router.PathPrefix("/api/v1/").Subrouter()

	subRouter.HandleFunc("/todos", personController.Index).Methods("GET")
	subRouter.HandleFunc("/todos", personController.AddTodo).Methods("POST")

	log.Fatal(http.ListenAndServe(":8080", router))
}
