// Package api provides an http.Handler that will serve the ANWORK API.
package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ankeesler/anwork/task2"
)

//go:generate counterfeiter . Authenticator

// Authenticator is an object that performs authentication for the ANWORK API.
type Authenticator interface {
	// Authenticate performs auth on an http.Request. If it passes, it should
	// return a nil error. If it fails, it should return an error.
	Authenticate(req *http.Request) error
}

type api struct {
	log           *log.Logger
	repo          task2.Repo
	authenticator Authenticator
}

// New creates an http.Handler that will perform the ANWORK API functionality.
func New(log *log.Logger, repo task2.Repo, authenticator Authenticator) http.Handler {
	return &api{log: log, repo: repo, authenticator: authenticator}
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := a.authenticator.Authenticate(r); err != nil {
		a.respond(w, http.StatusForbidden, err)
		return
	}

	hData := findHData(r.URL.Path)
	if hData == nil {
		a.respond(w, http.StatusNotFound, nil)
		return
	}

	hDatum := findHDatum(hData, r.Method, "")
	if hDatum == nil {
		a.respond(w, http.StatusMethodNotAllowed, nil)
		return
	}

	rspStatus, rspBody, err := hDatum.handler(hDatum, a.repo)
	if err != nil {
		a.respond(w, rspStatus, err)
		return
	}

	data, err := json.Marshal(rspBody)
	if err != nil {
		a.respond(w, http.StatusInternalServerError, err)
		return
	}

	if n, err := w.Write(data); err != nil {
		a.respond(w, http.StatusInternalServerError, err)
		return
	} else if n != len(data) {
		err = fmt.Errorf("write underflow; got %d, expected %d", n, len(data))
		a.respond(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(rspStatus)
}

func (a *api) respond(w http.ResponseWriter, statusCode int, err error) {
	a.log.Printf("responding with %d", statusCode)
	w.WriteHeader(statusCode)
	if err != nil {
		a.log.Printf("err: %s", err.Error())
		bytes, jsonErr := json.Marshal(Error{Message: err.Error()})
		if jsonErr == nil {
			w.Write(bytes)
		}
	}
}
