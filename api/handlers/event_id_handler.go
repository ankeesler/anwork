package handlers

import (
	"log"
	"net/http"
)

type eventIDHandler struct {
	log *log.Logger
}

func (h *eventIDHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.log.Printf("Handling /api/v1/events/:id...")

	eventID, err := parseLastPathSegment(r)
	if err != nil {
		h.log.Printf("...ERROR! Cannot parse eventID from end of URL %s: %s",
			r.URL.String(),
			err.Error())
		return
	}
	h.log.Printf("...eventID=%d", eventID)
}
