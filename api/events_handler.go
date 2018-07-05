package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ankeesler/anwork/task"
)

type eventsHandler struct {
	manager task.Manager
	log     *log.Logger
}

func NewEventsHandler(manager task.Manager, log *log.Logger) http.Handler {
	return &eventsHandler{manager: manager, log: log}
}

func (h *eventsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Printf("Handling %s /api/v1/events...", r.Method)

	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	events := h.manager.Events()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	eventsJson, err := json.Marshal(events)
	if err != nil {
		h.log.Printf("Failed to marshal events: %s", err.Error())
		return
	}

	h.log.Printf("Returning events %s", eventsJson)
	_, err = w.Write(eventsJson)
	if err != nil {
		h.log.Printf("Cannot write JSON body: %s", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
