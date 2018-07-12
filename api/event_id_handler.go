package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ankeesler/anwork/task"
)

type eventIDHandler struct {
	manager task.Manager
	log     *log.Logger
}

func NewEventIDHandler(manager task.Manager, log *log.Logger) http.Handler {
	return &eventIDHandler{manager: manager, log: log}
}

func (h *eventIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Printf("Handling %s /api/v1/events/:id...", r.Method)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	eventID, err := parseLastPathSegment(r)
	if err != nil {
		msg := fmt.Sprintf("Unable to parse last path segment: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, msg, h.log)
		return
	}

	h.log.Printf("Getting eventID %d", eventID)

	var e *task.Event
	for _, event := range h.manager.Events() {
		if event.Date == int64(eventID) {
			e = event
			break
		}
	}

	if e == nil {
		msg := fmt.Sprintf("No event with ID %d", eventID)
		respondWithError(w, http.StatusNotFound, msg, h.log)
		return
	}

	eJson, err := json.Marshal(e)
	if err != nil {
		msg := fmt.Sprintf("Cannot marshal event: %s", err.Error())
		respondWithError(w, http.StatusInternalServerError, msg, h.log)
		return
	}

	h.log.Printf("Returning event: %s", string(eJson))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(eJson)
	if err != nil {
		msg := fmt.Sprintf("Cannot write JSON body: %s", err.Error())
		respondWithError(w, http.StatusInternalServerError, msg, h.log)
		return
	}
}
