// This package contains an implementation of the anwork HTTP API.
//
// TODO: document me!
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

	listener, err := net.Listen("tcp", a.address)
	if err != nil {
		return err
	}

	go func() {
		err := http.Serve(listener, nil)
		if err != nil {
			a.log.Printf("API server exited with error: %s", err.Error())
		} else {
			a.log.Printf("API server exited successfully")
		}
	}()

	go func() {
		<-ctx.Done()
		err := listener.Close()
		if err != nil {
			a.log.Printf("API server failed to close listener socket: %s", err.Error())
		} else {
			a.log.Printf("API server successfully closed listener socket")
		}
	}()

	return nil
}
