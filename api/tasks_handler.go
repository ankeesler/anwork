package api

import (
	"log"
	"net/http"
)

type tasksHandler struct {
	log *log.Logger
}

func (h *tasksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Printf("Handling /api/v1/tasks...")
}
