package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	switch r.Method {
	case http.MethodGet:
		h.handleGet(w, r)
	case http.MethodPost:
		h.handlePost(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (h *eventsHandler) handleGet(w http.ResponseWriter, r *http.Request) {
	events := h.manager.Events()
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	eventsJson, err := json.Marshal(events)
	if err != nil {
		msg := fmt.Sprintf("Failed to marshal events: %s", err.Error())
		respondWithError(w, http.StatusInternalServerError, msg, h.log)
		return
	}

	h.log.Printf("Returning events %s", eventsJson)
	_, err = w.Write(eventsJson)
	if err != nil {
		msg := fmt.Sprintf("Cannot write JSON body: %s", err.Error())
		respondWithError(w, http.StatusInternalServerError, msg, h.log)
		return
	}
}

func (h *eventsHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	reqPayload, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error(), h.log)
	}

	var addEventReq AddEventRequest
	if err := json.Unmarshal(reqPayload, &addEventReq); err != nil {
		msg := fmt.Sprintf("Invalid request payload: %s", string(reqPayload))
		respondWithError(w, http.StatusBadRequest, msg, h.log)
		return
	}

	if addEventReq.Type != task.EventTypeNote {
		msg := fmt.Sprintf("Invalid event type %d, the only supported event type is %d",
			addEventReq.Type, task.EventTypeNote)
		respondWithError(w, http.StatusBadRequest, msg, h.log)
		return
	}

	t := h.manager.FindByID(addEventReq.TaskID)
	if t == nil {
		msg := fmt.Sprintf("Unknown task for ID %d", addEventReq.TaskID)
		respondWithError(w, http.StatusBadRequest, msg, h.log)
		return
	}

	if err := h.manager.Note(t.Name, addEventReq.Title); err != nil {
		msg := fmt.Sprintf("Failed to add note: %s", err.Error())
		respondWithError(w, http.StatusBadRequest, msg, h.log)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
