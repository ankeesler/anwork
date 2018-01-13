// This package contains the commands that can be passed to the anwork command line tool.
package command

import (
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/ankeesler/anwork/task"
)

//go:generate go run ../../genclidoc/main.go ../../../doc/CLI.md

// This is the version of this anwork application command set.
const Version = 2

// This is an indication of what the command wishes to happen after being run. See Response*
// constants below for more details.
type Response int

// These constants represent the valid Response values.
const (
	// Do not persist any data.
	ResponseNoPersist = iota
	// Persist data back to disk.
	ResponsePersist
	// Completely wipe out all data on disk. This is dangerous!
	ResponseReset
	// There was an error in the arguments passed (or not passed) to the Command.
	ResponseArgumentError
)

// A Command represents a keyword (see Name field) passed to the anwork executable that incites some
// behavior to run (via Command.Run).
type Command struct {
	Name, Description string

	// This slice holds the name(s) of the argument(s) that is(/are) expected by the Command.
	Args []string

	// This is the functionality that runs when this Command is invoked. The first parameter to the
	// function is the flag.FlagSet associated with this call to the Command. Implementers of the
	// action function call pull command line arguments off of the flags parameter with
	// flag.FlagSet.Arg(X) where X is the index of the argument. Note that f.Arg(0) is always the name
	// of the command. The second parameter to this function is an output stream to which all output
	// (for example, debug printing, or stuff that would normally be sent to stdout) should be
	// written. The function should return a Response value based on the next action that the caller
	// should take.
	action func(f *flag.FlagSet, o io.Writer, m *task.Manager) Response
}

// These are the Command's used by the anwork application.
var Commands = []Command{
	Command{
		Name:        "version",
		Description: "Print version information",
		Args:        []string{},
		action:      versionAction,
	},
	Command{
		Name:        "reset",
		Description: "Completely reset everything and blow away all data; USE CAREFULLY",
		Args:        []string{},
		action:      resetAction,
	},
	Command{
		Name:        "summary",
		Description: "Show a summary of the tasks completed in the past days",
		Args:        []string{"days"},
		action:      summaryAction,
	},
	Command{
		Name:        "create",
		Description: "Create a new task",
		Args:        []string{"task-name"},
		action:      createAction,
	},
	Command{
		Name:        "delete",
		Description: "Delete a task",
		Args:        []string{"task-name"},
		action:      deleteAction,
	},
	Command{
		Name:        "delete-all",
		Description: "Delete all tasks",
		Args:        []string{},
		action:      deleteAllAction,
	},
	Command{
		Name:        "show",
		Description: "Show the current tasks, or the details of a specific task",
		Args:        []string{"[task-name]"},
		action:      showAction,
	},
	Command{
		Name:        "note",
		Description: "Add a note to a task",
		Args:        []string{"task-name", "note"},
		action:      noteAction,
	},
	Command{
		Name:        "set-priority",
		Description: "Set the priority of a task",
		Args:        []string{"task-name", "priority"},
		action:      setPriorityAction,
	},
	Command{
		Name:        "set-running",
		Description: "Mark a task as running",
		Args:        []string{"task-name"},
		action:      setStateAction,
	},
	Command{
		Name:        "set-blocked",
		Description: "Mark a task as blocked",
		Args:        []string{"task-name"},
		action:      setStateAction,
	},
	Command{
		Name:        "set-waiting",
		Description: "Mark a task as waiting",
		Args:        []string{"task-name"},
		action:      setStateAction,
	},
	Command{
		Name:        "set-finished",
		Description: "Mark a task as finished",
		Args:        []string{"task-name"},
		action:      setStateAction,
	},
	Command{
		Name:        "journal",
		Description: "Show the journal; optionally pass a task name to only show events for that task",
		Args:        []string{"[task-name]"},
		action:      journalAction,
	},
}

// Print the usage information out for a single Command. The information will be printed to the
// provided output stream.
func (c *Command) Usage(output io.Writer) {
	fmt.Fprintf(output, "  %s %s\n", c.Name, strings.Join(c.Args, " "))
	fmt.Fprintf(output, "        %s\n", c.Description)
}

