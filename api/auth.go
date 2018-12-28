package api

import (
	"net/http"

	"code.cloudfoundry.org/lager"
)

type authHandler struct {
	logger        lager.Logger
	authenticator Authenticator
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, err := h.authenticator.Token()
	if err != nil {
		respondWithError(h.logger, w, http.StatusInternalServerError, err)
		return
	}

	auth := Auth{Token: token}
	respond(h.logger, w, http.StatusOK, auth)
}
