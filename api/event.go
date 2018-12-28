package api

import (
	"fmt"
	"net/http"
	"strconv"

	"code.cloudfoundry.org/lager"
	"github.com/ankeesler/anwork/task"
	"github.com/tedsuo/rata"
)

func findEvent(
	logger lager.Logger,
	repo task.Repo,
	w http.ResponseWriter,
	r *http.Request,
) (*task.Event, int) {
	id := rata.Param(r, "id")
	idN, err := strconv.Atoi(id)
	if err != nil {
		respondWithError(logger, w, http.StatusBadRequest, err)
		return nil, 0
	}

	event, err := repo.FindEventByID(idN)
	if err != nil {
		respondWithError(logger, w, http.StatusInternalServerError, err)
		return nil, 0
	}

	if event == nil {
		respondWithError(logger, w, http.StatusNotFound, fmt.Errorf("unknown event with ID %d", idN))
		return nil, 0
	}

	return event, idN
}

type getEventHandler struct {
	logger lager.Logger
	repo   task.Repo
}

func (h *getEventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if event, _ := findEvent(h.logger, h.repo, w, r); event != nil {
		respond(h.logger, w, http.StatusOK, event)
	}
}

type deleteEventHandler struct {
	logger lager.Logger
	repo   task.Repo
}

func (h *deleteEventHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if event, _ := findEvent(h.logger, h.repo, w, r); event != nil {
		if err := h.repo.DeleteEvent(event); err != nil {
			respondWithError(h.logger, w, http.StatusInternalServerError, err)
			return
		}

		respond(h.logger, w, http.StatusNoContent, event)
	}
}
