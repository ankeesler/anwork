package api

import (
	"log"
	"net/http"
)

type eventsHandler struct {
	log *log.Logger
}

func (h *eventsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Printf("Handling /api/v1/events...")
}
