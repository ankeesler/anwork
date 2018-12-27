// Package api provides an http.Handler that will serve the ANWORK API.
package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/ankeesler/anwork/lag"
	"github.com/ankeesler/anwork/task"
	"github.com/tedsuo/rata"
)

// This is the major verison of the API.
const Version = 1

//go:generate counterfeiter . Authenticator

// Authenticator is an object that performs authentication for the ANWORK API.
type Authenticator interface {
	// Authenticate performs auth on a token string. If it passes, it should
	// return a nil error. If it fails, it should return an error.
	Authenticate(token string) error
	// Generate a token for authentication. This token should probably be used
	// in the Authenticate method.
	Token() (string, error)
}

type api struct {
	l             *lag.L
	repo          task.Repo
	authenticator Authenticator
}

var routes = rata.Routes{
	{Name: "auth", Method: rata.POST, Path: "/api/v1/auth"},

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
func New(l *lag.L, repo task.Repo, authenticator Authenticator) http.Handler {
	return &api{
		l:             l,
		repo:          repo,
		authenticator: authenticator,
	}
}

func (a *api) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.l.P(lag.I, "api: handling %s %s", r.Method, r.URL.Path)

	if err, statusCode := a.authenticate(r); err != nil {
		respondWithError(a.l, w, statusCode, err)
		return
	}
	a.l.P(lag.I, "api: authentication succeeded")

	handlers := rata.Handlers{
		"auth": &authHandler{a.l, a.authenticator},

		"get_tasks":   &getTasksHandler{a.l, a.repo},
		"create_task": &createTaskHandler{a.l, a.repo},
		"get_task":    &getTaskHandler{a.l, a.repo},
		"update_task": &updateTaskHandler{a.l, a.repo},
		"delete_task": &deleteTaskHandler{a.l, a.repo},

		"get_events":   &getEventsHandler{a.l, a.repo},
		"create_event": &createEventHandler{a.l, a.repo},
		"get_event":    &getEventHandler{a.l, a.repo},
		"delete_event": &deleteEventHandler{a.l, a.repo},
	}
	router, err := rata.NewRouter(routes, handlers)
	if err != nil {
		respondWithError(a.l, w, http.StatusInternalServerError, err)
		return
	}

	router.ServeHTTP(w, r)
}

func (a *api) authenticate(r *http.Request) (error, int) {
	if r.URL.Path == "/api/v1/auth" {
		return nil, 0
	}

	tokenData := r.Header.Get("Authorization")
	if tokenData == "" {
		return errors.New("missing authorization header"), http.StatusUnauthorized
	}

	splitData := strings.Split(tokenData, " ")
	if len(splitData) != 2 || splitData[0] != "bearer" {
		return errors.New("invalid authorization data"), http.StatusBadRequest
	}

	return a.authenticator.Authenticate(splitData[1]), http.StatusForbidden
}

func respondWithError(l *lag.L, w http.ResponseWriter, statusCode int, err error) {
	respond(l, w, statusCode, Error{Message: err.Error()})
}

func respond(l *lag.L, w http.ResponseWriter, statusCode int, body interface{}) {
	l.P(lag.I, "api: responding with %d: %+v", statusCode, body)

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
