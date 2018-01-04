package main

import (
	"flag"
	"fmt"

	"github.com/ankeesler/anwork/storage"
	"github.com/ankeesler/anwork/task"
)

// This is the version of this anwork app.
const version = 1

// These variables are used to store command line flag values.
var (
	help, debug   bool
	context, root string
)

func dbgfln(format string, stuff ...interface{}) {
	if debug {
		fmt.Print("anwork: dbg: ")
		fmt.Printf(format, stuff...)
		fmt.Println()
	}
}

func readManager(persister *storage.Persister, context string, manager *task.Manager) {
	err := persister.Unpersist(context, manager)
	if err != nil {
		panic(err.Error())
	}
}

func writeManager(persister *storage.Persister, context string, manager *task.Manager) {
	err := persister.Persist(context, manager)
	if err != nil {
		panic(err.Error())
	}
}

func findAction(name string) func(string, *task.Manager) bool {
	for _, c := range commands {
		if c.name == name {
			return c.action
		}
	}
	return nil
}

func main() {
	flag.BoolVar(&help, "help", false, "Print this help message")
	flag.BoolVar(&debug, "debug", false, "Enable debug printing")

	flag.StringVar(&context, "context", "default-context", "Set the persistence context")
	flag.StringVar(&root, "root", ".", "Set the persistence root directory")

	flag.Parse()
	if help {
		// TODO: write our own usage with the commands!
		flag.Usage()
		return
	}

	if flag.NArg() == 0 {
		fmt.Println("Error! Expected command arguments")
		flag.Usage()
		return
	}
	firstArg := flag.Arg(0)

	persister := storage.Persister{root}
	manager := task.NewManager()
	if persister.Exists(context) {
		dbgfln("Context %s exists at root %s", context, root)
		readManager(&persister, context, manager)
	} else {
		dbgfln("Context %s does not exist at root %s; creating it", context, root)
		writeManager(&persister, context, manager)
	}
	dbgfln("Manager is %s", manager)

	action := findAction(firstArg)
	if action == nil {
		fmt.Println("Error! Unknown command line argument:", firstArg)
		return
	} else {
		if action(firstArg, manager) {
			dbgfln("Persisting manager back to disk")
			writeManager(&persister, context, manager)
		} else {
			dbgfln("NOT persisting manager back to disk")
		}
	}
}
