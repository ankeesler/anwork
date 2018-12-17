package main

import (
	"fmt"
	"log"
	"os"

	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/api/authenticator"
	"github.com/ankeesler/anwork/task/fs"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/http_server"
)

func main() {
	var address string
	if port, ok := os.LookupEnv("PORT"); ok {
		address = fmt.Sprintf(":%s", port)
	} else {
		address = ":12345"
	}

	log := log.New(os.Stdout, "ANWORK Service: ", log.Ldate|log.Ltime|log.Lshortfile)
	repo := fs.New("/tmp/default-context")
	authenticator := authenticator.New()

	log.Printf("hey")

	runner := http_server.New(address, api.New(log, repo, authenticator))
	process := ifrit.Invoke(runner)
	log.Printf("running")

	log.Fatal(<-process.Wait())
}
