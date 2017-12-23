package task

import (
	"sort"
)

// A Manager is an interface through which Task's can be created, read, updated, and deleted.
type Manager struct {
	tasks []*Task
}

// Create a new manager with an empty list of Task's.
func NewManager() *Manager {
	return &Manager{}
}

// Create a Task with the provided name.
func (m *Manager) Create(name string) {
	t := newTask(name)
	m.tasks = append(m.tasks, t)
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
		return true
	} else {
		return false
	}
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

// Return the length of the Task's held by this Manager.
func (m *Manager) Len() int {
	return len(m.tasks)
}

// Return true iff the i'th Task held by this Manager is more "important" than the j'th Task held by
// this Manager. See the documentation for Manager.Tasks for more discussion around this design.
func (m *Manager) Less(i, j int) bool {
	ti, tj := m.tasks[i], m.tasks[j]
	if ti.priority >= tj.priority {
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
