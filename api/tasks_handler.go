package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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
	payload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errMsg := fmt.Sprintf("Could not read request body: %s", err.Error())
		h.log.Printf(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		respondWithError(w, errMsg)
		return
	}

	var createReq CreateRequest
	if err := json.Unmarshal(payload, &createReq); err != nil {
		errMsg := fmt.Sprintf("Cannot unmarshal payload '%s': %s", string(payload), err.Error())
		h.log.Printf(errMsg)
		w.WriteHeader(http.StatusBadRequest)
		respondWithError(w, errMsg)
		return
	}
	h.log.Printf("Decoded create task request: %+v", createReq)

	if err := h.manager.Create(createReq.Name); err != nil {
		errMsg := fmt.Sprintf("Cannot create task: %s", err.Error())
		h.log.Printf(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		respondWithError(w, errMsg)
		return
	}

	t := h.manager.FindByName(createReq.Name)
	if t == nil {
		errMsg := fmt.Sprintf("Cannot find newly created task: %s", err.Error())
		h.log.Printf(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		respondWithError(w, errMsg)
		return
	}
	h.log.Printf("Created task %s", t.Name)

	tJson, err := json.Marshal(t)
	if err != nil {
		errMsg := fmt.Sprintf("Cannot marshal respond task: %s", err.Error())
		h.log.Printf(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		respondWithError(w, errMsg)
		return
	}
	h.log.Printf("Responding with new task %s", tJson)

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Location", fmt.Sprintf("/api/v1/tasks/%d", t.ID))
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(tJson); err != nil {
		errMsg := fmt.Sprintf("Cannot write response json: %s", err.Error())
		h.log.Printf(errMsg)
		w.WriteHeader(http.StatusInternalServerError)
		respondWithError(w, errMsg)
		return
	}
}
