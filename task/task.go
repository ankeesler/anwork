// This package contains all of the task-related data and functionality in the anwork project.
//
// A Task is something that someone is working on. It could be something like "mow the lawn" or "buy
// sister a holiday present."
//
// Every Task is in one of a number of different State's: Waiting, Blocked, Running, or Finished.
//
// A Manager is an interface through which Task's can be created, read, updated, and deleted.
package task

import (
	"time"

	pb "github.com/ankeesler/anwork/task/proto"
	"github.com/golang/protobuf/proto"
)

// A State describes the status of some Task.
type State uint32

// These are the states that a Task could be in.
const (
	StateWaiting  = State(0)
	StateBlocked  = State(1)
	StateRunning  = State(2)
	StateFinished = State(3)
)

// This is the default priority that a Task gets when created.
const DefaultPriority = 10

var nextTaskId int32 = 0

//go:generate protoc --proto_path=proto --go_out=proto task.proto

// A Task is something that someone is working on. It could be something like "mow the lawn" or "buy
// sister a holiday present."
type Task struct {
	// The name of the Task, i.e., "mow the lawn" or "PROJECT-123-fix-infinite-recursion."
	name string

	// This is a unique ID. Every Task has a different ID.
	id int32

	// This is a description of the Task.
	description string

	// This is when the Task was created.
	startDate time.Time

	// This is the priority of the Task. The lower the number, the higher the importance.
	priority int32

	// This is the State of the Task. See State* for possible values. A Task can go through any
	// number of State changes over the course of its life. All Tasks start out in the StateWaiting
	// State.
	state State
}

func (t *Task) Serialize() ([]byte, error) {
	tProtobuf := pb.TaskProtobuf{
		Name:        t.name,
		Id:          t.id,
		Description: t.description,
		StartDate:   t.startDate.Unix(),
		Priority:    t.priority,
		State:       pb.StateProtobuf(t.state),
	}
	return proto.Marshal(&tProtobuf)
}

func (t *Task) Unserialize(bytes []byte) error {
	tProtobuf := pb.TaskProtobuf{}
	err := proto.Unmarshal(bytes, &tProtobuf)
	if err != nil {
		return err
	}

	t.name = tProtobuf.Name
	t.id = tProtobuf.Id
	t.description = tProtobuf.Description
	t.startDate = time.Unix(tProtobuf.StartDate, 0) // sec, nsec
	t.priority = tProtobuf.Priority
	t.state = State(tProtobuf.State)

	// Increment the ID to one higher than what we just read in to make sure everyone is getting a
	// unique ID.
	if t.id >= nextTaskId {
		nextTaskId = t.id + 1
	}

	return nil
}

// Get the State for this task.
func (t *Task) State() State {
	return t.state
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
