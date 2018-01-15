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
	goodTaskID          = 0
	goodTaskDescription = "Here is a description!"
	goodTaskPriority    = 612
	goodState           = StateRunning

	taskWithID5Context = "task-with-id-5-context"
	taskWithID5Name    = "task-with-id-5"
	taskWithID5ID      = 5
)

var _ = Describe("State constants", func() {
	It("should line up with Protobuf definitions", func() {
		statePairs := [][]State{
			{StateWaiting, State(pb.State_WAITING)},
			{StateBlocked, State(pb.State_BLOCKED)},
			{StateRunning, State(pb.State_RUNNING)},
			{StateFinished, State(pb.State_FINISHED)},
		}
		for _, statePair := range statePairs {
			Expect(statePair[0]).To(Equal(statePair[1]))
		}
	})
	It("should have the expected state names", func() {
		Expect(StateNames[StateWaiting]).To(Equal("Waiting"))
		Expect(StateNames[StateBlocked]).To(Equal("Blocked"))
		Expect(StateNames[StateRunning]).To(Equal("Running"))
		Expect(StateNames[StateFinished]).To(Equal("Finished"))
	})
})
var _ = Describe("Task's", func() {
	It("are persistable", func() {

		persister := storage.Persister{Root: root}
		Expect(persister.Exists(tmpContext)).To(BeFalse(),
			"Cannot run this test when context (%s) already exists", tmpContext)
		defer persister.Delete(tmpContext)

		task := newTask("this is my task")
		task.description = "Yeah yeah yeah"
		task.priority = 21
		task.state = StateBlocked
		err := persister.Persist(tmpContext, task)
		Expect(err).ToNot(HaveOccurred(), "Failed to persist task (%v) to file: %s", task, err)

		unpersistedTask := &Task{}
		err = persister.Unpersist(tmpContext, unpersistedTask)
		Expect(err).ToNot(HaveOccurred(), "Failed to unpersist task (%v) from file: %s", task, err)

		Expect(unpersistedTask).To(Equal(task))
	})
	It("are unpersistable", func() {

		persister := storage.Persister{Root: root}
		Expect(persister.Exists(goodContext)).To(BeTrue(),
			"Cannot run this test when context (%s) does not exist", tmpContext)

		unpersistedTask := Task{}
		err := persister.Unpersist(goodContext, &unpersistedTask)
		Expect(err).ToNot(HaveOccurred(),
			"Could not unpersist task from context %s: %s", goodContext, err)

		expectedTask := Task{
			name:        goodTaskName,
			id:          goodTaskID,
			description: goodTaskDescription,
			priority:    goodTaskPriority,
			state:       goodState,

			startDate: unpersistedTask.StartDate(),
		}
		Expect(unpersistedTask).To(Equal(expectedTask))
		Expect(unpersistedTask.State()).To(Equal(goodState))
	})
	Context("have unique ID's", func() {
		It("that are larger", func() {

			persister := storage.Persister{Root: root}
			Expect(persister.Exists(taskWithID5Context)).To(BeTrue(),
				"Cannot run this test when context (%s) does not exist", taskWithID5Context)

			taskWithID5 := Task{}
			persister.Unpersist(taskWithID5Context, &taskWithID5)
			Expect(taskWithID5.id).To(BeEquivalentTo(taskWithID5ID))
			Expect(taskWithID5.name).To(Equal(taskWithID5Name))

			for i := 0; i < 10; i++ {
				name := fmt.Sprintf("task-%d", i)
				task := newTask(name)
				Expect(task.id).ToNot(Equal(taskWithID5ID))
			}
		})
		It("that are smaller", func() {

			persister := storage.Persister{Root: root}
			Expect(persister.Exists(tmpContext)).To(BeFalse(),
				"Cannot run this test when context (%s) exists", tmpContext)
			defer persister.Delete(tmpContext)

			task := Task{}
			task.id = 0
			err := persister.Persist(tmpContext, &task)
			Expect(err).ToNot(HaveOccurred(),
				"Got error when persisting task %#v to file: %s", task, err)

			ids := make(map[int]bool)
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

		persister := storage.Persister{Root: root}
		Expect(persister.Exists(badContext)).To(BeTrue(),
			"Cannot run this test when context (%s) does not exist", badContext)

		task := Task{}
		err := persister.Unpersist(badContext, &task)
		Expect(err).To(HaveOccurred(),
			"Expected error from load from bad context (%s) but didn't get one", badContext)
	})

	It("returns the correct name", func() {
		t := newTask("task-a")
		Expect(t.name).To(Equal(t.Name()))
	})

	It("returns the correct id", func() {
		t := newTask("task-a")
		Expect(t.id).To(Equal(t.ID()))
	})
})
