package api_test

import (
	"errors"
	"log"
	"net/http"
	"os"

	api "github.com/ankeesler/anwork/api2"
	"github.com/ankeesler/anwork/api2/api2fakes"
	"github.com/ankeesler/anwork/task2"
	"github.com/ankeesler/anwork/task2/task2fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/http_server"
)

var _ = Describe("Tasks", func() {
	var (
		repo          *task2fakes.FakeRepo
		authenticator *api2fakes.FakeAuthenticator

		process ifrit.Process
	)

	BeforeEach(func() {
		repo = &task2fakes.FakeRepo{}
		authenticator = &api2fakes.FakeAuthenticator{}

		a := api.New(log.New(GinkgoWriter, "api-test: ", 0), repo, authenticator)
		runner := http_server.New("127.0.0.1:12345", a)
		process = ifrit.Invoke(runner)
	})

	AfterEach(func() {
		process.Signal(os.Kill)
		Eventually(process.Wait()).Should(Receive())
	})

	Describe("Get", func() {
		var tasks []*task2.Task
		BeforeEach(func() {
			tasks = []*task2.Task{
				&task2.Task{Name: "task-a", ID: 1},
				&task2.Task{Name: "task-b", ID: 2},
				&task2.Task{Name: "task-c", ID: 3},
			}
			repo.TasksReturnsOnCall(0, tasks, nil)
		})

		It("responds with the tasks that the repo returns", func() {
			rsp, err := get("/api/v1/tasks")
			Expect(err).NotTo(HaveOccurred())
			defer rsp.Body.Close()

			Expect(rsp.StatusCode).To(Equal(http.StatusOK))
			assertTasks(rsp, tasks)

			Expect(repo.TasksCallCount()).To(Equal(1))
		})

		Context("when getting the tasks fails", func() {
			BeforeEach(func() {
				repo.TasksReturnsOnCall(0, nil, errors.New("some tasks error"))
			})

			It("returns a 500 with an error", func() {
				rsp, err := get("/api/v1/tasks")
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusInternalServerError))
				assertError(rsp, "some tasks error")
			})
		})

		Context("when a query parameter is used", func() {
			Context("when the query parameter is 'name'", func() {
				Context("when the task exists with that name", func() {
					var task *task2.Task
					BeforeEach(func() {
						task = &task2.Task{Name: "task-a", ID: 1}
						repo.FindTaskByNameReturnsOnCall(0, task, nil)
					})

					It("returns an array with that single task", func() {
						rsp, err := get("/api/v1/tasks?name=task-a")
						Expect(err).NotTo(HaveOccurred())
						defer rsp.Body.Close()

						Expect(rsp.StatusCode).To(Equal(http.StatusOK))
						assertTasks(rsp, []*task2.Task{task})

						Expect(repo.FindTaskByNameCallCount()).To(Equal(1))
						Expect(repo.FindTaskByNameArgsForCall(0)).To(Equal("task-a"))
					})
				})

				Context("when a task does not exist with that name", func() {
					It("returns an empty array of tasks", func() {
						rsp, err := get("/api/v1/tasks?name=task-a")
						Expect(err).NotTo(HaveOccurred())
						defer rsp.Body.Close()

						Expect(rsp.StatusCode).To(Equal(http.StatusOK))
						assertTasks(rsp, []*task2.Task{})

						Expect(repo.FindTaskByNameCallCount()).To(Equal(1))
						Expect(repo.FindTaskByNameArgsForCall(0)).To(Equal("task-a"))
					})
				})
			})

			Context("when the query parameter is not 'name'", func() {
				It("ignores it and returns the regular respond", func() {
					rsp, err := get("/api/v1/tasks")
					Expect(err).NotTo(HaveOccurred())
					defer rsp.Body.Close()

					Expect(rsp.StatusCode).To(Equal(http.StatusOK))
					assertTasks(rsp, tasks)

					Expect(repo.TasksCallCount()).To(Equal(1))
				})
			})
		})
	})

	Describe("Post", func() {
		var task *task2.Task
		BeforeEach(func() {
			task = &task2.Task{Name: "task-a", ID: 1}

			repo.CreateTaskStub = func(t *task2.Task) error {
				t.ID = 10
				return nil
			}
		})

		It("creates a task and responds with the location", func() {
			rsp, err := post("/api/v1/tasks", task)
			Expect(err).NotTo(HaveOccurred())
			defer rsp.Body.Close()

			Expect(rsp.StatusCode).To(Equal(http.StatusCreated))
			Expect(rsp.Header.Get("Location")).To(Equal("/api/v1/tasks/10"))

			task.ID = 10
			Expect(repo.CreateTaskCallCount()).To(Equal(1))
			Expect(repo.CreateTaskArgsForCall(0)).To(Equal(task))
		})

		Context("when the request payload is invalid", func() {
			It("responds with a 400 bad request", func() {
				rsp, err := post("/api/v1/tasks", "askjdnflkajnsfd")
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when we fail to create the task", func() {
			BeforeEach(func() {
				repo.CreateTaskReturnsOnCall(0, errors.New("some create error"))
			})

			It("responds with a 500 internal server error", func() {

				rsp, err := post("/api/v1/tasks", task)
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusInternalServerError))
				assertError(rsp, "some create error")
			})
		})
	})
})
