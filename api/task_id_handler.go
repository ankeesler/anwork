package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	taskID, err := parseLastPathSegment(r)
	if err != nil {
		respondWithError2(w, http.StatusInternalServerError, "Unable to parse last path segment", h.log)
		return
	}
	h.log.Printf("Getting taskID %d", taskID)

	t := h.manager.FindByID(taskID)
	if t == nil {
		respondWithError2(w, http.StatusNotFound, fmt.Sprintf("No task with ID %d", taskID), h.log)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r, t)
	case http.MethodPut:
		h.handlePut(w, r, t)
	case http.MethodDelete:
		h.handleDelete(w, r, t)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *taskIDHandler) handleGet(w http.ResponseWriter, r *http.Request, t *task.Task) {
	tJson, err := json.Marshal(t)
	if err != nil {
		msg := fmt.Sprintf("Cannot marshal task: %s", err.Error())
		respondWithError2(w, http.StatusInternalServerError, msg, h.log)
		return
	}

	h.log.Printf("Returning task: %s", string(tJson))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(tJson)
	if err != nil {
		msg := fmt.Sprintf("Cannot write JSON body: %s", err.Error())
		respondWithError2(w, http.StatusInternalServerError, msg, h.log)
		return
	}
}

func (h *taskIDHandler) handlePut(w http.ResponseWriter, r *http.Request, t *task.Task) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}
	h.log.Printf("handling request %s", string(body))

	var req UpdateTaskRequest
	if err := json.Unmarshal(body, &req); err != nil {
		msg := fmt.Sprintf("Cannot unmarshal request body: %s", err.Error())
		respondWithError2(w, http.StatusBadRequest, msg, h.log)
		return
	}

	if req.State != 0 {
		if int(req.State) < 0 || int(req.State) >= len(task.StateNames) {
			respondWithError2(w, http.StatusBadRequest, fmt.Sprintf("invalid state %d", req.State), h.log)
			return
		}

		if err := h.manager.SetState(t.Name, req.State); err != nil {
			msg := fmt.Sprintf("Failed to set state: %s", err.Error())
			respondWithError2(w, http.StatusInternalServerError, msg, h.log)
			return
		}
		h.log.Printf("set state %s", task.StateNames[req.State])
	}
	if req.Priority != 0 {
		if err := h.manager.SetPriority(t.Name, req.Priority); err != nil {
			msg := fmt.Sprintf("Failed to set priority: %s", err.Error())
			respondWithError2(w, http.StatusInternalServerError, msg, h.log)
			return
		}
		h.log.Printf("set priority %d", req.Priority)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *taskIDHandler) handleDelete(w http.ResponseWriter, r *http.Request, t *task.Task) {
	if err := h.manager.Delete(t.Name); err != nil {
		msg := fmt.Sprintf("Unable to delete task %s: %s", t.Name, err.Error())
		respondWithError2(w, http.StatusInternalServerError, msg, h.log)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
