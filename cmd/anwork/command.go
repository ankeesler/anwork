package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/ankeesler/anwork/task"
)

// A command is a keyword (see name field) passed to the anwork executable that provokes some
// functionality (see action field).
type command struct {
	name, description string

	// These are the name of the arguments that are expected by the command.
	args []string

	// This is the action that runs when this command is invoked. The first parameter to this function
	// is the name of the command. The second parameter to the function is the output stream to which
	// output should be written. The function should return true iff the task.Manager should be
	// persisted to disk after the action returns.
	action func(name string, output io.Writer, manager *task.Manager) bool
}

// These are the commands used by the anwork application.
var commands = []command{
	command{
		name:        "version",
		description: "Print version information",
		args:        []string{},
		action:      versionAction,
	},
	command{
		name:        "create",
		description: "Create a new task",
		args:        []string{"task-name"},
		action:      createAction,
	},
	command{
		name:        "delete",
		description: "Delete a task",
		args:        []string{"task-name"},
		action:      deleteAction,
	},
	command{
		name:        "show",
		description: "Show the current tasks",
		args:        []string{},
		action:      showAction,
	},
	command{
		name:        "note",
		description: "Add a note to a task",
		args:        []string{"task-name", "note"},
		action:      noteAction,
	},
	command{
		name:        "set-running",
		description: "Mark a task as running",
		args:        []string{"task-name"},
		action:      setStateAction,
	},
	command{
		name:        "set-blocked",
		description: "Mark a task as blocked",
		args:        []string{"task-name"},
		action:      setStateAction,
	},
	command{
		name:        "set-waiting",
		description: "Mark a task as waiting",
		args:        []string{"task-name"},
		action:      setStateAction,
	},
	command{
		name:        "set-finished",
		description: "Mark a task as finished",
		args:        []string{"task-name"},
		action:      setStateAction,
	},
	command{
		name:        "journal",
		description: "Show the journal",
		args:        []string{},
		action:      journalAction,
	},
}

func findCommand(name string) *command {
	for _, c := range commands {
		if c.name == name {
			return &c
		}
	}
	return nil
}

// This index starts at "1" since the arg at index "0" will be the command name itself.
var shiftIndex int = 1

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

func versionAction(command string, output io.Writer, manager *task.Manager) bool {
	fmt.Fprintln(output, "ANWORK Version =", version)
	return false
}

func createAction(command string, output io.Writer, manager *task.Manager) bool {
	dbgfln(os.Stdout, "Creating task...")
	name := shift()
	manager.Create(name)
	return true
}

func showAction(command string, output io.Writer, manager *task.Manager) bool {
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

func noteAction(command string, output io.Writer, manager *task.Manager) bool {
	name := shift()
	note := shift()
	manager.Note(name, note)
	return true
}

func deleteAction(command string, output io.Writer, manager *task.Manager) bool {
	dbgfln(os.Stdout, "Deleting task...")
	name := shift()
	if !manager.Delete(name) {
		fmt.Printf("Error! Unknown task: %s\n", name)
		return false
	} else {
		return true
	}
}

func setStateAction(command string, output io.Writer, manager *task.Manager) bool {
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

func journalAction(command string, output io.Writer, manager *task.Manager) bool {
	es := manager.Journal().Events
	for i := len(es) - 1; i >= 0; i-- {
		e := es[i]
		fmt.Printf("[%s %s %d %02d:%02d]: %s\n", e.Date.Weekday(), e.Date.Month(), e.Date.Day(),
			e.Date.Hour(), e.Date.Minute(), e.Title)
	}
	return false
}
