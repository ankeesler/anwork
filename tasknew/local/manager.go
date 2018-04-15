package local

import (
	"fmt"
	"sort"
	"time"

	task "github.com/ankeesler/anwork/tasknew"
)

type manager struct {
	tasks  []*task.Task
	events []*task.Event
}

var nextTaskID = 0

func NewManager() task.Manager {
	return &manager{}
}

func (m *manager) Create(name string) error {
	if m.FindByName(name) != nil {
		return fmt.Errorf("task '%s' has already been created", name)
	}

	t := &task.Task{
		Name:      name,
		ID:        nextTaskID,
		StartDate: time.Now().Unix(),
		Priority:  task.DefaultPriority,
		State:     task.StateWaiting,
	}
	nextTaskID++
	m.tasks = append(m.tasks, t)
	m.addEvent(fmt.Sprintf("created task '%s'", name), task.EventTypeCreate, t.ID)
	return nil
}

func (m *manager) Delete(name string) bool {
	index := -1
	for i, task := range m.tasks {
		if task.Name == name {
			index = i
			break
		}
	}

	if index != -1 {
		id := m.tasks[index].ID
		m.tasks = append(m.tasks[:index], m.tasks[index+1:]...)
		m.addEvent(fmt.Sprintf("deleted task '%s'", name), task.EventTypeDelete, id)
		return true
	} else {
		return false
	}
}

func (m *manager) FindByID(id int) *task.Task {
	for _, t := range m.tasks {
		if t.ID == id {
			return t
		}
	}
	return nil
}

func (m *manager) FindByName(name string) *task.Task {
	for _, t := range m.tasks {
		if t.Name == name {
			return t
		}
	}
	return nil
}

func (m *manager) mustFindByName(name string) *task.Task {
	t := m.FindByName(name)
	if t == nil {
		panic(fmt.Sprintf("cannot find task with name %s", name))
	}
	return t
}

func (m *manager) Tasks() []*task.Task {
	if !sort.IsSorted(m) {
		sort.Sort(m)
	}
	return m.tasks
}

func (m *manager) Note(name, note string) {
	t := m.mustFindByName(name)
	m.addEvent(fmt.Sprintf("note added to task '%s': %s", name, note), task.EventTypeNote, t.ID)
}

func (m *manager) SetPriority(name string, newPriority int) {
	t := m.mustFindByName(name)
	oldPriority := t.Priority
	t.Priority = newPriority
	m.addEvent(fmt.Sprintf("set priority on task '%s' from %d to %d", name,
		oldPriority, newPriority),
		task.EventTypeSetPriority, t.ID)
}

func (m *manager) SetState(name string, newState task.State) {
	t := m.mustFindByName(name)
	oldState := t.State
	t.State = newState
	m.addEvent(fmt.Sprintf("set state on task '%s' from %s to %s", name,
		task.StateNames[oldState], task.StateNames[newState]),
		task.EventTypeSetState, t.ID)
}

func (m *manager) addEvent(title string, teyep task.EventType, taskID int) {
	m.events = append(m.events, &task.Event{
		Title:  title,
		Date:   time.Now().Unix(),
		Type:   teyep,
		TaskID: taskID,
	})
}

func (m *manager) Events() []*task.Event {
	return m.events
}

// Return the length of the Task's held by this Manager.
func (m *manager) Len() int {
	return len(m.tasks)
}

// Return true iff the i'th Task held by this Manager is more "important" than the j'th Task held by
// this Manager. See the documentation for Manager. Tasks for more discussion around this design.
func (m *manager) Less(i, j int) bool {
	ti, tj := m.tasks[i], m.tasks[j]
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
	tmp := m.tasks[i]
	m.tasks[i] = m.tasks[j]
	m.tasks[j] = tmp
}
