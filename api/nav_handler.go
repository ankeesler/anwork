package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type navHandler struct {
	log *log.Logger
}

func NewNavHandler(log *log.Logger) http.Handler {
	return &navHandler{log: log}
}

func (nh *navHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	nh.log.Printf("Handling %s /api...", req.Method)

	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	nav := make(map[string]map[string]string)
	nav["links"] = make(map[string]string)
	links := nav["links"]

	links["tasks"] = "/api/v1/tasks"
	links["events"] = "/api/v1/events"
	links["health"] = "/api/v1/health"

	navJSON, err := json.Marshal(nav)
	if err != nil {
		msg := fmt.Sprintf("Failed to marshal nav %s: %s", nav, err.Error())
		respondWithError(w, http.StatusInternalServerError, msg, nh.log)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(navJSON); err != nil {
		msg := fmt.Sprintf("Failed to write nav JSON %s: %s", navJSON, err.Error())
		respondWithError(w, http.StatusInternalServerError, msg, nh.log)
	}
}
