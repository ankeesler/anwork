// Package api provides an http.Handler that will serve the ANWORK API.
package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ankeesler/anwork/task2"
	"github.com/tedsuo/rata"
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

var routes = rata.Routes{
	{Name: "get_tasks", Method: rata.GET, Path: "/api/v1/tasks"},
	{Name: "create_task", Method: rata.POST, Path: "/api/v1/tasks"},
	{Name: "get_task", Method: rata.GET, Path: "/api/v1/tasks/:id"},
	{Name: "update_task", Method: rata.PUT, Path: "/api/v1/tasks/:id"},
	{Name: "delete_task", Method: rata.DELETE, Path: "/api/v1/tasks/:id"},

	{Name: "get_events", Method: rata.GET, Path: "/api/v1/events"},
	{Name: "create_event", Method: rata.POST, Path: "/api/v1/events"},
	{Name: "get_event", Method: rata.GET, Path: "/api/v1/events/:id"},
	{Name: "delete_event", Method: rata.DELETE, Path: "/api/v1/events/:id"},
}

// New creates an http.Handler that will perform the ANWORK API functionality.
func New(log *log.Logger, repo task2.Repo, authenticator Authenticator) http.Handler {
	return &api{
		log:           log,
		repo:          repo,
		authenticator: authenticator,
	}
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.log.Printf("handling %s %s", r.Method, r.URL.Path)

	if err := a.authenticator.Authenticate(r); err != nil {
		respondWithError(a.log, w, http.StatusForbidden, err)
		return
	}
	a.log.Printf("authentication succeeded")

	handlers := rata.Handlers{
		"get_tasks":   &getTasksHandler{a.log, a.repo},
		"create_task": &createTaskHandler{a.log, a.repo},
		"get_task":    &getTaskHandler{a.log, a.repo},
		"update_task": &updateTaskHandler{a.log, a.repo},
		"delete_task": &deleteTaskHandler{a.log, a.repo},

		"get_events":   &getEventsHandler{a.log, a.repo},
		"create_event": &createEventHandler{a.log, a.repo},
		"get_event":    &getEventHandler{a.log, a.repo},
		"delete_event": &deleteEventHandler{a.log, a.repo},
	}
	router, err := rata.NewRouter(routes, handlers)
	if err != nil {
		respondWithError(a.log, w, http.StatusInternalServerError, err)
		return
	}

	router.ServeHTTP(w, r)
}

func respondWithError(log *log.Logger, w http.ResponseWriter, statusCode int, err error) {
	respond(log, w, statusCode, Error{Message: err.Error()})
}

func respond(log *log.Logger, w http.ResponseWriter, statusCode int, body interface{}) {
	log.Printf("responding with %d: %+v", statusCode, body)

	var bytes []byte = []byte{}
	var jsonErr error
	if body != nil {
		bytes, jsonErr = json.Marshal(body)
	}

	if jsonErr != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(jsonErr.Error()))
	} else {
		w.WriteHeader(statusCode)
		w.Write(bytes)
	}
}
