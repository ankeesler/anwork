// This package contains all of the task-related data and functionality in the anwork project.
//
// A Task is something that someone is working on. It could be something like "mow the lawn" or "buy
// sister a holiday present."
//
// Every Task is in one of a number of different State's: Waiting, Blocked, Running, or Finished.
//
// A Manager is an interface through which Task's can be created, read, updated, and deleted.
//
// A Manager also keeps track of the changes that are made to the Task's it oversees via a Journal.
// A Journal is simply a list of things (Event's) that happen to a Manager (i.e., a note is added, a
// Task is created, a Task is updated, etc.).
package task

import (
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

var nextTaskId int = 0

//go:generate protoc --proto_path=proto --go_out=proto task.proto

// A Task is something that someone is working on. It could be something like "mow the lawn" or "buy
// sister a holiday present."
type Task struct {
	// The name of the Task, i.e., "mow the lawn" or "PROJECT-123-fix-infinite-recursion."
	name string

	// This is a unique ID. Every Task has a different ID.
	id int

	// This is a description of the Task.
	// TODO: do we need this? We don't show it anywhere on the screen to the user.
	description string

	// This is when the Task was created.
	startDate time.Time

	// This is the priority of the Task. The lower the number, the higher the importance.
	priority int

	// This is the State of the Task. See State* for possible values. A Task can go through any
	// number of State changes over the course of its life. All Tasks start out in the StateWaiting
	// State.
	state State
}

func (t *Task) Serialize() ([]byte, error) {
	var tProtobuf pb.Task
	t.toProtobuf(&tProtobuf)
	return proto.Marshal(&tProtobuf)
}

func (t *Task) toProtobuf(tProtobuf *pb.Task) {
	tProtobuf.Name = t.name
	tProtobuf.Id = int32(t.id)
	tProtobuf.Description = t.description
	tProtobuf.StartDate = t.startDate.Unix()
	tProtobuf.Priority = int32(t.priority)
	tProtobuf.State = pb.State(t.state)
}

func (t *Task) Unserialize(bytes []byte) error {
	tProtobuf := pb.Task{}
	err := proto.Unmarshal(bytes, &tProtobuf)
	if err != nil {
		return err
	}

	t.fromProtobuf(&tProtobuf)

	return nil
}

func (t *Task) fromProtobuf(tProtobuf *pb.Task) {
	t.name = tProtobuf.Name
	t.id = int(tProtobuf.Id)
	t.description = tProtobuf.Description
	t.startDate = time.Unix(tProtobuf.StartDate, 0) // sec, nsec
	t.priority = int(tProtobuf.Priority)
	t.state = State(tProtobuf.State)

	// Increment the ID to one higher than what we just read in to make sure everyone is getting a
	// unique ID.
	if t.id >= nextTaskId {
		nextTaskId = t.id + 1
	}
}

// Get the name of this Task.
func (t *Task) Name() string {
	return t.name
}

// Get the State for this Task.
func (t *Task) State() State {
	return t.state
}

// Get the ID for this Task.
func (t *Task) ID() int {
	return t.id
}

// Create a new Task with a default priority (see DefaultPriority) in the waiting state (see
// StateWaiting).
func newTask(name string) *Task {
	t := &Task{
		name:      name,
		id:        nextTaskId,
		startDate: time.Now(),
		priority:  DefaultPriority,
		state:     StateWaiting,
	}

	nextTaskId++

	// Truncate the start time at the seconds since we only persist the seconds amount.
	t.startDate = t.startDate.Truncate(time.Second)

	return t
}
