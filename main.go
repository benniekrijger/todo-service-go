package main

import (
	"github.com/gorilla/mux"
	"todo-service-go/controllers"
	"todo-service-go/repositories"
	"todo-service-go/handlers"
	"net/http"
	"github.com/nats-io/go-nats"
	"todo-service-go/cassandra"
	"os"
	"github.com/Sirupsen/logrus"
)

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

func main() {
	// Start Cassandra/Scylla client
	cassandraUrl := os.Getenv("CASSANDRA_URL")
	if cassandraUrl == "" {
		cassandraUrl = cassandra.DefaultURL
	}
	dbSession, err := cassandra.Connect(cassandraUrl, "todos")
	if err != nil {
		panic(err)
	}
	logrus.Info("Initialized DB")
	defer dbSession.Close()

	// Start NATS client
	natsUrl := os.Getenv("NATS_URL")
	if natsUrl == "" {
		natsUrl = nats.DefaultURL
	}
	natsSession, err := nats.Connect(natsUrl)
	if err != nil {
		panic(err)
	}
	logrus.Info("Initialized NATS")
	defer natsSession.Close()

	// Start repository
	todoRepository, err := repositories.NewTodoRepository(dbSession)
	if err != nil {
		panic(err)
	}
	logrus.Info("Initialized Repositories")

	// Start event handler
	_, err = handlers.NewTodoHandler(todoRepository, natsSession)
	if err != nil {
		panic(err)
	}
	logrus.Info("Initialized Handlers")

	// Start controller
	todoController := controllers.NewTodoController(todoRepository, natsSession)
	logrus.Info("Initialized Controllers")

	// Start routers
	router := mux.NewRouter()
	apiRouter:= router.PathPrefix("/api/v1").Subrouter().StrictSlash(true)
	apiRouter.HandleFunc("/", healthCheck)
	apiRouter.HandleFunc("/todos", todoController.Index).Methods(http.MethodGet)
	apiRouter.HandleFunc("/todos", todoController.AddTodo).Methods(http.MethodPost)
	apiRouter.HandleFunc("/todos/{todo_id}", todoController.GetTodo).Methods(http.MethodGet)
	apiRouter.HandleFunc("/todos/{todo_id}", todoController.RemoveTodo).Methods(http.MethodDelete)
	logrus.Info("Initialized API")

	logrus.Fatal(http.ListenAndServe(":8080", router))
}

func healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ok": true}`))
}
