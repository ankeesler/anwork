// TODO: when we panic in this code, the command line interface looks really ugly. We should pass
// string error messages to the callers of these commands so that they are more nicely formatted
// when they appear on the command line.
package runner

import (
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ankeesler/anwork/task"
)

//go:generate go run ../../genclidoc/main.go ../../../doc/CLI.md

// This is the version of this anwork application command set.
const version = 3

// This error is used to indicate that the manager factory should be reset.
type resetError struct {
}

func (a resetError) Error() string {
	return ""
}

// A Command represents a keyword (see Name field) passed to the anwork executable that incites some
// behavior to run (via Command.Run).
type command struct {
	Name, Description string

	// This slice holds the name(s) of the argument(s) that is(/are) expected by the Command.
	Args []string

	// This is the functionality that runs when this Command is invoked. Note that args[0] is
	// always the Name of the command. The o parameter to this function is an output
	// stream to which all output should be written. The function should returns a non-nil error
	// iff an error occured.
	Action func(cmd *command, args []string, o io.Writer, m task.Manager) error
}

// These are the Command's used by the anwork application.
var commands = []command{
	command{
		Name:        "version",
		Description: "Print version information",
		Args:        []string{},
		Action:      versionAction,
	},
	command{
		Name:        "reset",
		Description: "Completely reset everything and blow away all data; USE CAREFULLY",
		Args:        []string{},
		Action:      resetAction,
	},
	//	command{
	//		Name:        "summary",
	//		Description: "Show a summary of the tasks completed in the past days",
	//		Args:        []string{"days"},
	//		action:      summaryAction,
	//	},
	command{
		Name:        "create",
		Description: "Create a new task",
		Args:        []string{"task-name"},
		Action:      createAction,
	},
	command{
		Name:        "delete",
		Description: "Delete a task",
		Args:        []string{"task-name"},
		Action:      deleteAction,
	},
	command{
		Name:        "delete-all",
		Description: "Delete all tasks",
		Args:        []string{},
		Action:      deleteAllAction,
	},
	command{
		Name:        "show",
		Description: "Show the current tasks, or the details of a specific task",
		Args:        []string{"[task-name]"},
		Action:      showAction,
	},
	command{
		Name:        "note",
		Description: "Add a note to a task",
		Args:        []string{"task-name", "note"},
		Action:      noteAction,
	},
	command{
		Name:        "set-priority",
		Description: "Set the priority of a task",
		Args:        []string{"task-name", "priority"},
		Action:      setPriorityAction,
	},
	//	command{
	//		Name:        "set-running",
	//		Description: "Mark a task as running",
	//		Args:        []string{"task-name"},
	//		action:      setStateAction,
	//	},
	//	command{
	//		Name:        "set-blocked",
	//		Description: "Mark a task as blocked",
	//		Args:        []string{"task-name"},
	//		action:      setStateAction,
	//	},
	//	command{
	//		Name:        "set-waiting",
	//		Description: "Mark a task as waiting",
	//		Args:        []string{"task-name"},
	//		action:      setStateAction,
	//	},
	//	command{
	//		Name:        "set-finished",
	//		Description: "Mark a task as finished",
	//		Args:        []string{"task-name"},
	//		action:      setStateAction,
	//	},
	//	command{
	//		Name:        "journal",
	//		Description: "Show the journal; optionally pass a task name to only show events for that task",
	//		Args:        []string{"[task-name]"},
	//		action:      journalAction,
	//	},
}

// Find the command with the provided name.
func findCommand(name string) *command {
	for _, c := range commands {
		if c.Name == name {
			return &c
		}
	}
	return nil
}

// Parse a "task spec" which is either the name of a task (i.e., "task-a") or the '@' symbol and an
// integer value indicating the ID of a task (i.e., "@37"). This function will never return nil. If
// the specifier is illegal, it will panic.
//func parseTaskSpec(str string, m task.Manager) *task.Task {
//	var t *task.Task = nil
//	if strings.HasPrefix(str, "@") {
//		num, err := strconv.Atoi(str[1:])
//		if err != nil {
//			panic("Error! Cannot parse task ID: " + err.Error())
//		}
//		for _, task := range m.Tasks() {
//			if task.ID == num {
//				t = task
//				break
//			}
//		}
//	} else {
//		t = m.FindByName(str) // str is the name of a task
//	}
//
//	if t == nil {
//		panic("Error! Unknown task for specifier: " + str)
//	}
//	return t
//}
//
func formatDate(seconds int64) string {
	date := time.Unix(seconds, 0)
	return fmt.Sprintf("%s %s %d %02d:%02d", date.Weekday(), date.Month(), date.Day(), date.Hour(),
		date.Minute())
}

//func formatDuration(duration time.Duration) string {
//	return fmt.Sprintf("%s", duration.String())
//}

func versionAction(cmd *command, args []string, o io.Writer, m task.Manager) error {
	fmt.Fprintln(o, "ANWORK Version =", version)
	return nil
}

func resetAction(cmd *command, args []string, o io.Writer, m task.Manager) error {
	fmt.Fprintf(o, "Are you sure you want to delete all data [y/n]: ")

	var answer string
	var ok bool
	if answer, ok = os.LookupEnv("ANWORK_TEST_RESET_ANSWER"); !ok {
		fmt.Scanf("%s", &answer)
	}

	if answer == "y" {
		fmt.Fprintln(o, "OK, deleting all data")
		return &resetError{}
	} else {
		fmt.Fprintln(o, "NOT deleting all data")
		return nil
	}
}

