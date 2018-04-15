package task

import (
	"fmt"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// This is a utility method to run tests with an object conforming to the
// Manager interface.
func RunManagerTests(factory ManagerFactory) {
	var (
		manager Manager
	)

	BeforeEach(func() {
		var err error
		manager, err = factory.Create()
		Expect(err).NotTo(HaveOccurred())
	})

	Context("when no tasks are created", func() {
		It("returns no tasks", func() {
			Expect(manager.Tasks()).To(BeEmpty())
		})
		It("returns no journal entries", func() {
			Expect(manager.Events()).To(BeEmpty())
		})
		It("returns false when deleting a task and does not add a journal entry", func() {
			Expect(manager.Delete("1")).To(BeFalse())
			Expect(manager.Events()).To(BeEmpty())
		})
	})

	Context("when three tasks are created", func() {
		var (
			names []string = []string{"1", "2", "3"}
			ids   []int
		)
		BeforeEach(func() {
			for _, name := range names {
				Expect(manager.Create(name)).To(Succeed())
			}

			ids = []int{}
			for _, t := range manager.Tasks() {
				ids = append(ids, t.ID)
			}
		})
		It("returns those three tasks in the order that they were added (since their priorities are all the same)", func() {
			tasks := manager.Tasks()
			Expect(tasks).To(HaveLen(3))
			for i := range names {
				Expect(tasks[i].Name).To(Equal(names[i]))
				Expect(tasks[i].Priority).To(Equal(DefaultPriority))
			}
		})
		It("returns an error when we try to recreate the existing tasks", func() {
			for _, name := range names {
				Expect(manager.Create(name)).NotTo(Succeed())
			}
		})
		It("can find those tasks by id", func() {
			for _, id := range ids {
				t := manager.FindByID(id)
				Expect(t.ID).To(Equal(id))
			}
		})
		It("can find those tasks by name", func() {
			for _, name := range names {
				t := manager.FindByName(name)
				Expect(t.Name).To(Equal(name))
			}
		})
		It("returns these events in the journal", func() {
			events := manager.Events()
			Expect(events).To(HaveLen(3))
			for i, name := range names {
				t := manager.FindByName(name)
				Expect(events[i].Title).To(Equal(fmt.Sprintf("created task '%s'", name)))
				Expect(events[i].Type).To(Equal(EventTypeCreate))
				Expect(events[i].TaskID).To(Equal(t.ID))
			}
		})
		It("uses different IDs for every task", func() {
			ids := make(map[int]bool)
			for _, t := range manager.Tasks() {
				if ids[t.ID] {
					Fail(fmt.Sprintf("Expected task %s to not have ID %d", t.Name, t.ID))
				}
				ids[t.ID] = true
			}
		})

		Context("when a note is added to those tasks", func() {
			var (
				notes = []string{"tuna", "fish", "marlin"}
			)
			BeforeEach(func() {
				for i := range notes {
					manager.Note(names[i], notes[i])
				}
			})
			It("returns the notes in the journal", func() {
				events := manager.Events()
				Expect(events).To(HaveLen(6))
				for i := 3; i < 6; i++ {
					t := manager.FindByName(names[i-3])
					Expect(events[i].Title).To(ContainSubstring(notes[i-3]))
					Expect(events[i].Type).To(Equal(EventTypeNote))
					Expect(events[i].TaskID).To(Equal(t.ID))
				}
			})
		})

		Context("when the priority is set on those tasks", func() {
			var (
				priorities = []int{6, 5, 7}
			)
			BeforeEach(func() {
				for i := range priorities {
					manager.SetPriority(names[i], priorities[i])
				}
			})
			It("returns the tasks in the order of priority", func() {
				tasks := manager.Tasks()
				Expect(tasks[0].Name).To(Equal("2"))
				Expect(tasks[1].Name).To(Equal("1"))
				Expect(tasks[2].Name).To(Equal("3"))
			})
			It("returns these events in the journal", func() {
				events := manager.Events()
				Expect(events).To(HaveLen(6))
				for i := 3; i < 6; i++ {
					t := manager.FindByName(names[i-3])
					msg := fmt.Sprintf("set priority on task '%s' from %d to %d", t.Name,
						DefaultPriority, priorities[i-3])
					Expect(events[i].Title).To(Equal(msg))
					Expect(events[i].Type).To(Equal(EventTypeSetPriority))
					Expect(events[i].TaskID).To(Equal(t.ID))
				}
			})
		})

		Context("when the state is set on those tasks", func() {
			var (
				states = []State{StateFinished, StateRunning, StateBlocked}
			)
			BeforeEach(func() {
				for i := range states {
					manager.SetState(names[i], states[i])
				}
			})
			It("returns the tasks with the correct state", func() {
				for i, task := range manager.Tasks() {
					Expect(task.State).To(Equal(states[i]))
				}
			})
			It("returns these events in the journal", func() {
				events := manager.Events()
				Expect(events).To(HaveLen(6))
				for i := 3; i < 6; i++ {
					t := manager.FindByName(names[i-3])
					msg := fmt.Sprintf("set state on task '%s' from %s to %s", t.Name,
						StateNames[StateWaiting], StateNames[states[i-3]])
					Expect(events[i].Title).To(Equal(msg))
					Expect(events[i].Type).To(Equal(EventTypeSetState))
					Expect(events[i].TaskID).To(Equal(t.ID))
				}
			})
		})

		Context("when those three tasks are deleted", func() {
			BeforeEach(func() {
				for _, name := range names {
					Expect(manager.Delete(name)).To(BeTrue())
				}
			})
			It("returns no tasks", func() {
				Expect(manager.Tasks()).To(BeEmpty())
			})
			It("returns false when trying to delete the tasks again", func() {
				for _, name := range names {
					Expect(manager.Delete(name)).To(BeFalse())
				}
			})
			It("cannot find those tasks by name", func() {
				for _, name := range names {
					Expect(manager.FindByName(name)).To(BeNil())
				}
			})
			It("returns these events in the journal", func() {
				events := manager.Events()
				Expect(events).To(HaveLen(6))
				for i := 3; i < 6; i++ {
					Expect(events[i].Title).To(Equal(fmt.Sprintf("deleted task '%s'", names[i-3])))
					Expect(events[i].Type).To(Equal(EventTypeDelete))
				}
			})
		})
	})
}
