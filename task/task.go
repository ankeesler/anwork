// This package contains all of the task-related data and functionality in the anwork project.
//
// A Task is something that someone is working on. It could be something like "mow the lawn" or "buy
// sister a holiday present."
//
// Every Task is in one of a number of different State's: Waiting, Blocked, Running, or Finished. A
// Task also has a priority which describes its relative importance to all other Task's.
//
// A Manager is an interface through which Task's can be created, read, updated, and deleted.
//
// A Manager also keeps track of the changes that are made to the Task's it oversees via a Journal.
// A Journal is simply a list of things (Event's) that happen to a Manager (i.e., a note is added, a
// Task is created, a Task is updated, etc.).
package task

import (
	"encoding/json"
	"time"

	pb "github.com/ankeesler/anwork/task/proto"
	"github.com/golang/protobuf/proto"
)

// A State describes the status of some Task.
type State int

// These are the states that a Task could be in.
const (
	StateWaiting  = State(0)
	StateBlocked  = State(1)
	StateRunning  = State(2)
	StateFinished = State(3)
)

// These are the names of the State's that a Task can occupy, indexed by the State integer value.
var StateNames = [...]string{
	"Waiting",
	"Blocked",
	"Running",
	"Finished",
}

// This is the default priority that a Task gets when created.
const DefaultPriority = 10

var nextTaskID int = 0

// A Task is something that someone is working on. It could be something like "mow the lawn" or "buy
// sister a holiday present." A Task also has a priority which describes its relative importance to
// all other Task's.
type Task struct {
	// The name of the Task, i.e., "mow the lawn" or "PROJECT-123-fix-infinite-recursion."
	Name string `json:"name"`

	// This is a unique ID. Every Task has a different ID.
	ID int `json:"id"`

	// This is a description of the Task.
	Description string `json:"description"`

	// This is when the Task was created.
	StartDate time.Time `json:"startDate"`

	// This is the priority of the Task. The lower the number, the higher the importance.
	Priority int `json:"priority"`

	// This is the State of the Task. See State* for possible values. A Task can go through any
	// number of State changes over the course of its life. All Tasks start out in the StateWaiting
	// State.
	State State `json:"state"`
}

func (t *Task) Serialize() ([]byte, error) {
	return json.Marshal(t)
}

func (t *Task) Unserialize(bytes []byte) error {
	tProtobuf := pb.Task{}
	err := proto.Unmarshal(bytes, &tProtobuf)
	if err == nil {
		t.fromProtobuf(&tProtobuf)
		return nil
	}

	return json.Unmarshal(bytes, t)
}

func (t *Task) fromProtobuf(tProtobuf *pb.Task) {
	t.Name = tProtobuf.Name
	t.ID = int(tProtobuf.ID)
	t.Description = tProtobuf.Description
	t.StartDate = time.Unix(tProtobuf.StartDate, 0) // sec, nsec
	t.Priority = int(tProtobuf.Priority)
	t.State = State(tProtobuf.State)

	noteTaskID(t.ID)
}

// Create a new Task with a default priority (see DefaultPriority) in the waiting state (see
// StateWaiting).
func NewTask(name string) *Task {
	t := &Task{
		Name:      name,
		ID:        nextTaskID,
		StartDate: time.Now(),
		Priority:  DefaultPriority,
		State:     StateWaiting,
	}

	nextTaskID++

	// Truncate the start time at the seconds since we only persist the seconds amount.
	t.StartDate = t.StartDate.Truncate(time.Second)

	return t
}

// This function should be called when a task ID is unpersisted from disk. This will make sure that
// the nextTaskID always points to a unique task ID.
func noteTaskID(id int) {
	// Increment the ID to one higher than what we just read in to make sure everyone is getting a
	// unique ID.
	if id >= nextTaskID {
		nextTaskID = id + 1
	}
}
