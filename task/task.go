// This package contains all of the task-related data and functionality in the anwork project.
//
// A Task is something that someone is working on. It could be something like "mow the lawn" or "buy
// sister a holiday present."
package task

import (
	"time"

	pb "github.com/ankeesler/anwork/task/proto"
	"github.com/golang/protobuf/proto"
)

// These are the states that a Task could be in.
const (
	TaskStateWaiting  = 0
	TaskStateBlocked  = 1
	TaskStateRunning  = 2
	TaskStateFinished = 3
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

	// This is the state of the Task. See TaskState* for possible values. A Task can go through any
	// number of state changes over the course of its life. All Tasks start out in the TaskStateWaiting
	// state.
	state int32
}

func (t *Task) Serialize() ([]byte, error) {
	tProtobuf := pb.TaskProtobuf{
		Name:        t.name,
		Id:          t.id,
		Description: t.description,
		StartDate:   t.startDate.Unix(),
		Priority:    t.priority,
		State:       pb.TaskStateProtobuf(t.state),
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
	t.state = int32(tProtobuf.State)

	// Increment the ID to one higher than what we just read in to make sure everyone is getting a
	// unique ID.
	if t.id >= nextTaskId {
		nextTaskId = t.id + 1
	}

	return nil
}

// Create a new Task with a default priority (see DefaultPriority) in the waiting state (see
// TaskStateWaiting).
func newTask(name string) *Task {
	t := &Task{
		name:      name,
		id:        nextTaskId,
		startDate: time.Now(),
		priority:  DefaultPriority,
		state:     TaskStateWaiting,
	}

	nextTaskId++

	// Truncate the start time at the seconds since we only persist the seconds amount.
	t.startDate = t.startDate.Truncate(time.Second)

	return t
}
