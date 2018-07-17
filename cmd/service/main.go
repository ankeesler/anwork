package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/task/local"
)

func main() {
	var address string
	if port, ok := os.LookupEnv("PORT"); ok {
		address = fmt.Sprintf(":%s", port)
	} else {
		address = ":12345"
	}

	factory := local.NewManagerFactory("/tmp", "default-context")
	log := log.New(os.Stdout, "ANWORK Service: ", log.Ldate|log.Ltime|log.Lshortfile)
	api := api.New(address, factory, log)

	ctx, cancel := context.WithCancel(context.Background())
	if err := api.Run(ctx); err != nil {
		log.Fatalf("ERROR! api.Run() returned: %s", err.Error())
	}
	defer cancel()

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt)
	<-c
}
