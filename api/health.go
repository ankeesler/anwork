package api

import (
	"net/http"
)

type healthHandler struct{}

func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("all good"))
}
