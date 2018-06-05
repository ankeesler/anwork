package handlers

import (
	"log"
	"net/http"
)

type healthHandler struct {
	log *log.Logger
}

func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Printf("Handling /api/v1/health...")
}
