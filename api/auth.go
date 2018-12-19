package api

import (
	"net/http"

	"github.com/ankeesler/anwork/lag"
)

type authHandler struct {
	l             *lag.L
	authenticator Authenticator
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, err := h.authenticator.Token()
	if err != nil {
		respondWithError(h.l, w, http.StatusInternalServerError, err)
		return
	}

	auth := Auth{Token: token}
	respond(h.l, w, http.StatusOK, auth)
}
