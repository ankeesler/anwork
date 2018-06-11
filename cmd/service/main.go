package main

import (
	"log"
	"os"

	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/task/local"
)

func main() {
	address := ":12345"
	factory := local.NewManagerFactory("/tmp", "default-context")
	log := log.New(os.Stdout, "ANWORK Service: ", log.Ldate|log.Ltime|log.Llongfile)
	api := api.New(address, factory, log)
	if err := api.Run(); err != nil {
		log.Fatalf("ERROR! api.Run() returned: %s", err.Error())
	} else {
		log.Fatalf("ERROR! api.Run() returned: nil")
	}
}
