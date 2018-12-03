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
	})
})
