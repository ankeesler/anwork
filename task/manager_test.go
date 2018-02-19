package task_test

import (
	"encoding/json"
	"fmt"

	"github.com/ankeesler/anwork/storage"
	"github.com/ankeesler/anwork/task"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	taskAName     = "task-a"
	taskAPriority = 20
	taskAState    = task.StateRunning

	taskBName     = "task-b"
	taskBPriority = 25

	taskCName     = "task-c"
	taskCPriority = 15
	taskCState    = task.StateBlocked

	goodProtobufManagerContext = "good-protobuf-manager-context"
)

var _ = Describe("Manager", func() {
	var (
		m *task.Manager
		p = storage.FilePersister{Root: root}
	)

	checkPersistence := func() {
		tmpM := task.NewManager()

		ExpectWithOffset(1, p.Exists(tmpContext)).To(BeFalse(),
			"Cannot run this test when context (%s) already exists", tmpContext)
		defer p.Delete(tmpContext)

		ExpectWithOffset(1, p.Persist(tmpContext, m)).To(Succeed())
		ExpectWithOffset(1, p.Unpersist(tmpContext, tmpM)).To(Succeed())
		ExpectWithOffset(1, m).To(Equal(tmpM))
	}

	BeforeEach(func() {
		m = task.NewManager()
	})

	Context("when no tasks are added", func() {
		It("has no tasks", func() {
			Expect(m.Tasks()).To(BeEmpty())
		})
		It("has no journal entries", func() {
			Expect(m.Journal().Events).To(BeEmpty())
		})
		It("persists correctly", checkPersistence)
		Context("when a task is deleted", func() {
			var (
				ret bool
			)
			BeforeEach(func() {
				ret = m.Delete(taskAName)
			})
			It("returns false", func() {
				Expect(ret).To(BeFalse())
			})
			It("still has no tasks", func() {
				Expect(m.Tasks()).To(HaveLen(0))
			})
			It("still has no journal entries", func() {
				Expect(m.Journal().Events).To(HaveLen(0))
			})
		})
	})

	Context("when one task is added", func() {
		BeforeEach(func() {
			m.Create(taskAName)
		})
		It("has one task with the expected name", func() {
			Expect(m.Tasks()).To(HaveLen(1))
			t := m.Tasks()[0]
			Expect(t.Name).To(Equal(taskAName))
		})
		It("has one journal entry (for the creation of the task)", func() {
			Expect(m.Journal().Events).To(HaveLen(1))
			Expect(m.Journal().Events[0].Title).To(Equal("Created task " + taskAName))
		})
		It("panics when we add a task with the same name", func() {
			Expect(func() { m.Create(taskAName) }).To(Panic())
		})
		It("panics when we try to set the state of a task that hasn't been added", func() {
			Expect(func() { m.SetState(taskBName, task.StateWaiting) }).To(Panic())
		})
		It("panics when we try to set the priority of a task that hasn't been added", func() {
			Expect(func() { m.SetPriority(taskCName, taskCPriority) }).To(Panic())
		})
		It("panics when we try to add a note for a task that hasn't been added", func() {
			Expect(func() { m.Note(taskCName, "tuna") }).To(Panic())
		})
		It("persists correctly", checkPersistence)
		Context("when that one task is modified", func() {
			BeforeEach(func() {
				m.SetPriority(taskAName, taskAPriority)
				m.SetState(taskAName, taskAState)
			})
			It("correctly tracks the task change", func() {
				t := m.Tasks()[0]
				Expect(t.Priority).To(BeEquivalentTo(taskAPriority))
				Expect(t.State).To(BeEquivalentTo(taskAState))
			})
			It("has 3 entries (creation, set priority, set state)", func() {
				Expect(m.Journal().Events).To(HaveLen(3))

				title := fmt.Sprintf("Created task %s", taskAName)
				Expect(m.Journal().Events[0].Title).To(Equal(title))

				title = fmt.Sprintf("Set priority on task %s from %d to %d",
					taskAName, task.DefaultPriority, taskAPriority)
				Expect(m.Journal().Events[1].Title).To(Equal(title))

				title = fmt.Sprintf("Set state on task %s from Waiting to %s",
					taskAName, task.StateNames[taskAState])
				Expect(m.Journal().Events[2].Title).To(Equal(title))
			})
			It("persists correctly", checkPersistence)
		})
		Context("when we add a note to that task", func() {
			var note string = "This is a note for task a"
			BeforeEach(func() {
				m.Note(taskAName, note)
			})
			It("adds an event to the journal", func() {
				Expect(m.Journal().Events).To(HaveLen(2))
			})
			It("stores the note in the journal", func() {
				journal := m.Journal()
				events := journal.Events
				event := events[len(events)-1]
				Expect(event.Title).To(ContainSubstring(note))
			})
			It("persists correctly", checkPersistence)
		})
		Context("when that one task is deleted", func() {
			var (
				ret bool
			)
			BeforeEach(func() {
				ret = m.Delete(taskAName)
			})
			It("returns true", func() {
				Expect(ret).To(BeTrue())
			})
			It("no longer stores that one task", func() {
				Expect(m.Tasks()).To(HaveLen(0))
			})
			It("does not successfully delete the task again", func() {
				Expect(m.Delete(taskAName)).To(BeFalse())
			})
			It("stores 2 events (creation, deletion)", func() {
				Expect(m.Journal().Events).To(HaveLen(2))
				Expect(m.Journal().Events[1].Title).To(Equal("Deleted task " + taskAName))
			})
			It("panics if we try to act on that deleted task", func() {
				Expect(func() { m.SetPriority(taskAName, taskAPriority) }).To(Panic())
				Expect(func() { m.SetState(taskAName, taskAState) }).To(Panic())
				Expect(func() { m.Note(taskAName, "tuna") }).To(Panic())
			})
			It("persists correctly", checkPersistence)
		})
	})

	Context("when more than one task is added", func() {
		BeforeEach(func() {
			m.Create(taskAName)
			m.Create(taskBName)
			m.Create(taskCName)
		})
		It("panics for all calls to Create with task names that have already been added", func() {
			Expect(func() { m.Create(taskAName) }).To(Panic())
			Expect(func() { m.Create(taskBName) }).To(Panic())
			Expect(func() { m.Create(taskCName) }).To(Panic())
		})
		It("has three tasks", func() {
			Expect(m.Tasks()).To(HaveLen(3))
		})
		It("orders the three tasks by their increasing ids, since the priorities are the same", func() {
			t := m.Tasks()[0]
			Expect(t.Name).To(Equal(taskAName))
			t = m.Tasks()[1]
			Expect(t.Name).To(Equal(taskBName))
			t = m.Tasks()[2]
			Expect(t.Name).To(Equal(taskCName))
		})
		It("has three journal entries for each of the creations", func() {
			Expect(m.Journal().Events).To(HaveLen(3))

			title := fmt.Sprintf("Created task %s", taskAName)
			Expect(m.Journal().Events[0].Title).To(Equal(title))

			title = fmt.Sprintf("Created task %s", taskBName)
			Expect(m.Journal().Events[1].Title).To(Equal(title))

			title = fmt.Sprintf("Created task %s", taskCName)
			Expect(m.Journal().Events[2].Title).To(Equal(title))
		})
		It("persists correctly", checkPersistence)
		Context("when tasks are updated", func() {
			BeforeEach(func() {
				taskA := m.FindByName(taskAName)
				Expect(taskA).ToNot(BeNil())
				m.SetPriority(taskA.Name, taskAPriority)
				m.SetState(taskA.Name, taskAState)

				taskB := m.FindByName(taskBName)
				Expect(taskB).ToNot(BeNil())
				m.SetPriority(taskB.Name, taskBPriority)

				taskC := m.FindByName(taskCName)
				Expect(taskC).ToNot(BeNil())
				m.SetPriority(taskC.Name, taskCPriority)
				m.SetState(taskC.Name, taskCState)
			})
			It("re-orders the tasks by priority ", func() {
				tasks := m.Tasks()
				Expect(tasks[0].Name).To(Equal(taskCName))
				Expect(tasks[1].Name).To(Equal(taskAName))
				Expect(tasks[2].Name).To(Equal(taskBName))
			})
			It("correctly tracks the changes", func() {
				tasks := m.Tasks()

				Expect(tasks[0].Priority).To(BeEquivalentTo(taskCPriority))
				Expect(tasks[0].State).To(BeEquivalentTo(taskCState))

				Expect(tasks[1].Priority).To(BeEquivalentTo(taskAPriority))
				Expect(tasks[1].State).To(BeEquivalentTo(taskAState))

				Expect(tasks[2].Priority).To(BeEquivalentTo(taskBPriority))
			})
			It("tracks the actions in the journal ", func() {
				// 3 creations
				// 1 set taskA pri
				// 1 set taskA state
				// 1 set taskB pri
				// 1 set taskC pri
				// 1 set taskC state
				// = 8 events
				Expect(m.Journal().Events).To(HaveLen(8))
			})
			It("persists correctly", checkPersistence)
			XIt("persists correctly through reset", func() {
				tmpM := task.NewManager()

				Expect(p.Exists(tmpContext)).To(BeFalse(),
					"Cannot run this test when context (%s) already exists", tmpContext)
				defer p.Delete(tmpContext)

				Expect(p.Persist(tmpContext, m)).To(Succeed())

				// Set the nextTaskID to 0 to simulate a new runtime.
				//nextTaskID = 0

				Expect(p.Unpersist(tmpContext, tmpM)).To(Succeed())
				Expect(m).To(Equal(tmpM))

				m.Create("new")
				newT := m.FindByName("new")
				Expect(newT.ID).ToNot(BeEquivalentTo(0))
			})
			XIt("maintains task ID invariant through reset", func() {
				tmpM := task.NewManager()

				Expect(p.Exists(tmpContext)).To(BeFalse(),
					"Cannot run this test when context (%s) already exists", tmpContext)
				defer p.Delete(tmpContext)

				taskC := m.FindByName(taskCName)
				Expect(taskC).ToNot(BeNil())
				taskC = m.FindByID(taskC.ID)
				Expect(taskC).ToNot(BeNil())

				Expect(m.Delete(taskCName)).To(BeTrue())
				Expect(p.Persist(tmpContext, m)).To(Succeed())

				// Set the nextTaskID to 0 to simulate a new runtime.
				//nextTaskID = 0

				Expect(p.Unpersist(tmpContext, tmpM)).To(Succeed())
				Expect(m).To(Equal(tmpM))

				m.Create("new")
				newT := m.FindByName("new")
				Expect(newT.ID).ToNot(BeEquivalentTo(taskC.ID))
			})
			Context("when one task is deleted", func() {
				var (
					ret bool
				)
				BeforeEach(func() {
					ret = m.Delete(taskBName)
				})
				It("holds one fewer task", func() {
					Expect(m.Tasks()).To(HaveLen(2))
				})
				It("no longer store the deleted task", func() {
					taskB := m.FindByName(taskBName)
					Expect(taskB).To(BeNil())
				})
				It("continue to store the other tasks in the correct order", func() {
					tasks := m.Tasks()
					Expect(tasks[0].Name).To(Equal(taskCName))
					Expect(tasks[0].Priority).To(BeEquivalentTo(taskCPriority))
					Expect(tasks[0].State).To(BeEquivalentTo(taskCState))

					Expect(tasks[1].Name).To(Equal(taskAName))
					Expect(tasks[1].Priority).To(BeEquivalentTo(taskAPriority))
					Expect(tasks[1].State).To(BeEquivalentTo(taskAState))
				})
				It("appends a journal entry for the deletion", func() {
					// 3 creations
					// 1 set taskA pri
					// 1 set taskA state
					// 1 set taskB pri
					// 1 set taskC pri
					// 1 set taskC state
					// 1 deletion
					// = 9 events
					Expect(m.Journal().Events).To(HaveLen(9))
				})
				It("persists correctly", checkPersistence)
				Context("when the other two tasks' priorities are set equal", func() {
					BeforeEach(func() {
						tasks := m.Tasks()
						tasks[0].Priority = tasks[1].Priority
					})
					It("sorts the tasks by their IDs", func() {
						tasks := m.Tasks()
						Expect(tasks[0].Name).To(Equal(taskAName))
						Expect(tasks[1].Name).To(Equal(taskCName))
					})
				})
				Context("when the other two tasks are deleted", func() {
					var retDeleteA, retDeleteC bool
					BeforeEach(func() {
						retDeleteA = m.Delete(taskAName)
						retDeleteC = m.Delete(taskCName)
					})
					It("returns true for both deletions", func() {
						Expect(retDeleteA).To(BeTrue())
						Expect(retDeleteC).To(BeTrue())
					})
					It("no longer stores any tasks", func() {
						Expect(m.Tasks()).To(BeEmpty())
					})
					It("returns false when trying to delete any deleted task", func() {
						Expect(m.Delete(taskAName)).To(BeFalse())
						Expect(m.Delete(taskBName)).To(BeFalse())
						Expect(m.Delete(taskCName)).To(BeFalse())
					})
					It("appends 2 journal entries for the deletions", func() {
						// 3 creations
						// 1 set taskA pri
						// 1 set taskA state
						// 1 set taskB pri
						// 1 set taskC pri
						// 1 set taskC state
						// 3 deletions
						// = 11 events
						Expect(m.Journal().Events).To(HaveLen(11))
					})
					It("persists correctly", checkPersistence)
				})
			})
		})
	})

	Context("when printed", func() {
		It("doesn't explode", func() {
			m.Create("task-a")
			m.Create("task-b")
			m.SetState("task-a", task.StateRunning)
			m.SetPriority("task-b", task.DefaultPriority-1)
			m.SetState("task-a", task.StateWaiting)
			m.SetState("task-b", task.StateRunning)
			Expect(fmt.Sprintf("%s", m)).ToNot(BeNil())
		})
	})

	It("fails gracefully when loaded from a bad context", func() {
		Expect(p.Exists(badContext)).To(BeTrue(),
			"Cannot run this test when context (%s) does not exist", badContext)
		Expect(p.Unpersist(badContext, &task.Manager{})).ToNot(Succeed())
	})

	It("serializes to json", func() {
		m.Create("task-a")
		m.Create("task-b")
		m.SetState("task-a", task.StateRunning)
		m.SetPriority("task-b", task.DefaultPriority-1)
		m.SetState("task-a", task.StateWaiting)
		m.SetState("task-b", task.StateRunning)
		bytes, err := m.Serialize()
		Expect(err).NotTo(HaveOccurred())

		var jsonM task.Manager
		Expect(json.Unmarshal(bytes, &jsonM)).To(Succeed())
		Expect(&jsonM).To(Equal(m))
	})

	It("successfully unpersists from protocol buffers (legacy)", func() {
		Expect(p.Exists(goodProtobufManagerContext)).To(BeTrue())
		Expect(p.Unpersist(goodProtobufManagerContext, m)).To(Succeed())
		Expect(m.Tasks()).To(BeEmpty())
		Expect(m.Journal().Events).To(HaveLen(11))
	})
})
