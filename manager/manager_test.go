package manager_test

import (
	"errors"
	"time"

	"code.cloudfoundry.org/clock/fakeclock"
	managerpkg "github.com/ankeesler/anwork/manager"
	"github.com/ankeesler/anwork/task2"
	"github.com/ankeesler/anwork/task2/task2fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manager", func() {
	var (
		repo    *task2fakes.FakeRepo
		now     time.Time
		clock   *fakeclock.FakeClock
		manager managerpkg.Manager
	)

	BeforeEach(func() {
		repo = &task2fakes.FakeRepo{}

		now = time.Now()
		clock = fakeclock.NewFakeClock(now)

		manager = managerpkg.New(repo, clock)
	})

	Describe("Create", func() {
		It("calls down to the repo to create a task", func() {
			Expect(manager.Create("task-a")).To(Succeed())
			Expect(repo.CreateTaskCallCount()).To(Equal(1))
			Expect(*repo.CreateTaskArgsForCall(0)).To(Equal(task2.Task{
				Name:      "task-a",
				StartDate: now.Unix(),
				Priority:  10,               // default
				State:     task2.StateReady, // default
			}))
		})

		Context("when the repo returns an error", func() {
			BeforeEach(func() {
				repo.CreateTaskReturnsOnCall(0, errors.New("some create task error"))
			})
			It("returns that error", func() {
				Expect(manager.Create("task-a")).To(MatchError("some create task error"))
			})
		})
	})

	//	Context("when no tasks are created", func() {
	//		It("returns no tasks", func() {
	//			ts := manager.Tasks()
	//			Expect(ts).NotTo(BeNil())
	//			Expect(ts).To(BeEmpty())
	//		})
	//		It("returns no journal entries", func() {
	//			Expect(manager.Events()).To(BeEmpty())
	//		})
	//		It("returns an error when deleting a task and does not add a journal entry", func() {
	//			err := manager.Delete("1")
	//			Expect(err).To(HaveOccurred())
	//			Expect(err.Error()).To(ContainSubstring("cannot find task with name 1"))
	//			Expect(manager.Events()).To(BeEmpty())
	//		})
	//		It("fails to add a note for any task", func() {
	//			Expect(manager.Note("1", "tuna")).To(HaveOccurred())
	//		})
	//		It("fails to add a set the state of a task", func() {
	//			Expect(manager.SetState("1", StateRunning)).To(HaveOccurred())
	//		})
	//		It("fails to add a set the priority of a task", func() {
	//			Expect(manager.SetPriority("1", 5)).To(HaveOccurred())
	//		})
	//		It("fails to delete any event", func() {
	//			Expect(manager.DeleteEvent(12345)).NotTo(Succeed())
	//		})
	//	})
	//
	//	Context("when three tasks are created", func() {
	//		var (
	//			names []string = []string{"1", "2", "3"}
	//			ids   []int
	//		)
	//		BeforeEach(func() {
	//			for _, name := range names {
	//				Expect(manager.Create(name)).To(Succeed())
	//			}
	//
	//			ids = []int{}
	//			for _, t := range manager.Tasks() {
	//				ids = append(ids, t.ID)
	//			}
	//		})
	//		It("returns those three tasks in the order that they were added (since their priorities are all the same)", func() {
	//			tasks := manager.Tasks()
	//			Expect(tasks).To(HaveLen(3))
	//			for i := range names {
	//				Expect(tasks[i].Name).To(Equal(names[i]))
	//				Expect(tasks[i].Priority).To(Equal(DefaultPriority))
	//			}
	//		})
	//		It("returns an error when we try to recreate the existing tasks", func() {
	//			for _, name := range names {
	//				Expect(manager.Create(name)).NotTo(Succeed())
	//			}
	//		})
	//		It("can find those tasks by id", func() {
	//			for _, id := range ids {
	//				t := manager.FindByID(id)
	//				Expect(t.ID).To(Equal(id))
	//			}
	//		})
	//		It("can find those tasks by name", func() {
	//			for _, name := range names {
	//				t := manager.FindByName(name)
	//				Expect(t.Name).To(Equal(name))
	//			}
	//		})
	//		It("returns these events in the journal", func() {
	//			events := manager.Events()
	//			Expect(events).To(HaveLen(3))
	//			for i, name := range names {
	//				t := manager.FindByName(name)
	//				Expect(events[i].Title).To(Equal(fmt.Sprintf("Created task '%s'", name)))
	//				Expect(events[i].Type).To(Equal(EventTypeCreate))
	//				Expect(events[i].TaskID).To(Equal(t.ID))
	//			}
	//		})
	//		It("can delete these events in the journal", func() {
	//			events := manager.Events()
	//			copiedEvents := make([]*Event, len(events))
	//			copy(copiedEvents, events)
	//			for _, event := range copiedEvents {
	//				Expect(manager.DeleteEvent(event.Date)).To(Succeed())
	//			}
	//			Expect(manager.Events()).To(HaveLen(0))
	//		})
	//		It("uses different IDs for every task", func() {
	//			ids := make(map[int]bool)
	//			for _, t := range manager.Tasks() {
	//				if ids[t.ID] {
	//					Fail(fmt.Sprintf("Expected task %s to not have ID %d", t.Name, t.ID))
	//				}
	//				ids[t.ID] = true
	//			}
	//		})
	//		It("can rename a task", func() {
	//			Expect(manager.Rename("1", "5")).To(Succeed())
	//
	//			t := manager.FindByName("5")
	//			Expect(t).ToNot(BeNil())
	//			Expect(t.Name).To(Equal("5"))
	//
	//			t = manager.FindByName("1")
	//			Expect(t).To(BeNil())
	//		})
	//		It("returns an error when trying to rename an unknown task", func() {
	//			Expect(manager.Rename("5", "6")).NotTo(Succeed())
	//		})
	//		It("returns an error when trying to rename a task to an existing task", func() {
	//			Expect(manager.Rename("5", "1")).NotTo(Succeed())
	//		})
	//
	//		Context("when a note is added to those tasks", func() {
	//			var (
	//				notes = []string{"tuna", "fish", "marlin"}
	//			)
	//			BeforeEach(func() {
	//				for i := range notes {
	//					Expect(manager.Note(names[i], notes[i])).To(Succeed())
	//				}
	//			})
	//			It("returns the notes in the journal", func() {
	//				events := manager.Events()
	//				Expect(events).To(HaveLen(6))
	//				for i := 3; i < 6; i++ {
	//					t := manager.FindByName(names[i-3])
	//					Expect(events[i].Title).To(ContainSubstring(notes[i-3]))
	//					Expect(events[i].Type).To(Equal(EventTypeNote))
	//					Expect(events[i].TaskID).To(Equal(t.ID))
	//				}
	//			})
	//		})
	//
	//		Context("when the priority is set on those tasks", func() {
	//			var (
	//				priorities = []int{6, 5, 7}
	//			)
	//			BeforeEach(func() {
	//				for i := range priorities {
	//					Expect(manager.SetPriority(names[i], priorities[i])).To(Succeed())
	//				}
	//			})
	//			It("returns the tasks in the order of priority", func() {
	//				tasks := manager.Tasks()
	//				Expect(tasks[0].Name).To(Equal("2"))
	//				Expect(tasks[1].Name).To(Equal("1"))
	//				Expect(tasks[2].Name).To(Equal("3"))
	//			})
	//			It("returns these events in the journal", func() {
	//				events := manager.Events()
	//				Expect(events).To(HaveLen(6))
	//				for i := 3; i < 6; i++ {
	//					t := manager.FindByName(names[i-3])
	//					msg := fmt.Sprintf("Set priority on task '%s' from %d to %d", t.Name,
	//						DefaultPriority, priorities[i-3])
	//					Expect(events[i].Title).To(Equal(msg))
	//					Expect(events[i].Type).To(Equal(EventTypeSetPriority))
	//					Expect(events[i].TaskID).To(Equal(t.ID))
	//				}
	//			})
	//		})
	//
	//		Context("when the state is set on those tasks", func() {
	//			var (
	//				states = []State{StateFinished, StateRunning, StateBlocked}
	//			)
	//			BeforeEach(func() {
	//				for i := range states {
	//					Expect(manager.SetState(names[i], states[i])).To(Succeed())
	//				}
	//			})
	//			It("returns the tasks with the correct state", func() {
	//				for i, task := range manager.Tasks() {
	//					Expect(task.State).To(Equal(states[i]))
	//				}
	//			})
	//			It("returns these events in the journal", func() {
	//				events := manager.Events()
	//				Expect(events).To(HaveLen(6))
	//				for i := 3; i < 6; i++ {
	//					t := manager.FindByName(names[i-3])
	//					msg := fmt.Sprintf("Set state on task '%s' from %s to %s", t.Name,
	//						StateNames[StateReady], StateNames[states[i-3]])
	//					Expect(events[i].Title).To(Equal(msg))
	//					Expect(events[i].Type).To(Equal(EventTypeSetState))
	//					Expect(events[i].TaskID).To(Equal(t.ID))
	//				}
	//			})
	//		})
	//
	//		Context("when reset is called", func() {
	//			It("deletes all of those tasks when Reset() is called", func() {
	//				Expect(manager.Reset()).To(Succeed())
	//				Expect(manager.Tasks()).To(BeEmpty())
	//				Expect(manager.Events()).To(BeEmpty())
	//			})
	//
	//			It("resets the ID that it uses back to 0", func() {
	//				Expect(manager.Reset()).To(Succeed())
	//				Expect(manager.Create("a")).To(Succeed())
	//				t := manager.FindByName("a")
	//				Expect(t).NotTo(BeNil())
	//				Expect(t.ID).To(Equal(0))
	//			})
	//		})
	//
	//		Context("when those three tasks are deleted", func() {
	//			BeforeEach(func() {
	//				for _, name := range names {
	//					Expect(manager.Delete(name)).To(Succeed())
	//				}
	//			})
	//			It("returns no tasks", func() {
	//				Expect(manager.Tasks()).To(BeEmpty())
	//			})
	//			It("returns an error when trying to delete the tasks again", func() {
	//				for _, name := range names {
	//					err := manager.Delete(name)
	//					Expect(err).To(HaveOccurred())
	//					Expect(err.Error()).To(ContainSubstring("cannot find task with name %s", name))
	//				}
	//			})
	//			It("cannot find those tasks by name", func() {
	//				for _, name := range names {
	//					Expect(manager.FindByName(name)).To(BeNil())
	//				}
	//			})
	//			It("returns these events in the journal", func() {
	//				events := manager.Events()
	//				Expect(events).To(HaveLen(6))
	//				for i := 3; i < 6; i++ {
	//					Expect(events[i].Title).To(Equal(fmt.Sprintf("Deleted task '%s'", names[i-3])))
	//					Expect(events[i].Type).To(Equal(EventTypeDelete))
	//				}
	//			})
	//			It("fails to add a note for those tasks", func() {
	//				for _, name := range names {
	//					Expect(manager.Note(name, "tuna")).To(HaveOccurred())
	//				}
	//			})
	//			It("fails to set the priority for those tasks", func() {
	//				for _, name := range names {
	//					Expect(manager.SetPriority(name, 5)).To(HaveOccurred())
	//				}
	//			})
	//			It("fails to set the state for those tasks", func() {
	//				for _, name := range names {
	//					Expect(manager.SetState(name, StateRunning)).To(HaveOccurred())
	//				}
	//			})
	//		})
	//	})
})
