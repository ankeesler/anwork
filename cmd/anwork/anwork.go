package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/ankeesler/anwork/storage"
	"github.com/ankeesler/anwork/task"
)

// These variables are used to store command line flag values.
var (
	help, debug   bool
	context, root string
)

// This map stores the functions associated with each of the command line actions.
var argActions = map[string]func(string, *task.Manager) bool{
	"create":       createAction,
	"show":         showAction,
	"delete":       deleteAction,
	"set-running":  setStateAction,
	"set-blocked":  setStateAction,
	"set-waiting":  setStateAction,
	"set-finished": setStateAction,
	"journal":      journalAction,
}

var shiftIndex int = 0

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

// This function gets the next command line argument (via flag.Args), or panics if there is no such
// argument.
func shift() string {
	if shiftIndex >= flag.NArg() {
		msg := fmt.Sprintf("Failed to retrieve arg %d", shiftIndex)
		panic(msg)
	}
	val := flag.Arg(shiftIndex)
	shiftIndex++
	return val
}

func createAction(command string, manager *task.Manager) bool {
	dbgfln("Creating task...")
	name := shift()
	manager.Create(name)
	return true
}

func showAction(command string, manager *task.Manager) bool {
	printer := func(state task.State) {
		fmt.Printf("%s tasks:\n", strings.ToUpper(task.StateNames[state]))
		for _, task := range manager.Tasks() {
			if task.State() == state {
				fmt.Printf("  %s (%d)\n", task.Name(), task.ID())
			}
		}
	}
	printer(task.StateRunning)
	printer(task.StateBlocked)
	printer(task.StateWaiting)
	printer(task.StateFinished)
	return false
}

func deleteAction(command string, manager *task.Manager) bool {
	dbgfln("Deleting task...")
	name := shift()
	if !manager.Delete(name) {
		fmt.Printf("Error! Unknown task: %s\n", name)
		return false
	} else {
		return true
	}
}

func setStateAction(command string, manager *task.Manager) bool {
	name := shift()

	var state task.State
	command = strings.TrimPrefix(command, "set-")
	switch command {
	case "running":
		state = task.StateRunning
	case "blocked":
		state = task.StateBlocked
	case "waiting":
		state = task.StateWaiting
	case "finished":
		state = task.StateFinished
	default:
		panic("Unknown state: " + command)
	}
	manager.SetState(name, state)
	return true
}

func journalAction(command string, manager *task.Manager) bool {
	es := manager.Journal().Events
	for i := len(es) - 1; i >= 0; i-- {
		e := es[i]
		fmt.Printf("[%s %s %d %02d:%02d]: %s\n", e.Date.Weekday(), e.Date.Month(), e.Date.Day(),
			e.Date.Hour(), e.Date.Minute(), e.Title)
	}
	return false
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
	dbgfln("Manager is %s", manager)

	command := shift()
	action := argActions[command]
	if action == nil {
		fmt.Println("Error! Unknown command line argument:", command)
		return
	} else {
		if action(command, manager) {
			dbgfln("Persisting manager back to disk")
			writeManager(&persister, context, manager)
		} else {
			dbgfln("NOT persisting manager back to disk")
		}
	}
}
