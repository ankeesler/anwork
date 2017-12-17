package task

import (
	"fmt"
	"strings"
	"testing"

	"github.com/ankeesler/anwork/storage"
	pb "github.com/ankeesler/anwork/task/proto"
)

const (
	root        = "test-data"
	tmpContext  = "tmp-context"
	goodContext = "good-context"
	badContext  = "bad-context"

	goodTaskName        = "task-a"
	goodTaskId          = 0
	goodTaskDescription = "Here is a description!"
	goodTaskPriority    = 612
	goodTaskState       = TaskStateRunning

	taskWithId5Context = "task-with-id-5-context"
	taskWithId5Name    = "task-with-id-5"
	taskWithId5Id      = 5
)

func expectTasksEqual(t *testing.T, t0, t1 *Task, expectTime bool) {
	if !strings.EqualFold(t0.name, t1.name) {
		t.Fatalf("Task names do not match: %s vs %s", t0.name, t1.name)
	} else if t0.id != t1.id {
		t.Fatalf("Task ids do not match: %d vs %d", t0.id, t1.id)
	} else if !strings.EqualFold(t0.description, t1.description) {
		t.Fatalf("Task descriptions do not match: %s vs %s", t0.description, t1.description)
	} else if expectTime && !t0.startDate.Equal(t1.startDate) {
		t.Fatalf("Task start date do not match: %s vs %s", t0.startDate, t1.startDate)
	} else if t0.priority != t1.priority {
		t.Fatalf("Task priorities do not match: %d vs %d", t0.priority, t1.priority)
	} else if t0.state != t1.state {
		t.Fatalf("Task states do not match: %d vs %d", t0.state, t1.state)
	}
}

func TestTaskStateProtobuf(t *testing.T) {
	statePairs := [][]int32{
		{TaskStateWaiting, int32(pb.TaskStateProtobuf_WAITING)},
		{TaskStateBlocked, int32(pb.TaskStateProtobuf_BLOCKED)},
		{TaskStateRunning, int32(pb.TaskStateProtobuf_RUNNING)},
		{TaskStateFinished, int32(pb.TaskStateProtobuf_FINISHED)},
	}
	for i, statePair := range statePairs {
		if statePair[0] != statePair[1] {
			t.Errorf("TaskState at index %d (%d) does not match TaskStateProtobuf %d",
				i, statePair[0], statePair[1])
		}
	}
}

func TestPersist(t *testing.T) {
	persister := storage.Persister{root}
	if persister.Exists(tmpContext) {
		t.Fatalf("Cannot run this test when context (%s) already exists", tmpContext)
	}
	defer persister.Delete(tmpContext)

	task := newTask("this is my task")
	task.description = "Yeah yeah yeah"
	task.priority = 21
	task.state = TaskStateBlocked
	err := persister.Persist(tmpContext, task)
	if err != nil {
		t.Fatalf("Failed to persist task (%v) to file: %s", task, err)
	}

	unpersistedTask := &Task{}
	err = persister.Unpersist(tmpContext, unpersistedTask)
	if err != nil {
		t.Fatalf("Failed to unpersist task (%v) from file: %s", task, err)
	}

	expectTasksEqual(t, task, unpersistedTask, true) // expectTime
}

func TestUnpersist(t *testing.T) {
	persister := storage.Persister{root}
	if !persister.Exists(goodContext) {
		t.Fatalf("Cannot run this test when context (%s) does not exist", tmpContext)
	}

	unpersistedTask := Task{}
	err := persister.Unpersist(goodContext, &unpersistedTask)
	if err != nil {
		t.Fatalf("Could not unpersist task from context %s: %s", goodContext, err)
	}

	expectedTask := Task{
		name:        goodTaskName,
		id:          goodTaskId,
		description: goodTaskDescription,
		priority:    goodTaskPriority,
		state:       goodTaskState,
	}
	expectTasksEqual(t, &unpersistedTask, &expectedTask, false) // expectTime
}

func TestTaskIdUniqueness(t *testing.T) {
	persister := storage.Persister{root}
	if !persister.Exists(taskWithId5Context) {
		t.Fatalf("Cannot run this test when context (%s) does not exist", taskWithId5Context)
	}

	taskWithId5 := Task{}
	persister.Unpersist(taskWithId5Context, &taskWithId5)
	if taskWithId5.id != taskWithId5Id {
		t.Errorf("Expected id %d from task %#v but got %d",
			taskWithId5Id, taskWithId5, taskWithId5.id)
	} else if !strings.EqualFold(taskWithId5Name, taskWithId5.name) {
		t.Errorf("Expected name %s from task %s but got %s",
			taskWithId5Name, taskWithId5Id, taskWithId5.name)
	}

	for i := 0; i < 10; i++ {
		name := fmt.Sprintf("task-%d", i)
		task := newTask(name)
		if task.id == taskWithId5Id {
			t.Errorf("New task %s has non-unique id %d", name, task.id)
		}
	}
}

func TestBadContext(t *testing.T) {
	persister := storage.Persister{root}
	if !persister.Exists(badContext) {
		t.Fatalf("Cannot run this test when context (%s) does not exist", badContext)
	}

	task := Task{}
	err := persister.Unpersist(badContext, &task)
	if err == nil {
		t.Fatalf("Expected error from load from bad context (%s) but didn't get one", badContext)
	}
}
