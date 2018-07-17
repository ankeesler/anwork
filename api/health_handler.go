package api

import (
	"log"
	"net/http"
)

type healthHandler struct {
	log *log.Logger
}

func NewHealthHandler(log *log.Logger) http.Handler {
	return &healthHandler{log: log}
}

func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Printf("Handling %s /api/v1/health...", r.Method)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	h.log.Printf("Returning healthy...")
	w.WriteHeader(http.StatusNoContent)
}
