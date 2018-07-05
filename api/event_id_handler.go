package api

import (
	"encoding/json"
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
		h.log.Printf("Unable to parse last path segment: %s", err.Error())
		w.WriteHeader(http.StatusBadRequest)
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
		h.log.Printf("No event with ID %d", eventID)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	eJson, err := json.Marshal(e)
	if err != nil {
		h.log.Printf("Cannot marshal event: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.log.Printf("Returning event: %s", string(eJson))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(eJson)
	if err != nil {
		h.log.Printf("Cannot write JSON body: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
