package task

import (
	"fmt"
	"sort"
	"time"
)

// A Manager is an interface through which Task's can be created, read, updated, and deleted.
type Manager struct {
	tasks   []*Task
	journal *Journal
}

// Create a new manager with an empty list of Task's.
func NewManager() *Manager {
	return &Manager{journal: &Journal{}}
}

// Create a Task with the provided name. This function will panic if a Task with the provided name
// already exists.
func (m *Manager) Create(name string) {
	t := m.Find(name)
	if t != nil {
		msg := fmt.Sprintf("Task with name %s already exists", name)
		panic(msg)
	}

	t = newTask(name)
	m.tasks = append(m.tasks, t)

	title := fmt.Sprintf("Created task %s", name)
	m.journal.Events = append(m.journal.Events, &Event{Title: title, T: time.Now()})
}

// Delete a Task with the provided name. Returns true iff the deletion was successful.
func (m *Manager) Delete(name string) bool {
	taskIndex := -1
	for index, task := range m.tasks {
		if task.name == name {
			taskIndex = index
			break
		}
	}

	if taskIndex != -1 {
		m.tasks = append(m.tasks[:taskIndex], m.tasks[taskIndex+1:]...)

		title := fmt.Sprintf("Deleted task %s", name)
		m.journal.Events = append(m.journal.Events, &Event{Title: title, T: time.Now()})

		return true
	} else {
		return false
	}
}

// Set the State of a Task currently in this Manager. This function will panic if there is no known
// Task with the provided name. The Task will be searched for via a call to Manager.Find(name).
func (m *Manager) SetState(name string, state State) {
	t := m.findOrPanic(name)
	beforeState := t.state
	t.state = state

	title := fmt.Sprintf("Set state on task %s from %s to %s",
		name, StateNames[beforeState], StateNames[state])
	m.journal.Events = append(m.journal.Events, &Event{Title: title, T: time.Now()})
}

// Set the priority of a Task currently in this Manager. This function will panic if there is no
// known Task with the provided name. The Task will be searched for via a call to Manager.Find(name).
func (m *Manager) SetPriority(name string, priority int32) {
	t := m.findOrPanic(name)
	beforePriority := t.priority
	t.priority = priority

	title := fmt.Sprintf("Set priority on task %s from %d to %d", name, beforePriority, priority)
	m.journal.Events = append(m.journal.Events, &Event{Title: title, T: time.Now()})
}

// Get all of the Tasks contained in this manager, ordered from highest priority (lowest integer
// value) to lowest priority (highest integer value).
//
// When multiple tasks have the same priority, the Task's will be ordered by their (unique) ID in
// ascending order. This means that the older Task's will come first. This is a conscious decision.
// The Task's that have been around the longest are assumed to need to be completed first.
func (m *Manager) Tasks() []*Task {
	if !sort.IsSorted(m) {
		sort.Sort(m)
	}
	return m.tasks
}

// Find a Task with the provided name in this Manager, or return nil if there is no such Task.
func (m *Manager) Find(name string) *Task {
	for _, task := range m.tasks {
		if task.name == name {
			return task
		}
	}
	return nil
}

func (m *Manager) findOrPanic(name string) *Task {
	t := m.Find(name)
	if t == nil {
		msg := fmt.Sprintf("Unknown task with name %s", name)
		panic(msg)
	}
	return t
}

// Return the length of the Task's held by this Manager.
func (m *Manager) Len() int {
	return len(m.tasks)
}

// Return true iff the i'th Task held by this Manager is more "important" than the j'th Task held by
// this Manager. See the documentation for Manager. Tasks for more discussion around this design.
func (m *Manager) Less(i, j int) bool {
	ti, tj := m.tasks[i], m.tasks[j]
	if ti.priority > tj.priority {
		return false
	} else if ti.priority == tj.priority {
		return ti.id < tj.id
	} else { // ti.priority < tj.priority
		return true
	}
}

// Swap the i'th Task held by this Manager with the j'th Task held by this Manager.
func (m *Manager) Swap(i, j int) {
	tmp := m.tasks[i]
	m.tasks[i] = m.tasks[j]
	m.tasks[j] = tmp
}

func (m *Manager) Journal() *Journal {
	return m.journal
}
