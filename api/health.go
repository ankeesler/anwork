package api

import (
	"net/http"

	"github.com/ankeesler/anwork/lag"
)

type healthHandler struct {
	l *lag.L
}

func (h *healthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("all good"))
}
