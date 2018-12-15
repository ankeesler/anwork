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
		BeforeEach(func() {
			repo.CreateTaskStub = func(task *task2.Task) error {
				task.ID = 10
				return nil
			}
		})

		It("calls down to the repo to create a task", func() {
			Expect(manager.Create("task-a")).To(Succeed())

			Expect(repo.CreateTaskCallCount()).To(Equal(1))
			Expect(*repo.CreateTaskArgsForCall(0)).To(Equal(task2.Task{
				Name:      "task-a",
				StartDate: now.Unix(),
				Priority:  10,               // default
				State:     task2.StateReady, // default
				ID:        10,
			}))
		})

		It("adds an event that the task was created", func() {
			Expect(manager.Create("task-a")).To(Succeed())

			Expect(repo.CreateEventCallCount()).To(Equal(1))
			Expect(repo.CreateEventArgsForCall(0)).To(Equal(&task2.Event{
				Title:  "Created task 'task-a'",
				Date:   now.Unix(),
				Type:   task2.EventTypeCreate,
				TaskID: 10,
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

	Describe("Delete", func() {
		BeforeEach(func() {
			repo.FindTaskByNameReturnsOnCall(0, &task2.Task{Name: "task-a", ID: 10}, nil)
		})

		It("find the task and then deletes it", func() {
			Expect(manager.Delete("task-a")).To(Succeed())

			Expect(repo.FindTaskByNameCallCount()).To(Equal(1))
			Expect(repo.FindTaskByNameArgsForCall(0)).To(Equal("task-a"))

			Expect(repo.DeleteTaskCallCount()).To(Equal(1))
			Expect(repo.DeleteTaskArgsForCall(0)).To(Equal(&task2.Task{
				Name: "task-a",
				ID:   10,
			}))
		})

		It("adds an event that the task was deleted", func() {
			Expect(manager.Delete("task-a")).To(Succeed())

			Expect(repo.CreateEventCallCount()).To(Equal(1))
			Expect(repo.CreateEventArgsForCall(0)).To(Equal(&task2.Event{
				Title:  "Deleted task 'task-a'",
				Date:   now.Unix(),
				Type:   task2.EventTypeDelete,
				TaskID: 10,
			}))
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				repo.FindTaskByNameReturnsOnCall(0, nil, nil)
			})

			It("returns an error", func() {
				err := manager.Delete("task-a")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("unknown task with name 'task-a'"))
			})
		})

		Context("when the repo fails to find the task", func() {
			BeforeEach(func() {
				repo.FindTaskByNameReturnsOnCall(0, nil, errors.New("some find error"))
			})

			It("returns an error", func() {
				err := manager.Delete("task-a")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("some find error"))
			})
		})

		Context("when the repo fails to delete the task", func() {
			BeforeEach(func() {
				repo.DeleteTaskReturnsOnCall(0, errors.New("some delete error"))
			})

			It("returns an error", func() {
				err := manager.Delete("task-a")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("some delete error"))
			})
		})
	})

	Describe("FindByID", func() {
		var task *task2.Task
		BeforeEach(func() {
			task = &task2.Task{Name: "task-a"}
			repo.FindTaskByIDReturnsOnCall(0, task, nil)
		})

		It("calls the repo to find the task", func() {
			t, err := manager.FindByID(10)
			Expect(err).NotTo(HaveOccurred())
			Expect(*t).To(Equal(*task))

			Expect(repo.FindTaskByIDCallCount()).To(Equal(1))
			Expect(repo.FindTaskByIDArgsForCall(0)).To(Equal(10))
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				repo.FindTaskByIDReturnsOnCall(0, nil, nil)
			})

			It("returns nil, nil", func() {
				t, err := manager.FindByID(10)
				Expect(err).NotTo(HaveOccurred())
				Expect(t).To(BeNil())
			})
		})

		Context("when the repo fails to find the task", func() {
			BeforeEach(func() {
				repo.FindTaskByIDReturnsOnCall(0, nil, errors.New("some find error"))
			})

			It("returns the error", func() {
				_, err := manager.FindByID(10)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("some find error"))
			})
		})
	})

	Describe("FindByName", func() {
		var task *task2.Task
		BeforeEach(func() {
			task = &task2.Task{Name: "task-a"}
			repo.FindTaskByNameReturnsOnCall(0, task, nil)
		})

		It("calls the repo to find the task", func() {
			t, err := manager.FindByName("task-a")
			Expect(err).NotTo(HaveOccurred())
			Expect(*t).To(Equal(*task))

			Expect(repo.FindTaskByNameCallCount()).To(Equal(1))
			Expect(repo.FindTaskByNameArgsForCall(0)).To(Equal("task-a"))
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				repo.FindTaskByNameReturnsOnCall(0, nil, nil)
			})

			It("returns nil, nil", func() {
				t, err := manager.FindByName("task-a")
				Expect(err).NotTo(HaveOccurred())
				Expect(t).To(BeNil())
			})
		})

		Context("when the repo fails to find the task", func() {
			BeforeEach(func() {
				repo.FindTaskByNameReturnsOnCall(0, nil, errors.New("some find error"))
			})

			It("returns the error", func() {
				_, err := manager.FindByName("task-a")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("some find error"))
			})
		})
	})

	Describe("Tasks", func() {
		var tasks []*task2.Task
		BeforeEach(func() {
			tasks = []*task2.Task{
				&task2.Task{Name: "task-a"},
				&task2.Task{Name: "task-b"},
				&task2.Task{Name: "task-c"},
			}
			repo.TasksReturnsOnCall(0, tasks, nil)
		})

		It("calls the repo to get the tasks", func() {
			t, err := manager.Tasks()
			Expect(err).NotTo(HaveOccurred())
			Expect(t).To(Equal(tasks))

			Expect(repo.TasksCallCount()).To(Equal(1))
		})

		Context("when the repo fails to get the tasks", func() {
			BeforeEach(func() {
				repo.TasksReturnsOnCall(0, nil, errors.New("some tasks error"))
			})

			It("returns the error", func() {
				_, err := manager.Tasks()
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("some tasks error"))
			})
		})

		Context("when the priorities change", func() {
			BeforeEach(func() {
				tasks = []*task2.Task{
					&task2.Task{Name: "task-a", Priority: 40, ID: 1},
					&task2.Task{Name: "task-b", Priority: 30, ID: 2},
					&task2.Task{Name: "task-c", Priority: 20, ID: 3},
					&task2.Task{Name: "task-d", Priority: 30, ID: 4},
				}
				tasksCopied := make([]*task2.Task, len(tasks))
				copy(tasksCopied, tasks)
				repo.TasksReturnsOnCall(0, tasksCopied, nil)
			})

			It("returns the tasks in order of highest priority and then newest id", func() {
				tasksSorted, err := manager.Tasks()
				Expect(err).NotTo(HaveOccurred())
				Expect(tasksSorted).To(HaveLen(4))
				Expect(tasksSorted[0]).To(Equal(tasks[2]))
				Expect(tasksSorted[1]).To(Equal(tasks[1]))
				Expect(tasksSorted[2]).To(Equal(tasks[3]))
				Expect(tasksSorted[3]).To(Equal(tasks[0]))
			})
		})
	})

	Describe("Note", func() {
		BeforeEach(func() {
			repo.FindTaskByNameReturnsOnCall(0, &task2.Task{Name: "task-a", ID: 10}, nil)
		})

		It("creates an event with type note", func() {
			Expect(manager.Note("task-a", "here is a note")).To(Succeed())

			Expect(repo.FindTaskByNameCallCount()).To(Equal(1))
			Expect(repo.FindTaskByNameArgsForCall(0)).To(Equal("task-a"))

			Expect(repo.CreateEventCallCount()).To(Equal(1))
			Expect(repo.CreateEventArgsForCall(0)).To(Equal(&task2.Event{
				Title:  "Note added to task 'task-a': here is a note",
				Date:   clock.Now().Unix(),
				Type:   task2.EventTypeNote,
				TaskID: 10,
			}))
		})

		Context("the find by name call fails", func() {
			BeforeEach(func() {
				repo.FindTaskByNameReturnsOnCall(0, nil, errors.New("some find by name error"))
			})

			It("returns the error", func() {
				err := manager.Note("task-a", "here is a note")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("some find by name error"))
			})
		})

		Context("the task does not exist", func() {
			BeforeEach(func() {
				repo.FindTaskByNameReturnsOnCall(0, nil, nil)
			})

			It("returns the error", func() {
				err := manager.Note("task-a", "here is a note")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("unknown task with name 'task-a'"))
			})
		})

		Context("the event cannot be created", func() {
			BeforeEach(func() {
				repo.CreateEventReturnsOnCall(0, errors.New("some create event error"))
			})

			It("returns the error", func() {
				err := manager.Note("task-a", "here is a note")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("some create event error"))
			})
		})
	})

	Describe("SetPriority", func() {
		BeforeEach(func() {
			repo.FindTaskByNameReturnsOnCall(0,
				&task2.Task{
					Name:      "task-a",
					ID:        10,
					Priority:  20,
					State:     task2.StateRunning,
					StartDate: 123,
				},
				nil)
		})

		It("updates the task and adds an event saying the priority was updated", func() {
			Expect(manager.SetPriority("task-a", 30)).To(Succeed())

			Expect(repo.FindTaskByNameCallCount()).To(Equal(1))
			Expect(repo.FindTaskByNameArgsForCall(0)).To(Equal("task-a"))

			Expect(repo.CreateEventCallCount()).To(Equal(1))
			Expect(repo.CreateEventArgsForCall(0)).To(Equal(&task2.Event{
				Title:  "Set priority on task 'task-a' from 20 to 30",
				Date:   clock.Now().Unix(),
				Type:   task2.EventTypeSetPriority,
				TaskID: 10,
			}))

			Expect(repo.UpdateTaskCallCount()).To(Equal(1))
			Expect(repo.UpdateTaskArgsForCall(0)).To(Equal(&task2.Task{
				Name:      "task-a",
				ID:        10,
				Priority:  30,
				State:     task2.StateRunning,
				StartDate: 123,
			}))
		})

		Context("the find by name call fails", func() {
			BeforeEach(func() {
				repo.FindTaskByNameReturnsOnCall(0, nil, errors.New("some find by name error"))
			})

			It("returns the error", func() {
				err := manager.SetPriority("task-a", 30)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("some find by name error"))
			})
		})

		Context("the task does not exist", func() {
			BeforeEach(func() {
				repo.FindTaskByNameReturnsOnCall(0, nil, nil)
			})

			It("returns the error", func() {
				err := manager.SetPriority("task-a", 30)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("unknown task with name 'task-a'"))
			})
		})

		Context("the task cannot be updated", func() {
			BeforeEach(func() {
				repo.UpdateTaskReturnsOnCall(0, errors.New("some update task error"))
			})

			It("returns the error", func() {
				err := manager.SetPriority("task-a", 30)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("some update task error"))
			})
		})

		Context("the event cannot be added", func() {
			BeforeEach(func() {
				repo.CreateEventReturnsOnCall(0, errors.New("some create event error"))
			})

			It("returns the error", func() {
				err := manager.SetPriority("task-a", 30)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("some create event error"))
			})
		})
	})

	Describe("SetState", func() {
		BeforeEach(func() {
			repo.FindTaskByNameReturnsOnCall(0,
				&task2.Task{
					Name:      "task-a",
					ID:        10,
					Priority:  20,
					State:     task2.StateRunning,
					StartDate: 123,
				},
				nil)
		})

		It("updates the task and adds an event saying the state was updated", func() {
			Expect(manager.SetState("task-a", task2.StateBlocked)).To(Succeed())

			Expect(repo.FindTaskByNameCallCount()).To(Equal(1))
			Expect(repo.FindTaskByNameArgsForCall(0)).To(Equal("task-a"))

			Expect(repo.CreateEventCallCount()).To(Equal(1))
			Expect(repo.CreateEventArgsForCall(0)).To(Equal(&task2.Event{
				Title:  "Set state on task 'task-a' from Running to Blocked",
				Date:   clock.Now().Unix(),
				Type:   task2.EventTypeSetState,
				TaskID: 10,
			}))

			Expect(repo.UpdateTaskCallCount()).To(Equal(1))
			Expect(repo.UpdateTaskArgsForCall(0)).To(Equal(&task2.Task{
				Name:      "task-a",
				ID:        10,
				Priority:  20,
				State:     task2.StateBlocked,
				StartDate: 123,
			}))
		})

		Context("the find by name call fails", func() {
			BeforeEach(func() {
				repo.FindTaskByNameReturnsOnCall(0, nil, errors.New("some find by name error"))
			})

			It("returns the error", func() {
				err := manager.SetState("task-a", task2.StateBlocked)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("some find by name error"))
			})
		})

		Context("the task does not exist", func() {
			BeforeEach(func() {
				repo.FindTaskByNameReturnsOnCall(0, nil, nil)
			})

			It("returns the error", func() {
				err := manager.SetState("task-a", task2.StateBlocked)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("unknown task with name 'task-a'"))
			})
		})

		Context("the task cannot be updated", func() {
			BeforeEach(func() {
				repo.UpdateTaskReturnsOnCall(0, errors.New("some update task error"))
			})

			It("returns the error", func() {
				err := manager.SetState("task-a", task2.StateBlocked)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("some update task error"))
			})
		})

		Context("the event cannot be added", func() {
			BeforeEach(func() {
				repo.CreateEventReturnsOnCall(0, errors.New("some create event error"))
			})

			It("returns the error", func() {
				err := manager.SetState("task-a", task2.StateBlocked)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("some create event error"))
			})
		})
	})

	Describe("Events", func() {
		var events []*task2.Event
		BeforeEach(func() {
			events = []*task2.Event{
				&task2.Event{Title: "event-a"},
				&task2.Event{Title: "event-b"},
				&task2.Event{Title: "event-c"},
			}
			repo.EventsReturnsOnCall(0, events, nil)
		})

		It("calls the repo to get the events", func() {
			t, err := manager.Events()
			Expect(err).NotTo(HaveOccurred())
			Expect(t).To(Equal(events))

			Expect(repo.EventsCallCount()).To(Equal(1))
		})

		Context("when the repo fails to get the events", func() {
			BeforeEach(func() {
				repo.EventsReturnsOnCall(0, nil, errors.New("some events error"))
			})

			It("returns the error", func() {
				_, err := manager.Events()
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("some events error"))
			})
		})
	})

	Describe("Reset", func() {
		var tasks []*task2.Task
		var tasksSize int
		var events []*task2.Event
		var eventsSize int
		BeforeEach(func() {
			tasks = []*task2.Task{
				&task2.Task{Name: "task-a", ID: 1},
				&task2.Task{Name: "task-b", ID: 2},
				&task2.Task{Name: "task-c", ID: 3},
			}
			tasksSize = len(tasks)
			repo.TasksReturnsOnCall(0, tasks, nil)
			events = []*task2.Event{
				&task2.Event{Title: "task-a", ID: 1},
				&task2.Event{Title: "task-b", ID: 2},
				&task2.Event{Title: "task-c", ID: 3},
			}
			eventsSize = len(events)
			repo.EventsReturnsOnCall(0, events, nil)
		})

		It("deletes all tasks and events", func() {
			Expect(manager.Reset()).To(Succeed())

			Expect(repo.TasksCallCount()).To(Equal(1))
			j := 0
			for i := tasksSize - 1; i >= 0; i-- {
				Expect(repo.DeleteTaskArgsForCall(j)).To(Equal(tasks[i]))
				j++
			}

			Expect(repo.EventsCallCount()).To(Equal(1))
			j = 0
			for i := eventsSize - 1; i >= 0; i-- {
				Expect(repo.DeleteEventArgsForCall(j)).To(Equal(events[i]))
				j++
			}
		})

		Context("getting the tasks fails", func() {
			BeforeEach(func() {
				repo.TasksReturnsOnCall(0, nil, errors.New("some tasks error"))
			})

			It("returns an error before doing anything", func() {
				Expect(manager.Reset()).To(MatchError("some tasks error"))

				Expect(repo.TasksCallCount()).To(Equal(1))
				Expect(repo.DeleteTaskCallCount()).To(Equal(0))

				Expect(repo.EventsCallCount()).To(Equal(0))
				Expect(repo.DeleteEventCallCount()).To(Equal(0))
			})
		})

		Context("getting the events fails", func() {
			BeforeEach(func() {
				repo.EventsReturnsOnCall(0, nil, errors.New("some events error"))
			})

			It("returns an error before doing anything", func() {
				Expect(manager.Reset()).To(MatchError("some events error"))

				Expect(repo.TasksCallCount()).To(Equal(1))
				Expect(repo.DeleteTaskCallCount()).To(Equal(0))

				Expect(repo.EventsCallCount()).To(Equal(1))
				Expect(repo.DeleteEventCallCount()).To(Equal(0))
			})
		})

		Context("at least one of the deletes fails", func() {
			BeforeEach(func() {
				repo.DeleteTaskReturnsOnCall(2, errors.New("some delete task error"))
				repo.DeleteEventReturnsOnCall(1, errors.New("some delete event error"))
			})

			It("does its best to delete the stuff it can and returns an error", func() {
				err := manager.Reset()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("some delete task error"))
				Expect(err.Error()).To(ContainSubstring("some delete event error"))

				Expect(repo.TasksCallCount()).To(Equal(1))
				j := 0
				for i := tasksSize - 1; i >= 0; i-- {
					Expect(repo.DeleteTaskArgsForCall(j)).To(Equal(tasks[i]))
					j++
				}

				Expect(repo.EventsCallCount()).To(Equal(1))
				j = 0
				for i := eventsSize - 1; i >= 0; i-- {
					Expect(repo.DeleteEventArgsForCall(j)).To(Equal(events[i]))
					j++
				}
			})
		})
	})

	Describe("Rename", func() {
		BeforeEach(func() {
			repo.FindTaskByNameReturnsOnCall(0,
				&task2.Task{
					Name:      "task-a",
					ID:        10,
					Priority:  20,
					State:     task2.StateRunning,
					StartDate: 123,
				},
				nil)
		})

		It("updates the task", func() {
			Expect(manager.Rename("task-a", "new-task-a")).To(Succeed())

			Expect(repo.FindTaskByNameCallCount()).To(Equal(1))
			Expect(repo.FindTaskByNameArgsForCall(0)).To(Equal("task-a"))

			Expect(repo.UpdateTaskCallCount()).To(Equal(1))
			Expect(repo.UpdateTaskArgsForCall(0)).To(Equal(&task2.Task{
				Name:      "new-task-a",
				ID:        10,
				Priority:  20,
				State:     task2.StateRunning,
				StartDate: 123,
			}))
		})

		Context("the find by name call fails", func() {
			BeforeEach(func() {
				repo.FindTaskByNameReturnsOnCall(0, nil, errors.New("some find by name error"))
			})

			It("returns the error", func() {
				err := manager.Rename("task-a", "new-task-a")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("some find by name error"))
			})
		})

		Context("the task does not exist", func() {
			BeforeEach(func() {
				repo.FindTaskByNameReturnsOnCall(0, nil, nil)
			})

			It("returns the error", func() {
				err := manager.Rename("task-a", "new-task-a")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("unknown task with name 'task-a'"))
			})
		})

		Context("the task cannot be updated", func() {
			BeforeEach(func() {
				repo.UpdateTaskReturnsOnCall(0, errors.New("some update task error"))
			})

			It("returns the error", func() {
				err := manager.Rename("task-a", "new-task-a")
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("some update task error"))
			})
		})
	})
})
