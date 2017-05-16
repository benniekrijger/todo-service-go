package api

import (
	"github.com/gorilla/mux"
	"todo-service-go/controllers"
	"net/http"
	"github.com/Sirupsen/logrus"
)

type Api struct {
	*mux.Router
}

func NewApi(todoController *controllers.TodoController) *Api {
	router := mux.NewRouter()
	api := Api{router}

	apiRouter:= router.PathPrefix("/api/v1").Subrouter().StrictSlash(true)
	apiRouter.HandleFunc("/", api.healthCheck)
	apiRouter.HandleFunc("/todos", todoController.Index).Methods(http.MethodGet)
	apiRouter.HandleFunc("/todos", todoController.AddTodo).Methods(http.MethodPost)
	apiRouter.HandleFunc("/todos/{todo_id}", todoController.GetTodo).Methods(http.MethodGet)
	apiRouter.HandleFunc("/todos/{todo_id}", todoController.RemoveTodo).Methods(http.MethodDelete)
	logrus.Info("Initialized API")

	return &api
}

func (a *Api) healthCheck(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"ok": true}`))
}