// This package contains a task.Manager that stores its state in RAM (during use)
// and on the local file system via a output/context pair.
package local

import (
	"fmt"
	"sort"
	"time"

	"github.com/ankeesler/anwork/task"
)

type unknownTaskError struct {
	name string
}

func (e unknownTaskError) Error() string {
	return fmt.Sprintf("cannot find task with name %s", e.name)
}

type manager struct {
	MyTasks    []*task.Task  `json:"tasks"`
	MyEvents   []*task.Event `json:"events"`
	NextTaskID int
}

func newManager() *manager {
	return &manager{
		MyTasks:    []*task.Task{},
		MyEvents:   []*task.Event{},
		NextTaskID: 0,
	}
}

func (m *manager) Create(name string) error {
	if m.FindByName(name) != nil {
		return fmt.Errorf("task '%s' has already been created", name)
	}

	t := &task.Task{
		Name:      name,
		ID:        m.NextTaskID,
		StartDate: time.Now().Unix(),
		Priority:  task.DefaultPriority,
		State:     task.StateWaiting,
	}
	m.NextTaskID++
	m.MyTasks = append(m.MyTasks, t)
	m.addEvent(fmt.Sprintf("Created task '%s'", name), task.EventTypeCreate, t.ID)
	return nil
}

func (m *manager) Delete(name string) error {
	index := -1
	for i, task := range m.MyTasks {
		if task.Name == name {
			index = i
			break
		}
	}

	if index != -1 {
		id := m.MyTasks[index].ID
		m.MyTasks = append(m.MyTasks[:index], m.MyTasks[index+1:]...)
		m.addEvent(fmt.Sprintf("Deleted task '%s'", name), task.EventTypeDelete, id)
		return nil
	} else {
		return unknownTaskError{name: name}
	}
}

func (m *manager) FindByID(id int) *task.Task {
	for _, t := range m.MyTasks {
		if t.ID == id {
			return t
		}
	}
	return nil
}

func (m *manager) FindByName(name string) *task.Task {
	for _, t := range m.MyTasks {
		if t.Name == name {
			return t
		}
	}
	return nil
}

func (m *manager) Tasks() []*task.Task {
	if !sort.IsSorted(m) {
		sort.Sort(m)
	}
	return m.MyTasks
}

func (m *manager) Note(name, note string) error {
	t := m.FindByName(name)
	if t == nil {
		return unknownTaskError{name: name}
	}
	m.addEvent(fmt.Sprintf("Note added to task '%s': %s", name, note), task.EventTypeNote, t.ID)
	return nil
}

func (m *manager) SetPriority(name string, newPriority int) error {
	t := m.FindByName(name)
	if t == nil {
		return unknownTaskError{name: name}
	}

	oldPriority := t.Priority
	t.Priority = newPriority
	m.addEvent(fmt.Sprintf("Set priority on task '%s' from %d to %d", name,
		oldPriority, newPriority),
		task.EventTypeSetPriority, t.ID)

	return nil
}

func (m *manager) SetState(name string, newState task.State) error {
	t := m.FindByName(name)
	if t == nil {
		return unknownTaskError{name: name}
	}

	oldState := t.State
	t.State = newState
	m.addEvent(fmt.Sprintf("Set state on task '%s' from %s to %s", name,
		task.StateNames[oldState], task.StateNames[newState]),
		task.EventTypeSetState, t.ID)

	return nil
}

func (m *manager) addEvent(title string, teyep task.EventType, taskID int) {
	m.MyEvents = append(m.MyEvents, &task.Event{
		Title:  title,
		Date:   time.Now().Unix(),
		Type:   teyep,
		TaskID: taskID,
	})
}

func (m *manager) Events() []*task.Event {
	return m.MyEvents
}

func (m *manager) DeleteEvent(startTime int64) error {
	foundItIndex := -1
	for i, event := range m.MyEvents {
		if event.Date == startTime {
			foundItIndex = i
			break
		}
	}

	if foundItIndex == -1 {
		return fmt.Errorf("No known event with start time %d", startTime)
	}

	m.MyEvents = append(m.MyEvents[:foundItIndex], m.MyEvents[foundItIndex+1:]...)

	return nil
}

func (m *manager) Reset() error {
	m.MyTasks = []*task.Task{}
	m.MyEvents = []*task.Event{}
	m.NextTaskID = 0
	return nil
}

// Return the length of the Task's held by this Manager.
func (m *manager) Len() int {
	return len(m.MyTasks)
}

// Return true iff the i'th Task held by this Manager is more "important" than the j'th Task held by
// this Manager. See the documentation for Manager. Tasks for more discussion around this design.
func (m *manager) Less(i, j int) bool {
	ti, tj := m.MyTasks[i], m.MyTasks[j]
	if ti.Priority > tj.Priority {
		return false
	} else if ti.Priority == tj.Priority {
		return ti.ID < tj.ID
	} else { // ti.Priority < tj.Priority
		return true
	}
}

// Swap the i'th Task held by this Manager with the j'th Task held by this Manager.
func (m *manager) Swap(i, j int) {
	tmp := m.MyTasks[i]
	m.MyTasks[i] = m.MyTasks[j]
	m.MyTasks[j] = tmp
}