// Run the functionality associated with this Command. This function will return a Response value
// indicating next steps for the caller to take.
func (c *Command) Run(f *flag.FlagSet, o io.Writer, m *task.Manager) Response {
	return c.action(f, o, m)
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

// Parse a "task spec" which is either the name of a task (i.e., "task-a") or the '@' symbol and an
// integer value indicating the ID of a task (i.e., "@37"). This function will never return nil. If
// the specifier is illegal, it will panic.
func parseTaskSpec(str string, m *task.Manager) *task.Task {
	var t *task.Task = nil
	if strings.HasPrefix(str, "@") {
		num, err := strconv.Atoi(str[1:])
		if err != nil {
			panic("Error! Cannot parse task ID: " + err.Error())
		}
		for _, task := range m.Tasks() {
			if task.ID() == num {
				t = task
				break
			}
		}
	} else {
		t = m.FindByName(str) // str is the name of a task
	}

	if t == nil {
		panic("Error! Unknown task for specifier: " + str)
	}
	return t
}

func formatDate(date time.Time) string {
	return fmt.Sprintf("%s %s %d %02d:%02d", date.Weekday(), date.Month(), date.Day(), date.Hour(),
		date.Minute())
}

func formatDuration(duration time.Duration) string {
	return fmt.Sprintf("%s", duration.String())
}

// Get the arg at index i or return false if it doesn't exist.
func arg(f *flag.FlagSet, i int) (string, bool) {
	if f.NArg() <= i {
		return "", false
	} else {
		return f.Arg(i), true
	}
}

func versionAction(f *flag.FlagSet, o io.Writer, m *task.Manager) Response {
	fmt.Fprintln(o, "ANWORK Version =", Version)
	return ResponseNoPersist
}

func resetAction(f *flag.FlagSet, o io.Writer, m *task.Manager) Response {
	fmt.Fprintf(o, "Are you sure you want to delete all data [y/n]: ")

	var answer string
	if f.NArg() > 1 {
		answer = f.Arg(1) // Trapdoor to ease testing
	} else {
		fmt.Scanf("%s", &answer)
	}

	if answer == "y" {
		fmt.Fprintln(o, "OK, deleting all data")
		return ResponseReset
	} else {
		fmt.Fprintln(o, "NOT deleting all data")
		return ResponseNoPersist
	}
}

func summaryAction(f *flag.FlagSet, o io.Writer, m *task.Manager) Response {
	days, ok := arg(f, 1)
	if !ok {
		return ResponseArgumentError
	}

	daysNum, err := strconv.Atoi(days)
	if err != nil {
		msg := fmt.Sprintf("Cannot convert days %s to number: %s", days, err.Error())
		panic(msg)
	}

	now := time.Now()
	es := m.Journal().Events
	for i := len(es) - 1; i >= 0; i-- {
		e := es[i]
		isFinished := e.Type == task.EventTypeSetState && strings.Contains(e.Title, "to Finished")
		isWithinDays := e.Date.Add(time.Duration(daysNum*24) * time.Hour).After(now)
		if isFinished && isWithinDays {
			t := m.FindById(e.TaskId)
			fmt.Fprintf(o, "[%s]: %s\n", formatDate(e.Date), e.Title)
			fmt.Fprintf(o, "  took %s\n", formatDuration(e.Date.Sub(t.StartDate())))
		}
	}
	return ResponseNoPersist
}

func createAction(f *flag.FlagSet, o io.Writer, m *task.Manager) Response {
	name, ok := arg(f, 1)
	if !ok {
		return ResponseArgumentError
	}

	m.Create(name)
	return ResponsePersist
}

func showAction(f *flag.FlagSet, o io.Writer, m *task.Manager) Response {
	var t *task.Task = nil
	if f.NArg() > 1 {
		t = parseTaskSpec(f.Arg(1), m)
	}

	if t == nil {
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
	} else {
		fmt.Fprintf(o, "Name: %s\n", t.Name())
		fmt.Fprintf(o, "ID: %d\n", t.ID())
		fmt.Fprintf(o, "Created: %s\n", formatDate(t.StartDate()))
		fmt.Fprintf(o, "Priority: %d\n", t.Priority())
		fmt.Fprintf(o, "State: %s\n", strings.ToUpper(task.StateNames[t.State()]))
	}
	return ResponseNoPersist
}

func noteAction(f *flag.FlagSet, o io.Writer, m *task.Manager) Response {
	spec, ok := arg(f, 1)
	if !ok {
		return ResponseArgumentError
	}
	note, ok := arg(f, 2)
	if !ok {
		return ResponseArgumentError
	}

	t := parseTaskSpec(spec, m)
	m.Note(t.Name(), note)
	return ResponsePersist
}

func deleteAction(f *flag.FlagSet, o io.Writer, m *task.Manager) Response {
	spec, ok := arg(f, 1)
	if !ok {
		return ResponseArgumentError
	}

	t := parseTaskSpec(spec, m)
	if !m.Delete(t.Name()) {
		fmt.Fprintf(o, "Error! Unknown task: %s\n", t.Name())
		return ResponseNoPersist
	} else {
		return ResponsePersist
	}
}

func deleteAllAction(f *flag.FlagSet, o io.Writer, m *task.Manager) Response {
	for len(m.Tasks()) > 0 {
		name := m.Tasks()[0].Name()
		if !m.Delete(name) {
			panic("Expected to be able to successfully delete task " + name)
		}
	}
	return ResponsePersist
}

func setPriorityAction(f *flag.FlagSet, o io.Writer, m *task.Manager) Response {
	spec, ok := arg(f, 1)
	if !ok {
		return ResponseArgumentError
	}

	t := parseTaskSpec(spec, m)
	priority := f.Arg(2)

	priorityInt, err := strconv.Atoi(priority)
	if err != nil {
		fmt.Fprintf(o, "Error! Could not parse priority from %s", priority)
		return ResponseNoPersist
	} else {
		m.SetPriority(t.Name(), priorityInt)
		return ResponsePersist
	}
}

func setStateAction(f *flag.FlagSet, o io.Writer, m *task.Manager) Response {
	spec, ok := arg(f, 1)
	if !ok {
		return ResponseArgumentError
	}

	t := parseTaskSpec(spec, m)

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
	m.SetState(t.Name(), state)
	return ResponsePersist
}

func journalAction(f *flag.FlagSet, o io.Writer, m *task.Manager) Response {
	var t *task.Task = nil
	if f.NArg() > 1 {
		t = parseTaskSpec(f.Arg(1), m)
	}

	es := m.Journal().Events
	for i := len(es) - 1; i >= 0; i-- {
		e := es[i]
		if t == nil || t.ID() == e.TaskId {
			fmt.Fprintf(o, "[%s]: %s\n", formatDate(e.Date), e.Title)
		}
	}
	return ResponseNoPersist
}
