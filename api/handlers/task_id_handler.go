package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ankeesler/anwork/task"
)

type taskIDHandler struct {
	manager task.Manager
	log     *log.Logger
}

func NewTaskIDHandler(manager task.Manager, log *log.Logger) http.Handler {
	return &taskIDHandler{manager: manager, log: log}
}

func (h *taskIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Printf("Handling %s /api/v1/tasks/:id...", r.Method)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	taskID, err := parseLastPathSegment(r)
	if err != nil {
		h.log.Printf("Unable to parse last path segment: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	h.log.Printf("Getting taskID %d", taskID)

	t := h.manager.FindByID(taskID)
	if t == nil {
		h.log.Printf("No task with ID %d", taskID)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tJson, err := json.Marshal(t)
	if err != nil {
		h.log.Printf("Cannot marshal task: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.log.Printf("Returning task: %s", string(tJson))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(tJson)
	if err != nil {
		h.log.Printf("Cannot write JSON body: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
