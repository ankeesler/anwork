package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ankeesler/anwork/task"
)

type getEventsHandler struct {
	log  *log.Logger
	repo task.Repo
}

func (h *getEventsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	events, err := h.repo.Events()
	if err != nil {
		respondWithError(h.log, w, http.StatusInternalServerError, err)
		return
	}

	respond(h.log, w, http.StatusOK, events)
}

type createEventHandler struct {
	log  *log.Logger
	repo task.Repo
}

func (h *createEventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(h.log, w, http.StatusInternalServerError, err)
		return
	}

	var event task.Event
	if err := json.Unmarshal(data, &event); err != nil {
		respondWithError(h.log, w, http.StatusBadRequest, err)
		return
	}

	if err := h.repo.CreateEvent(&event); err != nil {
		respondWithError(h.log, w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Location", fmt.Sprintf("/api/v1/events/%d", event.ID))
	respond(h.log, w, http.StatusCreated, nil)
}
