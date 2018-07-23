// This package contains an implementation of the anwork HTTP API.
//
// Here is the ANWORK API. All of the payloads are JSON formatted.
//   Show a list of links for this API:
//   GET /api -> return a map[string]map[string]string with links to the API endpoints
//
//   Get all of the current tasks:
//   GET  /api/v1/tasks -> returns an array of task.Task's
//
//   Create a new task:
//   POST /api/v1/tasks api.CreateRequest -> returns the created task.Task
//
//   Get the details about a task:
//   GET /api/v1/tasks/:id -> returns the task.Task
//
//   Update a task's state or priority:
//   PUT /api/v1/tasks/:id api.UpdateTaskRequest
//
//   Delete a task:
//   DELETE /api/v1/tasks/:id
//
//   Get all of the events:
//   GET /api/v1/events -> returns an array of task.Event's
//
//   Create a new event:
//   POST /api/v1/events api.AddEventRequest -> returns the created task.Event
//
//   Get the details about an event:
//   GET /api/v1/events/:startTime -> returns the task.Event that occurred at that time
//
//   Delete an event:
//   DELETE /api/v1/events/:startTime
package api

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/ankeesler/anwork/task"
)

type Api struct {
	address string
	factory task.ManagerFactory
	log     *log.Logger
}

func New(address string, factory task.ManagerFactory, log *log.Logger) *Api {
	return &Api{address: address, factory: factory, log: log}
}

// Run the API server. When the call returns, the server will be set up. The
// caller should use the provided context to determine when the server should be
// torn down.
func (a *Api) Run(ctx context.Context) error {
	a.log.Printf("API server starting on %s", a.address)

	server, err := a.makeServer()
	if err != nil {
		a.log.Printf("failed to make server: %s", err.Error())
		return err
	}

	listener, err := net.Listen("tcp", a.address)
	if err != nil {
		a.log.Printf("failed to listen on address %s: %s", a.address, err.Error())
		return err
	}

	go func() {
		<-ctx.Done()
		err = listener.Close()
		if err != nil {
			a.log.Printf("failed to close listener socket: %s", err.Error())
		}

		a.log.Printf("listener closed")
	}()

	go server.Serve(listener)

	return nil
}

func (a *Api) makeServer() (*http.Server, error) {
	manager, err := a.factory.Create()
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	mux.Handle("/", newFrontPageHandler(a.log))
	mux.Handle("/api", NewNavHandler(a.log))
	mux.Handle("/api/v1/health", NewHealthHandler(a.log))
	mux.Handle("/api/v1/tasks", NewTasksHandler(manager, a.log))
	mux.Handle("/api/v1/tasks/", NewTaskIDHandler(manager, a.log))
	mux.Handle("/api/v1/events", NewEventsHandler(manager, a.log))
	mux.Handle("/api/v1/events/", NewEventIDHandler(manager, a.log))

	return &http.Server{Handler: mux}, nil
}

func handleFrontPage(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
