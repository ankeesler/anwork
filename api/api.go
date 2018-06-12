// This package contains an implementation of the anwork HTTP API.
//
// TODO: document me!
package api

import (
	"log"

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

// Run the API server. This calls http.ListenAndServe so it will block.
func (a *Api) Run() error {
	a.log.Printf("API server starting on %s", a.address)
	//return http.ListenAndServe(a.address, nil)
	return nil
}
