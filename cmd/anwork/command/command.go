// This package contains the commands that can be passed to the anwork command line tool.
package command

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/ankeesler/anwork/task"
)

// This is the version of this anwork application command set.
const Version = 1

// A Command is a keyword (see name field) passed to the anwork executable that provokes some
// functionality (see action field).
type Command struct {
	Name, Description string

	// This slice holds the name(s) of the argument(s) that is(/are) expected by the Command.
	Args []string

	// This is the action that runs when this Command is invoked. The first parameter to this function
	// is the name of the Command. The second parameter to the function is the output stream to which
	// output should be written. The function should return true iff the task.Manager should be
	// persisted to disk after the action returns.
	Action func(name string, output io.Writer, manager *task.Manager) bool
}

// These are the Command's used by the anwork application.
var Commands = []Command{
	Command{
		Name:        "version",
		Description: "Print version information",
		Args:        []string{},
		Action:      versionAction,
	},
	Command{
		Name:        "create",
		Description: "Create a new task",
		Args:        []string{"task-name"},
		Action:      createAction,
	},
	Command{
		Name:        "delete",
		Description: "Delete a task",
		Args:        []string{"task-name"},
		Action:      deleteAction,
	},
	Command{
		Name:        "show",
		Description: "Show the current tasks",
		Args:        []string{},
		Action:      showAction,
	},
	Command{
		Name:        "note",
		Description: "Add a note to a task",
		Args:        []string{"task-name", "note"},
		Action:      noteAction,
	},
	Command{
		Name:        "set-running",
		Description: "Mark a task as running",
		Args:        []string{"task-name"},
		Action:      setStateAction,
	},
	Command{
		Name:        "set-blocked",
		Description: "Mark a task as blocked",
		Args:        []string{"task-name"},
		Action:      setStateAction,
	},
	Command{
		Name:        "set-waiting",
		Description: "Mark a task as waiting",
		Args:        []string{"task-name"},
		Action:      setStateAction,
	},
	Command{
		Name:        "set-finished",
		Description: "Mark a task as finished",
		Args:        []string{"task-name"},
		Action:      setStateAction,
	},
	Command{
		Name:        "journal",
		Description: "Show the journal",
		Args:        []string{},
		Action:      journalAction,
	},
}

// Print the usage information out for a single Command. The information will be printed to the
// provided output stream.
func (c *Command) Usage(output io.Writer) {
	fmt.Fprintf(output, "  %s %s\n", c.Name, strings.Join(c.Args, " "))
	fmt.Fprintf(output, "        %s\n", c.Description)
}

// Find the command with the provided name.
func FindCommand(name string) *Command {
	for _, c := range Commands {
		if c.Name == name {
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
	fmt.Fprintln(output, "ANWORK Version =", Version)
	return false
}

func createAction(command string, output io.Writer, manager *task.Manager) bool {
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
