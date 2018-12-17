// Package manager contains an interface through which task.Task's can be created,
// read, updated, and deleted.
package manager

import (
	"fmt"
	"sort"

	"code.cloudfoundry.org/clock"
	"github.com/ankeesler/anwork/task"
	taskpkg "github.com/ankeesler/anwork/task"
	multierror "github.com/hashicorp/go-multierror"
)

//go:generate counterfeiter . Manager

// A Manager is an interface through which Task's can be created, read, updated, and deleted.
type Manager interface {
	// Create a task with a name. Returns an error if the task name is not unique.
	// All Tasks start out in the task.StateReady task.State.
	Create(name string) error

	// Delete a task with a name. Returns an error if the task was not able to be deleted.
	Delete(name string) error

	// Find a task with an ID.
	FindByID(id int) (*taskpkg.Task, error)
	// Find a task with a name.
	FindByName(name string) (*taskpkg.Task, error)

	// Get all of the Tasks contained in this manager, ordered from highest priority (lowest integer
	// value) to lowest priority (highest integer value).
	//
	// When multiple tasks have the same priority, the Task's will be ordered by their (unique) ID in
	// ascending order. This means that the older Task's will come first. This is a conscious decision.
	// The Task's that have been around the longest are assumed to need to be completed first.
	Tasks() ([]*taskpkg.Task, error)

	// Add a note for a task.
	Note(name, note string) error
	// Set the priority of a task.
	SetPriority(name string, priority int) error
	// Set the state of a task.
	SetState(name string, state taskpkg.State) error

	// Get the events associated with this manager.
	Events() ([]*taskpkg.Event, error)

	// Perform a factory reset, e.g., make this manager new again.
	Reset() error

	// Rename a task.
	Rename(from, to string) error
}

const defaultPriority = 10

const defaultState = taskpkg.StateReady

type manager struct {
	repo  taskpkg.Repo
	clock clock.Clock
}

// New creates a new Manager that will use a task.Repo for CRUD task.Task operations.
func New(repo taskpkg.Repo, clock clock.Clock) Manager {
	return &manager{repo: repo, clock: clock}
}

func (m *manager) Create(name string) error {
	task := taskpkg.Task{
		Name:      name,
		StartDate: m.clock.Now().Unix(),
		Priority:  defaultPriority,
		State:     defaultState,
	}
	if err := m.repo.CreateTask(&task); err != nil {
		return err
	}

	return m.repo.CreateEvent(&taskpkg.Event{
		Title:  fmt.Sprintf("Created task '%s'", name),
		Date:   m.clock.Now().Unix(),
		Type:   taskpkg.EventTypeCreate,
		TaskID: task.ID,
	})
}

func (m *manager) Delete(name string) error {
	return m.doWithTask(name, func(task *taskpkg.Task) error {
		if err := m.repo.DeleteTask(task); err != nil {
			return err
		}

		return m.repo.CreateEvent(&taskpkg.Event{
			Title:  fmt.Sprintf("Deleted task '%s'", name),
			Date:   m.clock.Now().Unix(),
			Type:   taskpkg.EventTypeDelete,
			TaskID: task.ID,
		})
	})
}

func (m *manager) FindByID(id int) (*taskpkg.Task, error) {
	return m.repo.FindTaskByID(id)
}

func (m *manager) FindByName(name string) (*taskpkg.Task, error) {
	return m.repo.FindTaskByName(name)
}

func (m *manager) Tasks() ([]*taskpkg.Task, error) {
	tasks, err := m.repo.Tasks()
	if err != nil {
		return nil, err
	}

	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].Priority == tasks[j].Priority {
			return tasks[i].ID < tasks[j].ID
		} else {
			return tasks[i].Priority < tasks[j].Priority
		}
	})

	return tasks, nil
}

func (m *manager) Note(name, note string) error {
	return m.doWithTask(name, func(task *taskpkg.Task) error {
		return m.repo.CreateEvent(&taskpkg.Event{
			Title:  fmt.Sprintf("Note added to task '%s': %s", name, note),
			Date:   m.clock.Now().Unix(),
			Type:   taskpkg.EventTypeNote,
			TaskID: task.ID,
		})
	})
}

func (m *manager) SetPriority(name string, priority int) error {
	return m.doWithTask(name, func(task *taskpkg.Task) error {
		oldPriority := task.Priority
		task.Priority = priority
		if err := m.repo.UpdateTask(task); err != nil {
			return err
		}

		return m.repo.CreateEvent(&taskpkg.Event{
			Title: fmt.Sprintf("Set priority on task '%s' from %d to %d",
				name, oldPriority, priority),
			Date:   m.clock.Now().Unix(),
			Type:   taskpkg.EventTypeSetPriority,
			TaskID: task.ID,
		})
	})
}

func (m *manager) SetState(name string, state task.State) error {
	return m.doWithTask(name, func(task *task.Task) error {
		oldState := task.State
		task.State = state
		if err := m.repo.UpdateTask(task); err != nil {
			return err
		}

		return m.repo.CreateEvent(&taskpkg.Event{
			Title:  fmt.Sprintf("Set state on task '%s' from %s to %s", name, oldState, state),
			Date:   m.clock.Now().Unix(),
			Type:   taskpkg.EventTypeSetState,
			TaskID: task.ID,
		})
	})
}

func (m *manager) Events() ([]*task.Event, error) {
	return m.repo.Events()
}

func (m *manager) Reset() error {
	tasks, err := m.repo.Tasks()
	if err != nil {
		return err
	}
	tasksSize := len(tasks)

	events, err := m.repo.Events()
	if err != nil {
		return err
	}
	eventsSize := len(events)

	var result *multierror.Error
	for i := tasksSize - 1; i >= 0; i-- {
		result = multierror.Append(result, m.repo.DeleteTask(tasks[i]))
	}

	for i := eventsSize - 1; i >= 0; i-- {
		result = multierror.Append(result, m.repo.DeleteEvent(events[i]))
	}

	return result.ErrorOrNil()
}

func (m *manager) Rename(from, to string) error {
	return m.doWithTask(from, func(task *task.Task) error {
		task.Name = to
		return m.repo.UpdateTask(task)
	})
}

func (m *manager) doWithTask(name string, do func(*task.Task) error) error {
	task, err := m.repo.FindTaskByName(name)
	if err != nil {
		return err
	}

	if task == nil {
		return unknownTaskError{name: name}
	}

	return do(task)
}
