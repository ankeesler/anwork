// This package contains an implementation of the anwork HTTP API.
//
// TODO: document me!
package api

import (
	"context"
	"log"
	"net"
	"net/http"

	"github.com/ankeesler/anwork/api/handlers"
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
	mux.Handle("/api/v1/tasks", handlers.NewTasksHandler(manager, a.log))

	return &http.Server{Handler: mux}, nil
}
