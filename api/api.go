// This package contains an implementation of the anwork HTTP API.
//
// TODO: document me!
package api

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strconv"
	"strings"

	"github.com/ankeesler/anwork/task"
)

type Api struct {
	manager task.Manager

	httpServer http.Server
	log        *log.Logger
}

func New(manager task.Manager, l *log.Logger) *Api {
	return &Api{manager: manager, log: l}
}

// Start the API server. This function will start the HTTP server in a goroutine
// and return. The caller should use the errorChan parameter to detect if the server
// falls over for some reason.
//
// Here is a sample use.
//   a := api.New(...)
//   errorChan := make(chan error)
//   if err := a.Start(":12345", errorChan); err != nil {
//     // Handle error...
//   }
//   // Start doing stuff...
//   err := <-errorChan
//   if errIsReallyBad(err) {
//     if stopErr := api.Stop(); stopErr != nil {
//       // Handle stop error...
//     }
//     // Handle really bad error...
//   } else {
//     // Do some custom error handling and let the server keep running...
//   }
func (a *Api) Start(address string, errChan chan error) error {
	a.log.Printf("API server starting on %s", address)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		myErr := fmt.Errorf("Cannot start server: %s", err.Error())
		a.log.Printf("ERROR! %s", err.Error())
		return myErr
	}

	mux := http.NewServeMux()
	mux.Handle("/api/v1/health", &healthHandler{log: a.log})
	mux.Handle("/api/v1/tasks", &tasksHandler{log: a.log})
	mux.Handle("/api/v1/tasks/", &taskIDHandler{log: a.log})
	mux.Handle("/api/v1/events", &eventsHandler{log: a.log})
	mux.Handle("/api/v1/events/", &eventIDHandler{log: a.log})

	a.httpServer.Addr = address
	a.httpServer.Handler = mux

	go func() {
		errChan <- a.httpServer.Serve(listener)
	}()

	return nil
}

// Stop running the API server. Returns an error if we fail to do so.
func (a *Api) Stop() error {
	return a.httpServer.Close()
}

func parseLastPathSegment(r *http.Request) (int, error) {
	segs := strings.Split(r.URL.EscapedPath(), "/")
	lastSeg := segs[len(segs)-1]
	return strconv.Atoi(lastSeg)
}
