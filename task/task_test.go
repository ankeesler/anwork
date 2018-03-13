package task_test

import (
	"encoding/json"
	"fmt"

	"github.com/ankeesler/anwork/storage"
	"github.com/ankeesler/anwork/task"
	pb "github.com/ankeesler/anwork/task/proto"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	goodContext         = "good-context"
	goodProtobufContext = "good-protobuf-context"
)

var _ = Describe("State constants", func() {
	It("should line up with Protobuf definitions", func() {
		statePairs := [][]task.State{
			{task.StateWaiting, task.State(pb.State_WAITING)},
			{task.StateBlocked, task.State(pb.State_BLOCKED)},
			{task.StateRunning, task.State(pb.State_RUNNING)},
			{task.StateFinished, task.State(pb.State_FINISHED)},
		}
		for _, statePair := range statePairs {
			Expect(statePair[0]).To(Equal(statePair[1]))
		}
	})
	It("should have the expected state names", func() {
		Expect(task.StateNames[task.StateWaiting]).To(Equal("Waiting"))
		Expect(task.StateNames[task.StateBlocked]).To(Equal("Blocked"))
		Expect(task.StateNames[task.StateRunning]).To(Equal("Running"))
		Expect(task.StateNames[task.StateFinished]).To(Equal("Finished"))
	})
})

var _ = Describe("Task", func() {
	var (
		t *task.Task
	)

	BeforeEach(func() {
		t = &task.Task{
			Name:        "this is a task",
			ID:          15,
			Description: "this is a task description",
			StartDate:   150,
			Priority:    27,
			State:       task.StateRunning,
		}
	})

	It("are serialized via JSON", func() {
		bytes, err := t.Serialize()
		Expect(err).NotTo(HaveOccurred())

		var jsonTask *task.Task = &task.Task{}
		Expect(json.Unmarshal(bytes, jsonTask)).To(Succeed())
		Expect(jsonTask).To(Equal(t))
	})

	Describe("Persistence", func() {
		var (
			p storage.Persister
		)

		BeforeEach(func() {
			p = &storage.FilePersister{Root: root}
		})

		successfulPersistence := func() {
			Expect(p.Exists(tmpContext)).To(BeFalse(),
				"Cannot run this test when context (%s) already exists", tmpContext)
			defer p.Delete(tmpContext)

			Expect(p.Persist(tmpContext, t)).To(Succeed())

			var unpersistedTask task.Task
			Expect(p.Unpersist(tmpContext, &unpersistedTask)).To(Succeed())

			Expect(t).To(Equal(&unpersistedTask))
		}

		Context("when persisting", func() {
			It("successfully persists and then unpersists", func() {
				successfulPersistence()
			})
		})

		Context("when unpersisting", func() {
			var (
				expectedTask *task.Task
			)
			BeforeEach(func() {
				expectedTask = &task.Task{
					Name:        "task-a",
					ID:          0,
					Description: "Here is a description!",
					StartDate:   1513548268,
					Priority:    612,
					State:       task.StateRunning,
				}
			})

			successfulUnpersistence := func(context string) {
				Expect(p.Exists(context)).To(BeTrue(),
					"Cannot run this test when context (%s) does not exist", context)

				var unpersistedTask task.Task
				Expect(p.Unpersist(context, &unpersistedTask)).To(Succeed())

				Expect(&unpersistedTask).To(Equal(expectedTask))
			}

			It("successfully unpersists from a known-good protobuf context (legacy)", func() {
				successfulUnpersistence(goodProtobufContext)
			})

			It("successfully unpersists from a known-good json context", func() {
				successfulUnpersistence(goodContext)
			})

			It("fails to unpersist from a known-bad context", func() {
				Expect(p.Exists(badContext)).To(BeTrue(),
					"Cannot run this test when context (%s) does not exist", badContext)

				var unpersistedTask task.Task
				Expect(p.Unpersist(badContext, &unpersistedTask)).NotTo(Succeed())
			})
		})
	})

	Context("have unique ID's", func() {

		It("that are larger", func() {
			const taskWithID5Context = "protobuf-task-with-id-5-context"
			persister := storage.FilePersister{Root: root}
			Expect(persister.Exists(taskWithID5Context)).To(BeTrue(),
				"Cannot run this test when context (%s) does not exist", taskWithID5Context)

			taskWithID5 := task.Task{}
			persister.Unpersist(taskWithID5Context, &taskWithID5)
			Expect(taskWithID5.ID).To(BeEquivalentTo(5))
			Expect(taskWithID5.Name).To(Equal("task-with-id-5"))

			for i := 0; i < 10; i++ {
				name := fmt.Sprintf("task-%d", i)
				t := task.NewTask(name)
				Expect(t.ID).ToNot(Equal(5))
			}
		})

		It("that are smaller", func() {
			persister := storage.FilePersister{Root: root}
			Expect(persister.Exists(tmpContext)).To(BeFalse(),
				"Cannot run this test when context (%s) exists", tmpContext)
			defer persister.Delete(tmpContext)

			t := task.Task{}
			t.ID = 0
			err := persister.Persist(tmpContext, &t)
			Expect(err).ToNot(HaveOccurred(),
				"Got error when persisting task %#v to file: %s", t, err)

			ids := make(map[int]bool)
			for i := 0; i < 25; i++ {
				name := fmt.Sprintf("task-%d", i)
				tt := task.NewTask(name)
				ids[tt.ID] = true
			}

			err = persister.Unpersist(tmpContext, &t)
			Expect(err).ToNot(HaveOccurred(),
				"Got error when unpersisting task %#v from file: %s", t, err)
			Expect(t.ID).To(BeZero(),
				"Expected task id to be %d but was %d", 0, t.ID)

			for i := 0; i < 25; i++ {
				name := fmt.Sprintf("task-%d", i)
				tt := task.NewTask(name)
				_, ok := ids[tt.ID]
				Expect(ok).To(BeFalse(), "Failure! Task ID %d already exists!", tt.ID)
			}
		})
	})
})
