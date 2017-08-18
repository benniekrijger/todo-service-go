package main

import (
	"todo-service-go/controllers"
	"todo-service-go/repositories"
	"todo-service-go/handlers"
	"net/http"
	"todo-service-go/cassandra"
	"os"
	"github.com/Sirupsen/logrus"
	"todo-service-go/api"
	"github.com/nats-io/go-nats-streaming"
)

const (
	natsClientName = "service_todo"
	natsClusterName = "test-cluster"
)

func init() {
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.DebugLevel)
	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})
}

func main() {
	// Start Cassandra/Scylla client
	var (
		cassandraUrl = envString("CASSANDRA_URL", cassandra.DefaultURL)
		natsUrl = envString("NATS_URL", stan.DefaultNatsURL)
	)

	dbSession, err := cassandra.Connect(cassandraUrl, "todos")
	if err != nil {
		panic(err)
	}
	logrus.Info("Initialized DB")
	defer dbSession.Close()

	// Start NATS client
	natsSession, err := stan.Connect(natsClusterName, natsClientName, stan.NatsURL(natsUrl))
	if err != nil {
		panic(err)
	}
	logrus.Info("Initialized NATS streaming")
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

	// Start api
	a := api.NewApi(todoController)

	logrus.Fatal(http.ListenAndServe(":8080", a.Router))
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}