//
//func findCreateEvent(m task.Manager, taskID int) *task.Event {
//	for _, e := range m.Events() {
//		if e.Type == task.EventTypeCreate && e.TaskID == taskID {
//			return e
//		}
//	}
//	return nil
//}
//
//func summaryAction(args []string, o io.Writer, m task.Manager) response {
//	days, ok := arg(f, 1)
//	if !ok {
//		return responseArgumentError
//	}
//
//	daysNum, err := strconv.Atoi(days)
//	if err != nil {
//		msg := fmt.Sprintf("Cannot convert days %s to number: %s", days, err.Error())
//		panic(msg)
//	}
//
//	now := time.Now()
//	es := m.Events()
//	for i := len(es) - 1; i >= 0; i-- {
//		e := es[i]
//		isFinished := e.Type == task.EventTypeSetState && strings.Contains(e.Title, "to Finished")
//		eDate := time.Unix(e.Date, 0)
//		isWithinDays := eDate.Add(time.Duration(daysNum*24) * time.Hour).After(now)
//		if isFinished && isWithinDays {
//			createE := findCreateEvent(m, e.TaskID)
//			fmt.Fprintf(o, "[%s]: %s\n", formatDate(e.Date), e.Title)
//			if createE != nil {
//				createEDate := time.Unix(createE.Date, 0)
//				fmt.Fprintf(o, "  took %s\n", formatDuration(eDate.Sub(createEDate)))
//			}
//		}
//	}
//	return responseNoPersist
//}

func createAction(cmd *command, args []string, o io.Writer, m task.Manager) error {
	name := args[1]
	if err := m.Create(name); err != nil {
		return err
	}

	return nil
}

func deleteAction(cmd *command, args []string, o io.Writer, m task.Manager) error {
	spec := args[1]

	//t := parseTaskSpec(spec, m)
	if !m.Delete(spec) {
		fmt.Fprintf(o, "Error! Unknown task: %s\n", spec)
	}
	return nil
}

func deleteAllAction(cmd *command, args []string, o io.Writer, m task.Manager) error {
	tasks := m.Tasks()
	for i := 0; i < len(tasks); i++ {
		name := tasks[i].Name
		if !m.Delete(name) {
			fmt.Fprintf(o, "Error! Unable to delete task %s", name)
		}
	}
	return nil
}

func showAction(cmd *command, args []string, o io.Writer, m task.Manager) error {
	if len(args) == 1 {
		tasks := m.Tasks()
		printer := func(state task.State) {
			fmt.Fprintf(o, "%s tasks:\n", strings.ToUpper(task.StateNames[state]))
			for _, task := range tasks {
				if task.State == state {
					fmt.Fprintf(o, "  %s (%d)\n", task.Name, task.ID)
				}
			}
		}
		printer(task.StateRunning)
		printer(task.StateBlocked)
		printer(task.StateWaiting)
		printer(task.StateFinished)
	} else {
		t := m.FindByName(args[1])
		if t == nil {
			fmt.Fprintf(o, "Error! Unknown task: %s", args[1])
			return nil
		}

		fmt.Fprintf(o, "Name: %s\n", t.Name)
		fmt.Fprintf(o, "ID: %d\n", t.ID)
		fmt.Fprintf(o, "Created: %s\n", formatDate(t.StartDate))
		fmt.Fprintf(o, "Priority: %d\n", t.Priority)
		fmt.Fprintf(o, "State: %s\n", strings.ToUpper(task.StateNames[t.State]))
	}
	return nil
}

func noteAction(cmd *command, args []string, o io.Writer, m task.Manager) error {
	//t := parseTaskSpec(spec, m)
	err := m.Note(args[1], args[2])
	if err != nil {
		fmt.Fprintf(o, "Error! Cannot add note: %s", err.Error())
	}
	return nil
}

func setPriorityAction(cmd *command, args []string, o io.Writer, m task.Manager) error {
	prio, err := strconv.Atoi(args[2])
	if err != nil {
		fmt.Fprintf(o, "Error! Cannot set priority: invalid priority: '%s'", args[2])
		return nil
	}

	if err := m.SetPriority(args[1], prio); err != nil {
		fmt.Fprintf(o, "Error! Cannot set priority: %s", err.Error())
	}

	return nil
}

//
//func setStateAction(args []string, o io.Writer, m task.Manager) response {
//	spec, ok := arg(f, 1)
//	if !ok {
//		return responseArgumentError
//	}
//
//	t := parseTaskSpec(spec, m)
//
//	var state task.State
//	switch command := strings.TrimPrefix(f.Arg(0), "set-"); command {
//	case "running":
//		state = task.StateRunning
//	case "blocked":
//		state = task.StateBlocked
//	case "waiting":
//		state = task.StateWaiting
//	case "finished":
//		state = task.StateFinished
//	default:
//		panic("Unknown state: " + command)
//	}
//	m.SetState(t.Name, state)
//	return responsePersist
//}
//
//func journalAction(args []string, o io.Writer, m task.Manager) response {
//	var t *task.Task = nil
//	if f.NArg() > 1 {
//		t = parseTaskSpec(f.Arg(1), m)
//	}
//
//	es := m.Events()
//	for i := len(es) - 1; i >= 0; i-- {
//		e := es[i]
//		if t == nil || t.ID == e.TaskID {
//			fmt.Fprintf(o, "[%s]: %s\n", formatDate(e.Date), e.Title)
//		}
//	}
//	return responseNoPersist
//}
