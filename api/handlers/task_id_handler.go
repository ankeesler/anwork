package handlers

import (
	"log"
	"net/http"
)

type taskIDHandler struct {
	log *log.Logger
}

func (h *taskIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Printf("Handling /api/v1/tasks/:id...")

	taskID, err := parseLastPathSegment(r)
	if err != nil {
		h.log.Printf("...ERROR! Cannot parse taskID from end of URL %s: %s",
			r.URL.String(),
			err.Error())
		return
	}
	h.log.Printf("...taskID=%d", taskID)
}
