package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/ankeesler/anwork/lag"
	"github.com/ankeesler/anwork/task"
)

type getEventsHandler struct {
	l    *lag.L
	repo task.Repo
}

func (h *getEventsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	events, err := h.repo.Events()
	if err != nil {
		respondWithError(h.l, w, http.StatusInternalServerError, err)
		return
	}

	respond(h.l, w, http.StatusOK, events)
}

type createEventHandler struct {
	l    *lag.L
	repo task.Repo
}

func (h *createEventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(h.l, w, http.StatusInternalServerError, err)
		return
	}

	var event task.Event
	if err := json.Unmarshal(data, &event); err != nil {
		respondWithError(h.l, w, http.StatusBadRequest, err)
		return
	}

	if err := h.repo.CreateEvent(&event); err != nil {
		respondWithError(h.l, w, http.StatusInternalServerError, err)
		return
	}

	w.Header().Add("Location", fmt.Sprintf("/api/v1/events/%d", event.ID))
	respond(h.l, w, http.StatusCreated, nil)
}
