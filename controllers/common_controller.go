package controllers

import (
	"net/http"
	"encoding/json"
	"log"
	"fmt"
	"io"
)

type CommonController struct {}

func (c *CommonController) SendJSON(w http.ResponseWriter, r *http.Request, v interface{}, code int) {
	w.Header().Add("Content-Type", "application/json")
	b, err := json.Marshal(v)
	if err != nil {
		log.Print(fmt.Sprintf("Error while encoding JSON: %v", err))
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, `{"error": "Internal server error"}`)
	} else {
		w.WriteHeader(code)
		io.WriteString(w, string(b))
	}
}
