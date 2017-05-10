package handlers

import (
	"github.com/nats-io/go-nats"
	"fmt"
)

type TodoHandler struct {
	CommonHandler
}

func (h *TodoHandler) Init() error {
	_, err := h.NatsSession.Subscribe("todos.new", func(m *nats.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})

	if err != nil {
		return err
	}

	return nil
}