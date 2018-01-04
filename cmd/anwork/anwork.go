package main

import (
	"flag"
	"fmt"
	"strings"

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

// This map stores the functions associated with each of the command line actions.
var commands = []command{
	command{name: "version", usage: "(no arguments)", action: versionAction},
	command{name: "create", usage: "<task-name>", action: createAction},
	command{name: "delete", usage: "<task-name>", action: deleteAction},
	command{name: "show", usage: "(no arguments)", action: showAction},
	command{name: "note", usage: "<task-name> <note>", action: noteAction},
	command{name: "set-running", usage: "<task-name>", action: setStateAction},
	command{name: "set-blocked", usage: "<task-name>", action: setStateAction},
	command{name: "set-waiting", usage: "<task-name>", action: setStateAction},
	command{name: "set-finished", usage: "<task-name>", action: setStateAction},
	command{name: "journal", usage: "(no arguments)", action: journalAction},
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

func versionAction(command string, manager *task.Manager) bool {
	fmt.Println("ANWORK Version =", version)
	return false
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

func noteAction(command string, manager *task.Manager) bool {
	name := shift()
	note := shift()
	manager.Note(name, note)
	return true
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

	firstArg := shift()
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
