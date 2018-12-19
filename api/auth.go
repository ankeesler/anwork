package api

import (
	"log"
	"net/http"
)

type authHandler struct {
	log           *log.Logger
	authenticator Authenticator
}

func (a *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	token, err := a.authenticator.Token()
	if err != nil {
		respondWithError(a.log, w, http.StatusInternalServerError, err)
		return
	}

	auth := Auth{Token: token}
	respond(a.log, w, http.StatusOK, auth)
}
