package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/task/local"
)

func main() {
	factory := local.NewManagerFactory("/tmp", "default-context")
	manager, err := factory.Create()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	log := log.New(os.Stdout, "ANWORK API: ", log.Ldate|log.Ltime|log.Lshortfile)
	a := api.New(manager, log)
	errChan := make(chan error)
	if err := a.Start(":12345", errChan); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}

	if err, ok := <-errChan; err != nil && ok {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
		os.Exit(1)
	}
}
