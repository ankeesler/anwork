// This package contains the commands that can be passed to the anwork command line tool.
package command

import (
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/ankeesler/anwork/task"
)

// This is the version of this anwork application command set.
const Version = 1

// A Command is a keyword (see Name field) passed to the anwork executable that provokes some
// functionality (see Action field).
type Command struct {
	Name, Description string

	// This slice holds the name(s) of the argument(s) that is(/are) expected by the Command.
	Args []string

	// This is the functionality that runs when this Command is invoked. The first parameter to the
	// function is the flag.FlagSet associated with this call to the Command. Implementers of the
	// Action function call pull command line arguments off of the flags parameter with
	// flag.FlagSet.Arg(X) where X is the index of the argument. Note that f.Arg(0) is always the name
	// of the command. The second parameter to this function is an output stream to which all output
	// (for example, debug printing, or stuff that would normally be sent to stdout) should be
	// written. The function should return true iff the task.Manager should be persisted to disk after
	// the Action returns.
	Action func(f *flag.FlagSet, o io.Writer, m *task.Manager) bool
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
		Name:        "set-priority",
		Description: "Set the priority of a task",
		Args:        []string{"task-name", "priority"},
		Action:      setPriorityAction,
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
		Description: "Show the journal; optionally pass a task name to only show events for that task",
		Args:        []string{"[task-name]"},
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

func versionAction(f *flag.FlagSet, o io.Writer, m *task.Manager) bool {
	fmt.Fprintln(o, "ANWORK Version =", Version)
	return false
}

func createAction(f *flag.FlagSet, o io.Writer, m *task.Manager) bool {
	name := f.Arg(1)
	m.Create(name)
	return true
}

func showAction(f *flag.FlagSet, o io.Writer, m *task.Manager) bool {
	printer := func(state task.State) {
		fmt.Fprintf(o, "%s tasks:\n", strings.ToUpper(task.StateNames[state]))
		for _, task := range m.Tasks() {
			if task.State() == state {
				fmt.Fprintf(o, "  %s (%d)\n", task.Name(), task.ID())
			}
		}
	}
	printer(task.StateRunning)
	printer(task.StateBlocked)
	printer(task.StateWaiting)
	printer(task.StateFinished)
	return false
}

func noteAction(f *flag.FlagSet, o io.Writer, m *task.Manager) bool {
	name := f.Arg(1)
	note := f.Arg(2)
	m.Note(name, note)
	return true
}

func deleteAction(f *flag.FlagSet, o io.Writer, m *task.Manager) bool {
	name := f.Arg(1)
	if !m.Delete(name) {
		fmt.Fprintf(o, "Error! Unknown task: %s\n", name)
		return false
	} else {
		return true
	}
}

func setPriorityAction(f *flag.FlagSet, o io.Writer, m *task.Manager) bool {
	name := f.Arg(1)
	priority := f.Arg(2)

	priorityInt, err := strconv.Atoi(priority)
	if err != nil {
		fmt.Fprintf(o, "Error! Could not parse priority from %s", priority)
		return false
	} else {
		m.SetPriority(name, priorityInt)
		return true
	}
}

func setStateAction(f *flag.FlagSet, o io.Writer, m *task.Manager) bool {
	name := f.Arg(1)

	var state task.State
	switch command := strings.TrimPrefix(f.Arg(0), "set-"); command {
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
	m.SetState(name, state)
	return true
}

func journalAction(f *flag.FlagSet, o io.Writer, m *task.Manager) bool {
	var t *task.Task = nil
	if f.NArg() > 1 {
		name := f.Arg(1)
		t = m.Find(name)
		if t == nil {
			panic("Unknown task: " + name)
		}
	}

	es := m.Journal().Events
	for i := len(es) - 1; i >= 0; i-- {
		e := es[i]
		if t == nil || t.ID() == e.TaskId {
			fmt.Fprintf(o, "[%s %s %d %02d:%02d]: %s\n", e.Date.Weekday(), e.Date.Month(), e.Date.Day(),
				e.Date.Hour(), e.Date.Minute(), e.Title)
		}
	}
	return false
}
