package api

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/ankeesler/anwork/task2"
	"github.com/tedsuo/rata"
)

func findEvent(log *log.Logger, repo task2.Repo, w http.ResponseWriter, r *http.Request) (*task2.Event, int) {
	id := rata.Param(r, "id")
	idN, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(log, w, http.StatusBadRequest, err)
		return nil, 0
	}

	event, err := repo.FindEventByID(idN)
	if err != nil {
		respondWithError(log, w, http.StatusInternalServerError, err)
		return nil, 0
	}

	if event == nil {
		respondWithError(log, w, http.StatusNotFound, fmt.Errorf("unknown event with ID %d", idN))
		return nil, 0
	}

	return event, idN
}

type getEventHandler struct {
	log  *log.Logger
	repo task2.Repo
}

func (h *getEventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if event, _ := findEvent(h.log, h.repo, w, r); event != nil {
		respond(h.log, w, http.StatusOK, event)
	}
}

type deleteEventHandler struct {
	log  *log.Logger
	repo task2.Repo
}

func (h *deleteEventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if event, _ := findEvent(h.log, h.repo, w, r); event != nil {
		if err := h.repo.DeleteEvent(event); err != nil {
			respondWithError(h.log, w, http.StatusInternalServerError, err)
			return
		}

		respond(h.log, w, http.StatusNoContent, event)
	}
}
