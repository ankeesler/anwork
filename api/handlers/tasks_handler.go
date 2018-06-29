package handlers

import (
	"encoding/json"
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

	tasks := h.manager.Tasks()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	tasksJson, err := json.Marshal(tasks)
	if err != nil {
		h.log.Printf("Failed to marshal tasks: %s", err.Error())
		return
	}

	w.Write(tasksJson)
}
