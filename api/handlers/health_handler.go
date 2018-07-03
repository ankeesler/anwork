package handlers

import (
	"log"
	"net/http"

	"github.com/ankeesler/anwork/task"
)

type healthHandler struct {
	manager task.Manager
	log     *log.Logger
}

func NewHealthHandler(manager task.Manager, log *log.Logger) http.Handler {
	return &healthHandler{manager: manager, log: log}
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
