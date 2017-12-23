package task

import (
	"fmt"

	"github.com/ankeesler/anwork/storage"
	pb "github.com/ankeesler/anwork/task/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
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

var _ = Describe("TaskState constants", func() {
	It("should line up with Protobuf definitions", func() {
		statePairs := [][]int32{
			{TaskStateWaiting, int32(pb.TaskStateProtobuf_WAITING)},
			{TaskStateBlocked, int32(pb.TaskStateProtobuf_BLOCKED)},
			{TaskStateRunning, int32(pb.TaskStateProtobuf_RUNNING)},
			{TaskStateFinished, int32(pb.TaskStateProtobuf_FINISHED)},
		}
		for _, statePair := range statePairs {
			Expect(statePair[0]).To(Equal(statePair[1]))
		}
	})
})
var _ = Describe("Task's", func() {
	It("are persistable", func() {

		persister := storage.Persister{root}
		Expect(persister.Exists(tmpContext)).To(BeFalse(),
			"Cannot run this test when context (%s) already exists", tmpContext)
		defer persister.Delete(tmpContext)

		task := newTask("this is my task")
		task.description = "Yeah yeah yeah"
		task.priority = 21
		task.state = TaskStateBlocked
		err := persister.Persist(tmpContext, task)
		Expect(err).ToNot(HaveOccurred(), "Failed to persist task (%v) to file: %s", task, err)

		unpersistedTask := &Task{}
		err = persister.Unpersist(tmpContext, unpersistedTask)
		Expect(err).ToNot(HaveOccurred(), "Failed to unpersist task (%v) from file: %s", task, err)

		Expect(unpersistedTask).To(Equal(task))
	})
	It("are unpersistable", func() {

		persister := storage.Persister{root}
		Expect(persister.Exists(goodContext)).To(BeTrue(),
			"Cannot run this test when context (%s) does not exist", tmpContext)

		unpersistedTask := Task{}
		err := persister.Unpersist(goodContext, &unpersistedTask)
		Expect(err).ToNot(HaveOccurred(),
			"Could not unpersist task from context %s: %s", goodContext, err)

		expectedTask := Task{
			name:        goodTaskName,
			id:          goodTaskId,
			description: goodTaskDescription,
			priority:    goodTaskPriority,
			state:       goodTaskState,

			startDate: unpersistedTask.startDate,
		}
		Expect(unpersistedTask).To(Equal(expectedTask))
	})
	Context("have unique ID's", func() {
		It("that are larger", func() {

			persister := storage.Persister{root}
			Expect(persister.Exists(taskWithId5Context)).To(BeTrue(),
				"Cannot run this test when context (%s) does not exist", taskWithId5Context)

			taskWithId5 := Task{}
			persister.Unpersist(taskWithId5Context, &taskWithId5)
			Expect(taskWithId5.id).To(BeEquivalentTo(taskWithId5Id))
			Expect(taskWithId5.name).To(Equal(taskWithId5Name))

			for i := 0; i < 10; i++ {
				name := fmt.Sprintf("task-%d", i)
				task := newTask(name)
				Expect(task.id).ToNot(Equal(taskWithId5Id))
			}
		})
		It("that are smaller", func() {

			persister := storage.Persister{root}
			Expect(persister.Exists(tmpContext)).To(BeFalse(),
				"Cannot run this test when context (%s) exists", tmpContext)
			defer persister.Delete(tmpContext)

			task := Task{}
			task.id = 0
			err := persister.Persist(tmpContext, &task)
			Expect(err).ToNot(HaveOccurred(),
				"Got error when persisting task %#v to file: %s", task, err)

			ids := make(map[int32]bool)
			for i := 0; i < 25; i++ {
				name := fmt.Sprintf("task-%d", i)
				task := newTask(name)
				ids[task.id] = true
			}

			err = persister.Unpersist(tmpContext, &task)
			Expect(err).ToNot(HaveOccurred(),
				"Got error when unpersisting task %#v from file: %s", task, err)
			Expect(task.id).To(BeZero(),
				"Expected task id to be %d but was %d", 0, task.id)

			for i := 0; i < 25; i++ {
				name := fmt.Sprintf("task-%d", i)
				task := newTask(name)
				_, ok := ids[task.id]
				Expect(ok).To(BeFalse(), "Failure! Task ID %d already exists!", task.id)
			}
		})
	})
	It("fail gracefully from bad contexts", func() {

		persister := storage.Persister{root}
		Expect(persister.Exists(badContext)).To(BeTrue(),
			"Cannot run this test when context (%s) does not exist", badContext)

		task := Task{}
		err := persister.Unpersist(badContext, &task)
		Expect(err).To(HaveOccurred(),
			"Expected error from load from bad context (%s) but didn't get one", badContext)
	})
})
