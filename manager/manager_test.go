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

	Describe("Delete", func() {
		BeforeEach(func() {
			repo.FindTaskByNameReturnsOnCall(0, &task2.Task{Name: "task-a"}, nil)
		})

		It("find the task and then deletes it", func() {
			Expect(manager.Delete("task-a")).To(Succeed())

			Expect(repo.FindTaskByNameCallCount()).To(Equal(1))
			Expect(repo.FindTaskByNameArgsForCall(0)).To(Equal("task-a"))

			Expect(repo.DeleteTaskCallCount()).To(Equal(1))
			Expect(repo.DeleteTaskArgsForCall(0)).To(Equal(&task2.Task{Name: "task-a"}))
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
})
