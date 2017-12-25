package task

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	taskAName     = "task-a"
	taskAPriority = 20
	taskAState    = TaskStateRunning

	taskBName     = "task-b"
	taskBPriority = 25

	taskCName     = "task-c"
	taskCPriority = 15
	taskCState    = TaskStateBlocked
)

var _ = Describe("Manager", func() {
	var (
		m *Manager
	)

	BeforeEach(func() {
		m = NewManager()
	})

	Context("when no tasks are added", func() {
		It("has no tasks", func() {
			Expect(m.Tasks()).To(BeEmpty())
		})
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
		})
	})

	Context("when one task is added", func() {
		BeforeEach(func() {
			m.Create(taskAName)
		})
		It("has one task with the expected name", func() {
			Expect(m.Tasks()).To(HaveLen(1))
			t := m.Tasks()[0]
			Expect(t.name).To(Equal(taskAName))
		})
		It("panics when we add a task with the same name", func() {
			Expect(func() { m.Create(taskAName) }).To(Panic())
		})
		Context("when that one task is modified", func() {
			BeforeEach(func() {
				t := m.Tasks()[0]
				t.priority = taskAPriority
				t.state = taskAState
			})
			It("correctly tracks the task change", func() {
				t := m.Tasks()[0]
				Expect(t.priority).To(BeEquivalentTo(taskAPriority))
				Expect(t.state).To(BeEquivalentTo(taskAState))
			})
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
			Expect(t.name).To(Equal(taskAName))
			t = m.Tasks()[1]
			Expect(t.name).To(Equal(taskBName))
			t = m.Tasks()[2]
			Expect(t.name).To(Equal(taskCName))
		})
		Context("when tasks are updated", func() {
			BeforeEach(func() {
				taskA := m.Find(taskAName)
				Expect(taskA).ToNot(BeNil())
				taskA.priority = taskAPriority
				taskA.state = taskAState

				taskB := m.Find(taskBName)
				Expect(taskB).ToNot(BeNil())
				taskB.priority = taskBPriority

				taskC := m.Find(taskCName)
				Expect(taskC).ToNot(BeNil())
				taskC.priority = taskCPriority
				taskC.state = taskCState
			})
			It("re-orders the tasks by priority ", func() {
				tasks := m.Tasks()
				Expect(tasks[0].name).To(Equal(taskCName))
				Expect(tasks[1].name).To(Equal(taskAName))
				Expect(tasks[2].name).To(Equal(taskBName))
			})
			It("correctly tracks the changes", func() {
				tasks := m.Tasks()

				Expect(tasks[0].priority).To(BeEquivalentTo(taskCPriority))
				Expect(tasks[0].state).To(BeEquivalentTo(taskCState))

				Expect(tasks[1].priority).To(BeEquivalentTo(taskAPriority))
				Expect(tasks[1].state).To(BeEquivalentTo(taskAState))

				Expect(tasks[2].priority).To(BeEquivalentTo(taskBPriority))
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
					taskB := m.Find(taskBName)
					Expect(taskB).To(BeNil())
				})
				It("continue to store the other tasks in the correct order", func() {
					tasks := m.Tasks()
					Expect(tasks[0].name).To(Equal(taskCName))
					Expect(tasks[0].priority).To(BeEquivalentTo(taskCPriority))
					Expect(tasks[0].state).To(BeEquivalentTo(taskCState))

					Expect(tasks[1].name).To(Equal(taskAName))
					Expect(tasks[1].priority).To(BeEquivalentTo(taskAPriority))
					Expect(tasks[1].state).To(BeEquivalentTo(taskAState))
				})
				Context("when the other two tasks' priorities are set equal", func() {
					BeforeEach(func() {
						tasks := m.Tasks()
						tasks[0].priority = tasks[1].priority
					})
					It("sorts the tasks by their IDs", func() {
						tasks := m.Tasks()
						Expect(tasks[0].name).To(Equal(taskAName))
						Expect(tasks[1].name).To(Equal(taskCName))
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
				})
			})
		})
	})
})
