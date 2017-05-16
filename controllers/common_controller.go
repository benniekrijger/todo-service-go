package controllers

import (
	"net/http"
	"encoding/json"
	"fmt"
	"io"
	"todo-service-go/utils"
	"github.com/Sirupsen/logrus"
	"github.com/nats-io/go-nats-streaming"
)

type CommonController struct {
	natsSession stan.Conn
}

func (c *CommonController) decodeAndValidate(r *http.Request, v utils.InputValidation) error {
	// json decode the payload - obviously this could be abstracted
	// to handle many content types
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	defer r.Body.Close()
	// peform validation on the InputValidation implementation
	return v.Validate(r)
}

func (c *CommonController) sendJSON(w http.ResponseWriter, r *http.Request, v interface{}, code int) {
	w.Header().Add("Content-Type", "application/json")
	b, err := json.Marshal(v)
	if err != nil {
		logrus.Print(fmt.Sprintf("Error while encoding JSON: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Internal server error"}`)
	} else {
		w.WriteHeader(code)
		io.WriteString(w, string(b))
	}
}
