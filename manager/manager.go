// Package manager contains an interface through which task.Task'fs can be created, read, updated, and deleted.
package manager

import (
	"fmt"

	"code.cloudfoundry.org/clock"
	"github.com/ankeesler/anwork/task2"
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
	return m.repo.CreateTask(&task)
}

func (m *manager) Delete(name string) error {
	task, err := m.repo.FindTaskByName(name)
	if err != nil {
		return err
	}

	if task == nil {
		return fmt.Errorf("unknown task with name '%s'", name)
	}

	return m.repo.DeleteTask(task)
}

func (m *manager) FindByID(id int) (*task2.Task, error)          { return nil, nil }
func (m *manager) FindByName(name string) (*task2.Task, error)   { return nil, nil }
func (m *manager) Tasks() ([]*task2.Task, error)                 { return nil, nil }
func (m *manager) Note(name, note string) error                  { return nil }
func (m *manager) SetPriority(name string, priority int) error   { return nil }
func (m *manager) SetState(name string, state task2.State) error { return nil }
func (m *manager) Events() ([]*task2.Event, error)               { return nil, nil }
func (m *manager) DeleteEvent(startTime int64) error             { return nil }
func (m *manager) Reset() error                                  { return nil }
func (m *manager) Rename(from, to string) error                  { return nil }
