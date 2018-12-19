package api

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/ankeesler/anwork/lag"
	"github.com/ankeesler/anwork/task"
	"github.com/tedsuo/rata"
)

func findEvent(l *lag.L, repo task.Repo, w http.ResponseWriter, r *http.Request) (*task.Event, int) {
	id := rata.Param(r, "id")
	idN, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(l, w, http.StatusBadRequest, err)
		return nil, 0
	}

	event, err := repo.FindEventByID(idN)
	if err != nil {
		respondWithError(l, w, http.StatusInternalServerError, err)
		return nil, 0
	}

	if event == nil {
		respondWithError(l, w, http.StatusNotFound, fmt.Errorf("unknown event with ID %d", idN))
		return nil, 0
	}

	return event, idN
}

type getEventHandler struct {
	l    *lag.L
	repo task.Repo
}

func (h *getEventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if event, _ := findEvent(h.l, h.repo, w, r); event != nil {
		respond(h.l, w, http.StatusOK, event)
	}
}

type deleteEventHandler struct {
	l    *lag.L
	repo task.Repo
}

func (h *deleteEventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if event, _ := findEvent(h.l, h.repo, w, r); event != nil {
		if err := h.repo.DeleteEvent(event); err != nil {
			respondWithError(h.l, w, http.StatusInternalServerError, err)
			return
		}

		respond(h.l, w, http.StatusNoContent, event)
	}
}
