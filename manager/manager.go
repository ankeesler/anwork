// Package manager contains an interface through which task.Task'fs can be created, read, updated, and deleted.
package manager

import (
	"fmt"

	"code.cloudfoundry.org/clock"
	"github.com/ankeesler/anwork/task2"
	multierror "github.com/hashicorp/go-multierror"
)

//go:generate counterfeiter . Manager

// A Manager is an interface through which Task's can be created, read, updated, and deleted.
type Manager interface {
	// Create a task with a name. Return an error if the task name is not unique.
	Create(name string) error

	// Delete a task with a name. Returns an error if the task was not able to be deleted.
	Delete(name string) error

	// Find a task with an ID.
	FindByID(id int) (*task2.Task, error)
	// Find a task with a name.
	FindByName(name string) (*task2.Task, error)

	// Get all of the Tasks contained in this manager, ordered from highest priority (lowest integer
	// value) to lowest priority (highest integer value).
	//
	// When multiple tasks have the same priority, the Task's will be ordered by their (unique) ID in
	// ascending order. This means that the older Task's will come first. This is a conscious decision.
	// The Task's that have been around the longest are assumed to need to be completed first.
	//
	// This function will never return nil!
	Tasks() ([]*task2.Task, error)

	// Add a note for a task.
	Note(name, note string) error
	// Set the priority of a task.
	SetPriority(name string, priority int) error
	// Set the state of a task.
	SetState(name string, state task2.State) error

	// Get the events associated with this manager.
	Events() ([]*task2.Event, error)
	// Delete an event, identified by its start time.
	DeleteEvent(startTime int64) error

	// Perform a factory reset, e.g., make this manager new again.
	Reset() error

	// Rename a task.
	Rename(from, to string) error
}

const defaultPriority = 10

const defaultState = task2.StateReady

type manager struct {
	repo  task2.Repo
	clock clock.Clock
}

// New creates a new Manager that will use a task.Repo for CRUD task.Task operations.
func New(repo task2.Repo, clock clock.Clock) Manager {
	return &manager{repo: repo, clock: clock}
}

func (m *manager) Create(name string) error {
	task := task2.Task{
		Name:      name,
		StartDate: m.clock.Now().Unix(),
		Priority:  defaultPriority,
		State:     defaultState,
	}
	if err := m.repo.CreateTask(&task); err != nil {
		return err
	}

	return m.repo.CreateEvent(&task2.Event{
		Title:  fmt.Sprintf("Created task '%s'", name),
		Date:   m.clock.Now().Unix(),
		Type:   task2.EventTypeCreate,
		TaskID: task.ID,
	})
}

func (m *manager) Delete(name string) error {
	return m.doWithTask(name, func(task *task2.Task) error {
		if err := m.repo.DeleteTask(task); err != nil {
			return err
		}

		return m.repo.CreateEvent(&task2.Event{
			Title:  fmt.Sprintf("Deleted task '%s'", name),
			Date:   m.clock.Now().Unix(),
			Type:   task2.EventTypeDelete,
			TaskID: task.ID,
		})
	})
}

func (m *manager) FindByID(id int) (*task2.Task, error) {
	return m.repo.FindTaskByID(id)
}

func (m *manager) FindByName(name string) (*task2.Task, error) {
	return m.repo.FindTaskByName(name)
}

func (m *manager) Tasks() ([]*task2.Task, error) {
	return m.repo.Tasks()
}

func (m *manager) Note(name, note string) error {
	return m.doWithTask(name, func(task *task2.Task) error {
		return m.repo.CreateEvent(&task2.Event{
			Title:  fmt.Sprintf("Note: %s", note),
			Date:   m.clock.Now().Unix(),
			Type:   task2.EventTypeNote,
			TaskID: task.ID,
		})
	})
}

func (m *manager) SetPriority(name string, priority int) error {
	return m.doWithTask(name, func(task *task2.Task) error {
		oldPriority := task.Priority
		task.Priority = priority
		if err := m.repo.UpdateTask(task); err != nil {
			return err
		}

		return m.repo.CreateEvent(&task2.Event{
			Title: fmt.Sprintf("Set priority on task '%s' from %d to %d",
				name, oldPriority, priority),
			Date:   m.clock.Now().Unix(),
			Type:   task2.EventTypeNote,
			TaskID: task.ID,
		})
	})
}

func (m *manager) SetState(name string, state task2.State) error {
	return m.doWithTask(name, func(task *task2.Task) error {
		oldState := task.State
		task.State = state
		if err := m.repo.UpdateTask(task); err != nil {
			return err
		}

		return m.repo.CreateEvent(&task2.Event{
			Title: fmt.Sprintf("Set state on task '%s' from %s to %s",
				name, task2.StateNames[oldState], task2.StateNames[state]),
			Date:   m.clock.Now().Unix(),
			Type:   task2.EventTypeNote,
			TaskID: task.ID,
		})
	})
}

func (m *manager) Events() ([]*task2.Event, error) {
	return m.repo.Events()
}

func (m *manager) DeleteEvent(startTime int64) error {
	event := task2.Event{Date: startTime}
	return m.repo.DeleteEvent(&event)
}

func (m *manager) Reset() error {
	tasks, err := m.repo.Tasks()
	if err != nil {
		return err
	}

	events, err := m.repo.Events()
	if err != nil {
		return err
	}

	var result *multierror.Error
	for _, task := range tasks {
		result = multierror.Append(result, m.repo.DeleteTask(task))
	}

	for _, event := range events {
		result = multierror.Append(result, m.repo.DeleteEvent(event))
	}

	return result.ErrorOrNil()
}

func (m *manager) Rename(from, to string) error {
	return m.doWithTask(from, func(task *task2.Task) error {
		task.Name = to
		return m.repo.UpdateTask(task)
	})
}

func (m *manager) doWithTask(name string, do func(*task2.Task) error) error {
	task, err := m.repo.FindTaskByName(name)
	if err != nil {
		return err
	}

	if task == nil {
		return unknownTaskError{name: name}
	}

	return do(task)
}
