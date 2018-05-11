package runner

import (
	"errors"
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

type unknownTaskError struct {
	name string
}

func (u unknownTaskError) Error() string {
	return fmt.Sprintf("unknown task: %s", u.name)
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
	command{
		Name:        "summary",
		Description: "Show a summary of the tasks completed in the past days",
		Args:        []string{"days"},
		Action:      summaryAction,
	},
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
	command{
		Name:        "set-running",
		Description: "Mark a task as running",
		Args:        []string{"task-name"},
		Action:      setStateAction,
	},
	command{
		Name:        "set-blocked",
		Description: "Mark a task as blocked",
		Args:        []string{"task-name"},
		Action:      setStateAction,
	},
	command{
		Name:        "set-waiting",
		Description: "Mark a task as waiting",
		Args:        []string{"task-name"},
		Action:      setStateAction,
	},
	command{
		Name:        "set-finished",
		Description: "Mark a task as finished",
		Args:        []string{"task-name"},
		Action:      setStateAction,
	},
	command{
		Name:        "journal",
		Description: "Show the journal; optionally pass a task name to only show events for that task",
		Args:        []string{"[task-name]"},
		Action:      journalAction,
	},
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
// integer value indicating the ID of a task (i.e., "@37").
func parseTaskSpec(str string, m task.Manager) (*task.Task, error) {
	var t *task.Task = nil
	if strings.HasPrefix(str, "@") {
		num, err := strconv.Atoi(str[1:])
		if err != nil {
			return nil, fmt.Errorf("cannot parse task ID: %s", err.Error())
		}

		t = m.FindByID(num)
		if t == nil {
			return nil, fmt.Errorf("unknown task ID in task spec: %d", num)
		}
	} else {
		t = m.FindByName(str) // str is the name of a task
		if t == nil {
			return nil, unknownTaskError{name: str}
		}
	}

	return t, nil
}

func formatDate(seconds int64) string {
	date := time.Unix(seconds, 0)
	return fmt.Sprintf("%s %s %d %02d:%02d", date.Weekday(), date.Month(), date.Day(), date.Hour(),
		date.Minute())
}

func formatDuration(duration time.Duration) string {
	return fmt.Sprintf("%s", duration.String())
}

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

func findCreateEvent(m task.Manager, taskID int) *task.Event {
	for _, e := range m.Events() {
		if e.Type == task.EventTypeCreate && e.TaskID == taskID {
			return e
		}
	}
	return nil
}

func summaryAction(cmd *command, args []string, o io.Writer, m task.Manager) error {
	daysNum, err := strconv.Atoi(args[1])
	if err != nil {
		return fmt.Errorf("Cannot convert days %s to number: %s", args[1], err.Error())
	}
	_ = daysNum

	now := time.Now()
	es := m.Events()
	for i := len(es) - 1; i >= 0; i-- {
		e := es[i]
		isFinished := e.Type == task.EventTypeSetState && strings.Contains(e.Title, "to Finished")
		eDate := time.Unix(e.Date, 0)
		isWithinDays := eDate.Add(time.Duration(daysNum*24) * time.Hour).After(now)
		if isFinished && isWithinDays {
			createE := findCreateEvent(m, e.TaskID)
			fmt.Fprintf(o, "[%s]: %s\n", formatDate(e.Date), e.Title)
			if createE != nil {
				createEDate := time.Unix(createE.Date, 0)
				fmt.Fprintf(o, "  took %s\n", formatDuration(eDate.Sub(createEDate)))
			}
		}
	}

	return nil
}

func createAction(cmd *command, args []string, o io.Writer, m task.Manager) error {
	name := args[1]
	if err := m.Create(name); err != nil {
		return err
	}

	return nil
}

func deleteAction(cmd *command, args []string, o io.Writer, m task.Manager) error {
	spec := args[1]

	t, err := parseTaskSpec(spec, m)
	if err != nil {
		return err
	}

	if !m.Delete(t.Name) {
		return unknownTaskError{name: spec}
	}

	return nil
}

func deleteAllAction(cmd *command, args []string, o io.Writer, m task.Manager) error {
	tasks := m.Tasks()
	errMsgs := []string{}
	for i := 0; i < len(tasks); i++ {
		name := tasks[i].Name
		if !m.Delete(name) {
			msg := fmt.Sprintf("\n\tunable to delete task %s", name)
			errMsgs = append(errMsgs, msg)
		}
	}

	if len(errMsgs) > 0 {
		return errors.New(strings.Join(errMsgs, ""))
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
		t, err := parseTaskSpec(args[1], m)
		if err != nil {
			return err
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
	t, err := parseTaskSpec(args[1], m)
	if err != nil {
		return err
	}

	if err = m.Note(t.Name, args[2]); err != nil {
		return fmt.Errorf("cannot add note: %s", err.Error())
	}

	return nil
}

func setPriorityAction(cmd *command, args []string, o io.Writer, m task.Manager) error {
	t, err := parseTaskSpec(args[1], m)
	if err != nil {
		return err
	}

	prio, err := strconv.Atoi(args[2])
	if err != nil {
		return fmt.Errorf("cannot set priority: invalid priority: '%s'", args[2])
	}

	if err := m.SetPriority(t.Name, prio); err != nil {
		return fmt.Errorf("cannot set priority: %s", err.Error())
	}

	return nil
}

func setStateAction(cmd *command, args []string, o io.Writer, m task.Manager) error {
	t, err := parseTaskSpec(args[1], m)
	if err != nil {
		return err
	}

	var state task.State
	switch command := strings.TrimPrefix(args[0], "set-"); command {
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

	if err := m.SetState(t.Name, state); err != nil {
		return fmt.Errorf("cannot set state: %s", err.Error())
	}

	return nil
}

func journalAction(cmd *command, args []string, o io.Writer, m task.Manager) error {
	var t *task.Task = nil
	if len(args) > 1 {
		var err error
		t, err = parseTaskSpec(args[1], m)
		if err != nil {
			return err
		}
	}

	es := m.Events()
	for i := len(es) - 1; i >= 0; i-- {
		e := es[i]
		if t == nil || t.ID == e.TaskID {
			fmt.Fprintf(o, "[%s]: %s\n", formatDate(e.Date), e.Title)
		}
	}
	return nil
}
