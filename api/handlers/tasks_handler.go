package handlers

import (
	"log"
	"net/http"

	"github.com/ankeesler/anwork/task"
)

type tasksHandler struct {
	manager task.Manager
	log     *log.Logger
}

func NewTasksHandler(manager task.Manager, log *log.Logger) http.Handler {
	return &tasksHandler{manager: manager, log: log}
}

func (h *tasksHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Printf("Handling /api/v1/tasks...")

	h.manager.Tasks()
	w.WriteHeader(http.StatusOK)
}