package api

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

	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodDelete:
		h.handleDelete(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *taskIDHandler) handleGet(w http.ResponseWriter, r *http.Request) {
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

func (h *taskIDHandler) handleDelete(w http.ResponseWriter, r *http.Request) {
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

	err = h.manager.Delete(t.Name)
	if err != nil {
		h.log.Printf("Unable to delete task %s: %s", t.Name, err.Error())
		respondWithError(w, err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
