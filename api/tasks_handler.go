package api

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
	h.log.Printf("Handling %s /api/v1/tasks...", r.Method)

	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *tasksHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	tasks := h.manager.Tasks()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	tasksJson, err := json.Marshal(tasks)
	if err != nil {
		h.log.Printf("Failed to marshal tasks: %s", err.Error())
		return
	}

	h.log.Printf("Returning tasks %s", tasksJson)
	_, err = w.Write(tasksJson)
	if err != nil {
		h.log.Printf("Cannot write JSON body: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (h *tasksHandler) handlePost(w http.ResponseWriter, r *http.Request) {
}
