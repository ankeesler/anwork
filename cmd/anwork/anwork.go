package main

import (
	"flag"
	"fmt"

	"github.com/ankeesler/anwork/storage"
	"github.com/ankeesler/anwork/task"
)

// These variables are used to store command line flag values.
var (
	help, debug   bool
	context, root string
)

// This map stores the functions associated with each of the command line actions.
var argActions = map[string]func(*task.Manager) bool{
	"create": createAction,
}

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

func createAction(manager *task.Manager) bool {
	dbgfln("Creating task...")
	return true
}

func main() {
	flag.BoolVar(&help, "help", false, "Print this help message")
	flag.BoolVar(&debug, "debug", false, "Enable debug printing")

	flag.StringVar(&context, "context", "default-context", "Set the persistence context")
	flag.StringVar(&root, "root", ".", "Set the persistence root directory")

	flag.Parse()
	if help {
		flag.Usage()
		return
	}

	if flag.NArg() == 0 {
		fmt.Println("Error! Expected command arguments")
		flag.Usage()
		return
	}

	persister := storage.Persister{root}
	manager := task.NewManager()
	if persister.Exists(context) {
		dbgfln("Context %s exists at root %s", context, root)
		readManager(&persister, context, manager)
	} else {
		dbgfln("Context %s does not exist at root %s; creating it", context, root)
		writeManager(&persister, context, manager)
	}
	dbgfln("Manager is %#v", manager)

	action := argActions[flag.Arg(0)]
	if action == nil {
		fmt.Println("Error! Unknown command line argument:", flag.Arg(0))
		return
	} else {
		if action(manager) {
			dbgfln("Persisting manager back to disk")
			writeManager(&persister, context, manager)
		} else {
			dbgfln("NOT persisting manager back to disk")
		}
	}
}
