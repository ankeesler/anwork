package task2

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// RunRepoTests will run a set of tests to verify that the provided repo
// is a valid Repo implementation.
func RunRepoTests(createRepoFunc func() Repo) {
	var (
		repo                   Repo
		taskA, taskB, taskC    *Task
		eventA, eventB, eventC *Event
	)
	BeforeEach(func() {
		repo = createRepoFunc()

		taskA = &Task{Name: "task-a", ID: 1}
		taskB = &Task{Name: "task-b", ID: 2}
		taskC = &Task{Name: "task-c", ID: 3}

		eventA = &Event{Title: "event-a", Date: 1}
		eventB = &Event{Title: "event-b", Date: 2}
		eventC = &Event{Title: "event-c", Date: 3}
	})

	Describe("CreateTask", func() {
		Context("when tasks are created", func() {
			BeforeEach(func() {
				Expect(repo.CreateTask(taskA)).To(Succeed())
				Expect(repo.CreateTask(taskB)).To(Succeed())
				Expect(repo.CreateTask(taskC)).To(Succeed())
			})
			It("returns then with Tasks()", func() {
				tasks, err := repo.Tasks()
				Expect(err).NotTo(HaveOccurred())
				Expect(*tasks[0]).To(Equal(*taskA))
				Expect(*tasks[1]).To(Equal(*taskB))
				Expect(*tasks[2]).To(Equal(*taskC))
			})
		})
		Context("when a task with that ID already exists", func() {
			BeforeEach(func() {
				Expect(repo.CreateTask(taskA)).To(Succeed())
				Expect(repo.CreateTask(taskB)).To(Succeed())
				Expect(repo.CreateTask(taskC)).To(Succeed())
			})
			It("returns an error", func() {
				dupTaskA := *taskA
				dupTaskA.Name = "dup-task-a"
				Expect(repo.CreateTask(&dupTaskA)).NotTo(Succeed())
			})
		})
	})

	Describe("Tasks", func() {
		Context("when no tasks exist", func() {
			It("returns no tasks", func() {
				tasks, err := repo.Tasks()
				Expect(err).NotTo(HaveOccurred())
				Expect(tasks).To(HaveLen(0))
			})
		})
		Context("tasks exist", func() {
			BeforeEach(func() {
				Expect(repo.CreateTask(taskA)).To(Succeed())
				Expect(repo.CreateTask(taskB)).To(Succeed())
				Expect(repo.CreateTask(taskC)).To(Succeed())
			})
			It("returns then with Tasks()", func() {
				tasks, err := repo.Tasks()
				Expect(err).NotTo(HaveOccurred())
				Expect(*tasks[0]).To(Equal(*taskA))
				Expect(*tasks[1]).To(Equal(*taskB))
				Expect(*tasks[2]).To(Equal(*taskC))
			})
		})
	})

	Describe("FindTaskByID", func() {
		Context("when the task does not exist", func() {
			BeforeEach(func() {
				Expect(repo.CreateTask(taskA)).To(Succeed())
				Expect(repo.CreateTask(taskC)).To(Succeed())
			})
			It("returns nil and nil error", func() {
				task, err := repo.FindTaskByID(2)
				Expect(err).NotTo(HaveOccurred())
				Expect(task).To(BeNil())
			})
		})

		Context("when the task exists", func() {
			BeforeEach(func() {
				Expect(repo.CreateTask(taskA)).To(Succeed())
				Expect(repo.CreateTask(taskB)).To(Succeed())
			})
			It("returns the task", func() {
				task, err := repo.FindTaskByID(2)
				Expect(err).NotTo(HaveOccurred())
				Expect(*task).To(Equal(*taskB))
			})
		})
	})

	Describe("FindTaskByName", func() {
		Context("when the task does not exist", func() {
			BeforeEach(func() {
				Expect(repo.CreateTask(taskA)).To(Succeed())
				Expect(repo.CreateTask(taskC)).To(Succeed())
			})
			It("returns nil and nil error", func() {
				task, err := repo.FindTaskByName("task-b")
				Expect(err).NotTo(HaveOccurred())
				Expect(task).To(BeNil())
			})
		})

		Context("when the task exists", func() {
			BeforeEach(func() {
				Expect(repo.CreateTask(taskA)).To(Succeed())
				Expect(repo.CreateTask(taskB)).To(Succeed())
			})
			It("returns the task", func() {
				task, err := repo.FindTaskByName("task-b")
				Expect(err).NotTo(HaveOccurred())
				Expect(*task).To(Equal(*taskB))
			})
		})
	})

	Describe("UpdateTask", func() {
		Context("when the task does not exist", func() {
			BeforeEach(func() {
				Expect(repo.CreateTask(taskA)).To(Succeed())
				Expect(repo.CreateTask(taskC)).To(Succeed())
			})
			It("returns an error", func() {
				Expect(repo.UpdateTask(taskB)).NotTo(Succeed())
			})
		})

		Context("when the task exists", func() {
			BeforeEach(func() {
				Expect(repo.CreateTask(taskA)).To(Succeed())
				Expect(repo.CreateTask(taskB)).To(Succeed())
			})
			It("updates the task", func() {
				newTaskB := *taskB
				newTaskB.Name = "new-task-b"
				err := repo.UpdateTask(&newTaskB)
				Expect(err).NotTo(HaveOccurred())

				task, err := repo.FindTaskByName("task-b")
				Expect(err).NotTo(HaveOccurred())
				Expect(task).To(BeNil())

				task, err = repo.FindTaskByName("new-task-b")
				Expect(err).NotTo(HaveOccurred())
				Expect(task).To(Equal(&newTaskB))
			})
			It("does not add a new task", func() {
				newTaskB := *taskB
				newTaskB.Name = "new-task-b"
				err := repo.UpdateTask(&newTaskB)
				Expect(err).NotTo(HaveOccurred())

				tasks, err := repo.Tasks()
				Expect(err).NotTo(HaveOccurred())
				Expect(tasks).To(HaveLen(2))
			})
		})
	})

	Describe("DeleteTask", func() {
		Context("when the task does not exist", func() {
			BeforeEach(func() {
				Expect(repo.CreateTask(taskA)).To(Succeed())
				Expect(repo.CreateTask(taskC)).To(Succeed())
			})
			It("returns an error", func() {
				Expect(repo.DeleteTask(taskB)).NotTo(Succeed())
			})
		})

		Context("when the task exists", func() {
			BeforeEach(func() {
				Expect(repo.CreateTask(taskA)).To(Succeed())
				Expect(repo.CreateTask(taskB)).To(Succeed())
			})
			It("deletes the task", func() {
				err := repo.DeleteTask(taskB)
				Expect(err).NotTo(HaveOccurred())

				task, err := repo.FindTaskByName("task-b")
				Expect(err).NotTo(HaveOccurred())
				Expect(task).To(BeNil())

				task, err = repo.FindTaskByID(2)
				Expect(err).NotTo(HaveOccurred())
				Expect(task).To(BeNil())
			})
		})
	})

	Describe("CreateEvent", func() {
		Context("when events are created", func() {
			BeforeEach(func() {
				Expect(repo.CreateEvent(eventA)).To(Succeed())
				Expect(repo.CreateEvent(eventB)).To(Succeed())
				Expect(repo.CreateEvent(eventC)).To(Succeed())
			})
			It("returns then with Events()", func() {
				events, err := repo.Events()
				Expect(err).NotTo(HaveOccurred())
				Expect(*events[0]).To(Equal(*eventA))
				Expect(*events[1]).To(Equal(*eventB))
				Expect(*events[2]).To(Equal(*eventC))
			})
		})
		Context("when an event with that date already exists", func() {
			BeforeEach(func() {
				Expect(repo.CreateEvent(eventA)).To(Succeed())
				Expect(repo.CreateEvent(eventB)).To(Succeed())
				Expect(repo.CreateEvent(eventC)).To(Succeed())
			})
			It("returns an error", func() {
				dupEventA := *eventA
				dupEventA.Title = "dup-event-a"
				Expect(repo.CreateEvent(&dupEventA)).NotTo(Succeed())
			})
		})
	})

	Describe("Events", func() {
		Context("when no events exist", func() {
			It("returns no events", func() {
				events, err := repo.Events()
				Expect(err).NotTo(HaveOccurred())
				Expect(events).To(HaveLen(0))
			})
		})
		Context("events exist", func() {
			BeforeEach(func() {
				Expect(repo.CreateEvent(eventA)).To(Succeed())
				Expect(repo.CreateEvent(eventB)).To(Succeed())
				Expect(repo.CreateEvent(eventC)).To(Succeed())
			})
			It("returns then with Events()", func() {
				events, err := repo.Events()
				Expect(err).NotTo(HaveOccurred())
				Expect(*events[0]).To(Equal(*eventA))
				Expect(*events[1]).To(Equal(*eventB))
				Expect(*events[2]).To(Equal(*eventC))
			})
		})
	})

	Describe("DeleteEvent", func() {
		Context("when the event does not exist", func() {
			BeforeEach(func() {
				Expect(repo.CreateEvent(eventA)).To(Succeed())
				Expect(repo.CreateEvent(eventC)).To(Succeed())
			})
			It("returns an error", func() {
				Expect(repo.DeleteEvent(eventB)).NotTo(Succeed())
			})
		})

		Context("when the event exists", func() {
			BeforeEach(func() {
				Expect(repo.CreateEvent(eventA)).To(Succeed())
				Expect(repo.CreateEvent(eventB)).To(Succeed())
			})
			It("deletes the event", func() {
				err := repo.DeleteEvent(eventB)
				Expect(err).NotTo(HaveOccurred())

				events, err := repo.Events()
				Expect(err).NotTo(HaveOccurred())
				Expect(events).To(HaveLen(1))
				Expect(*events[0]).To(Equal(*eventA))
			})
		})
	})

	Describe("Persistance", func() {
		Context("when tasks are created with one repo", func() {
			BeforeEach(func() {
				Expect(repo.CreateTask(taskA)).To(Succeed())
				Expect(repo.CreateTask(taskB)).To(Succeed())
				Expect(repo.CreateTask(taskC)).To(Succeed())
			})
			It("returns them from another repo with Tasks()", func() {
				repo2 := createRepoFunc()
				tasks, err := repo2.Tasks()
				Expect(err).NotTo(HaveOccurred())
				Expect(*tasks[0]).To(Equal(*taskA))
				Expect(*tasks[1]).To(Equal(*taskB))
				Expect(*tasks[2]).To(Equal(*taskC))
			})
		})
		Context("when tasks are updated with one repo", func() {
			var newTaskB Task
			BeforeEach(func() {
				Expect(repo.CreateTask(taskA)).To(Succeed())
				Expect(repo.CreateTask(taskB)).To(Succeed())
				Expect(repo.CreateTask(taskC)).To(Succeed())

				newTaskB = *taskB
				newTaskB.Name = "new-task-b"
				Expect(repo.UpdateTask(&newTaskB)).To(Succeed())
			})
			It("returns then from another repo with Tasks()", func() {
				repo2 := createRepoFunc()
				tasks, err := repo2.Tasks()
				Expect(err).NotTo(HaveOccurred())
				Expect(*tasks[0]).To(Equal(*taskA))
				Expect(*tasks[1]).To(Equal(newTaskB))
				Expect(*tasks[2]).To(Equal(*taskC))
			})
		})
		Context("when tasks are deleted with one repo", func() {
			BeforeEach(func() {
				Expect(repo.CreateTask(taskA)).To(Succeed())
				Expect(repo.CreateTask(taskB)).To(Succeed())
				Expect(repo.CreateTask(taskC)).To(Succeed())

				Expect(repo.DeleteTask(taskC)).To(Succeed())
			})
			It("returns then from another repo with Tasks()", func() {
				repo2 := createRepoFunc()
				tasks, err := repo2.Tasks()
				Expect(err).NotTo(HaveOccurred())
				Expect(tasks).To(HaveLen(2))
				Expect(*tasks[0]).To(Equal(*taskA))
				Expect(*tasks[1]).To(Equal(*taskB))
			})
		})
		Context("when events are created with one repo", func() {
			BeforeEach(func() {
				Expect(repo.CreateEvent(eventA)).To(Succeed())
				Expect(repo.CreateEvent(eventB)).To(Succeed())
				Expect(repo.CreateEvent(eventC)).To(Succeed())
			})
			It("returns then from another repo with Events()", func() {
				repo2 := createRepoFunc()
				events, err := repo2.Events()
				Expect(err).NotTo(HaveOccurred())
				Expect(*events[0]).To(Equal(*eventA))
				Expect(*events[1]).To(Equal(*eventB))
				Expect(*events[2]).To(Equal(*eventC))
			})
		})
		Context("when events are deleted with one repo", func() {
			BeforeEach(func() {
				Expect(repo.CreateEvent(eventA)).To(Succeed())
				Expect(repo.CreateEvent(eventB)).To(Succeed())
				Expect(repo.CreateEvent(eventC)).To(Succeed())

				Expect(repo.DeleteEvent(eventC)).To(Succeed())
			})
			It("returns then from another repo with Events()", func() {
				repo2 := createRepoFunc()
				events, err := repo2.Events()
				Expect(err).NotTo(HaveOccurred())
				Expect(events).To(HaveLen(2))
				Expect(*events[0]).To(Equal(*eventA))
				Expect(*events[1]).To(Equal(*eventB))
			})
		})
	})
}
