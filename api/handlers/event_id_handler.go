package handlers

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

	t := h.manager.FindByID(eventID)
	if t == nil {
		h.log.Printf("No event with ID %d", eventID)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	tJson, err := json.Marshal(t)
	if err != nil {
		h.log.Printf("Cannot marshal event: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.log.Printf("Returning event: %s", string(tJson))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(tJson)
	if err != nil {
		h.log.Printf("Cannot write JSON body: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
