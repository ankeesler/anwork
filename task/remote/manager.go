// This package contains a task.Manager that acts as an HTTP client for the anwork
// API. All of the state that this type manipulates lives on a remote system.
package remote

import (
	"fmt"
	"strings"
	"time"

	"github.com/ankeesler/anwork/task"
)

//go:generate counterfeiter . APIClient
type APIClient interface {
	CreateTask(string) error
	DeleteTask(int) error
	GetTasks() ([]*task.Task, error)
	GetTask(int) (*task.Task, error)
	UpdatePriority(int, int) error
	UpdateState(int, task.State) error
	CreateEvent(string, task.EventType, int64, int) error
	GetEvents() ([]*task.Event, error)
	DeleteEvent(int64) error
}

type manager struct {
	client APIClient
}

func newManager(client APIClient) *manager {
	return &manager{client: client}
}

type unknownTaskError struct {
	name string
}

func (ute *unknownTaskError) Error() string {
	return fmt.Sprintf("Unknown task with name %s", ute.name)
}

func (m *manager) Create(name string) error {
	return m.client.CreateTask(name)
}

func (m *manager) Delete(name string) error {
	task := m.FindByName(name)
	if task == nil {
		return &unknownTaskError{name: name}
	}
	return m.client.DeleteTask(task.ID)
}

func (m *manager) Tasks() []*task.Task {
	tasks, err := m.client.GetTasks()
	if err != nil {
		panic(err)
	}
	return tasks
}

func (m *manager) FindByName(name string) *task.Task {
	tasks, err := m.client.GetTasks()
	if err != nil {
		panic(err)
	}

	var foundIt *task.Task
	for _, task := range tasks {
		if task.Name == name {
			foundIt = task
			break
		}
	}

	return foundIt
}

func (m *manager) FindByID(id int) *task.Task {
	task, err := m.client.GetTask(id)
	if err != nil {
		panic(err)
	}
	return task
}

func (m *manager) SetPriority(name string, prio int) error {
	task := m.FindByName(name)
	if task == nil {
		return &unknownTaskError{name: name}
	}
	return m.client.UpdatePriority(task.ID, prio)
}

func (m *manager) SetState(name string, state task.State) error {
	task := m.FindByName(name)
	if task == nil {
		return &unknownTaskError{name: name}
	}
	return m.client.UpdateState(task.ID, state)
}

func (m *manager) Note(name, note string) error {
	t := m.FindByName(name)
	if t == nil {
		return &unknownTaskError{name: name}
	}

	return m.client.CreateEvent(note, task.EventTypeNote, time.Now().Unix(), t.ID)
}

func (m *manager) DeleteEvent(startTime int64) error {
	return m.client.DeleteEvent(startTime)
}

func (m *manager) Events() []*task.Event {
	events, err := m.client.GetEvents()
	if err != nil {
		panic(err)
	}
	return events
}

func (m *manager) Reset() error {
	tasks, err := m.client.GetTasks()
	if err != nil {
		return err
	}

	events, err := m.client.GetEvents()
	if err != nil {
		return err
	}

	errs := make(map[string]string)
	for _, t := range tasks {
		if err := m.client.DeleteTask(t.ID); err != nil {
			errs[fmt.Sprintf("delete task %d", t.ID)] = err.Error()
		}
	}

	for _, e := range events {
		if err := m.client.DeleteEvent(e.Date); err != nil {
			errs[fmt.Sprintf("delete event %d", e.Date)] = err.Error()
		}
	}

	if len(errs) > 0 {
		errMsgs := []string{}
		for k, v := range errs {
			errMsgs = append(errMsgs, fmt.Sprintf("%s: %s", k, v))
		}
		errsMsg := strings.Join(errMsgs, "\n  ")
		return fmt.Errorf("Encountered errors during reset:\n  %s", errsMsg)
	} else {
		return nil
	}
}
