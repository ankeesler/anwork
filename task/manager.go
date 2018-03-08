package task

import (
	"bytes"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/golang/protobuf/proto"

	pb "github.com/ankeesler/anwork/task/proto"
)

// A Manager is an interface through which Task's can be created, read, updated, and deleted.
type Manager struct {
	tasks   []*Task
	journal *Journal
}

// Create a new manager with an empty list of Task's.
func NewManager() *Manager {
	return &Manager{tasks: make([]*Task, 0), journal: NewJournal()}
}

// Create a Task with the provided name. This function will panic if a Task with the provided name
// already exists.
func (m *Manager) Create(name string) {
	t := m.FindByName(name)
	if t != nil {
		msg := fmt.Sprintf("Task with name %s already exists", name)
		panic(msg)
	}

	t = NewTask(name)
	m.tasks = append(m.tasks, t)

	title := fmt.Sprintf("Created task %s", name)
	m.journal.Events = append(m.journal.Events, newEvent(title, EventTypeCreate, t.ID))
}

// Delete a Task with the provided name. Returns true iff the deletion was successful.
func (m *Manager) Delete(name string) bool {
	taskIndex := -1
	for index, task := range m.tasks {
		if task.Name == name {
			taskIndex = index
			break
		}
	}

	if taskIndex != -1 {
		t := m.tasks[taskIndex]
		m.tasks = append(m.tasks[:taskIndex], m.tasks[taskIndex+1:]...)

		title := fmt.Sprintf("Deleted task %s", name)
		m.journal.Events = append(m.journal.Events, newEvent(title, EventTypeDelete, t.ID))

		return true
	} else {
		return false
	}
}

// Set the State of a Task currently in this Manager. This function will panic if there is no known
// Task with the provided name. The Task will be searched for via a call to Manager.FindByName(name).
func (m *Manager) SetState(name string, state State) {
	t := m.findOrPanic(name)
	beforeState := t.State
	t.State = state

	title := fmt.Sprintf("Set state on task %s from %s to %s",
		name, StateNames[beforeState], StateNames[state])
	m.journal.Events = append(m.journal.Events, newEvent(title, EventTypeSetState, t.ID))
}

// Set the priority of a Task currently in this Manager. This function will panic if there is no
// known Task with the provided name. The Task will be searched for via a call to
// Manager.FindByName(name).
func (m *Manager) SetPriority(name string, priority int) {
	t := m.findOrPanic(name)
	beforePriority := t.Priority
	t.Priority = priority

	title := fmt.Sprintf("Set priority on task %s from %d to %d", name, beforePriority, priority)
	m.journal.Events = append(m.journal.Events, newEvent(title, EventTypeSetPriority, t.ID))
}

// Add a note that relates to a Task. This function will panic if there is no known Task with the
// provided name. The Task will be searched for via a call to Manager.FindByName(name).
func (m *Manager) Note(name string, note string) {
	t := m.findOrPanic(name) // this ensures that a Task exists for the provided name.
	title := fmt.Sprintf("Note added to task %s: %s", name, note)
	m.journal.Events = append(m.journal.Events, newEvent(title, EventTypeNote, t.ID))
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
func (m *Manager) FindByName(name string) *Task {
	for _, task := range m.tasks {
		if task.Name == name {
			return task
		}
	}
	return nil
}

// Find a Task with the provided id in this Manager, or return nil if there is no such Task.
func (m *Manager) FindByID(id int) *Task {
	for _, task := range m.tasks {
		if task.ID == id {
			return task
		}
	}
	return nil
}

func (m *Manager) findOrPanic(name string) *Task {
	t := m.FindByName(name)
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
	if ti.Priority > tj.Priority {
		return false
	} else if ti.Priority == tj.Priority {
		return ti.ID < tj.ID
	} else { // ti.Priority < tj.Priority
		return true
	}
}

// Swap the i'th Task held by this Manager with the j'th Task held by this Manager.
func (m *Manager) Swap(i, j int) {
	tmp := m.tasks[i]
	m.tasks[i] = m.tasks[j]
	m.tasks[j] = tmp
}

// Get the Journal associated with this Manager.
func (m *Manager) Journal() *Journal {
	return m.journal
}

func (m *Manager) Serialize() ([]byte, error) {
	return json.Marshal(m)
}

func (m *Manager) Unserialize(bytes []byte) error {
	mProtobuf := pb.Manager{}
	err := proto.Unmarshal(bytes, &mProtobuf)
	if err == nil {
		tsProtobuf := mProtobuf.GetTasks()
		m.tasks = make([]*Task, len(tsProtobuf))
		for index, tProtobuf := range tsProtobuf {
			m.tasks[index] = &Task{}
			m.tasks[index].fromProtobuf(tProtobuf)
		}

		m.journal.fromProtobuf(mProtobuf.Journal)
		return nil
	}

	return json.Unmarshal(bytes, m)
}

func (m *Manager) String() string {
	var buf bytes.Buffer
	buf.WriteString("manager{")

	buf.WriteString("tasks:")
	for _, t := range m.tasks {
		str := fmt.Sprintf("%s(%d),", t.Name, t.ID)
		buf.WriteString(str)
	}

	buf.WriteString(";journal:")
	for _, e := range m.journal.Events {
		str := fmt.Sprintf("'%s',", e.Title)
		buf.WriteString(str)
	}

	buf.WriteString(";}")
	return buf.String()
}

var jsonManager struct {
	Tasks   []*Task  `json:"tasks"`
	Journal *Journal `json:"journal"`
}

func (m *Manager) MarshalJSON() ([]byte, error) {
	jsonManager.Tasks = m.tasks
	jsonManager.Journal = m.journal
	return json.Marshal(jsonManager)
}

func (m *Manager) UnmarshalJSON(bytes []byte) error {
	if err := json.Unmarshal(bytes, &jsonManager); err != nil {
		return err
	}

	m.tasks = jsonManager.Tasks
	m.journal = jsonManager.Journal

	for _, task := range m.tasks {
		noteTaskID(task.ID)
	}
	for _, event := range m.journal.Events {
		noteTaskID(event.TaskID)
	}

	return nil
}
