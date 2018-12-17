package api_test

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/api/apifakes"
	taskpkg "github.com/ankeesler/anwork/task"
	"github.com/ankeesler/anwork/task/taskfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/http_server"
)

var _ = Describe("Task", func() {
	var (
		repo          *taskfakes.FakeRepo
		authenticator *apifakes.FakeAuthenticator

		process ifrit.Process
	)

	testAllCommonFailures := func(doFunc func(path string) (*http.Response, error)) {
		Context("when the id in the path is invalid", func() {
			It("returns with 400 bad request", func() {
				rsp, err := doFunc("/api/v1/tasks/tuna")
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when the repo fails to get the task", func() {
			BeforeEach(func() {
				repo.FindTaskByIDReturnsOnCall(0, nil, errors.New("some find error"))
			})

			It("responds with a 500 internal server error plus the error", func() {
				rsp, err := doFunc("/api/v1/tasks/10")
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusInternalServerError))
				assertError(rsp, "some find error")
			})
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				repo.FindTaskByIDReturnsOnCall(0, nil, nil)
			})

			It("responds with a 404 not found", func() {
				rsp, err := doFunc("/api/v1/tasks/10")
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	}

	BeforeEach(func() {
		repo = &taskfakes.FakeRepo{}
		authenticator = &apifakes.FakeAuthenticator{}

		a := api.New(log.New(GinkgoWriter, "api-test: ", 0), repo, authenticator)
		runner := http_server.New("127.0.0.1:12345", a)
		process = ifrit.Invoke(runner)
	})

	AfterEach(func() {
		process.Signal(os.Kill)
		Eventually(process.Wait()).Should(Receive())
	})

	Describe("Get", func() {
		var task *taskpkg.Task
		BeforeEach(func() {
			task = &taskpkg.Task{Name: "task-a", ID: 1}
			repo.FindTaskByIDReturnsOnCall(0, task, nil)
		})

		It("responds with the task", func() {
			rsp, err := get("/api/v1/tasks/10")
			Expect(err).NotTo(HaveOccurred())
			defer rsp.Body.Close()

			Expect(rsp.StatusCode).To(Equal(http.StatusOK))
			assertTask(rsp, task)

			Expect(repo.FindTaskByIDCallCount()).To(Equal(1))
			Expect(repo.FindTaskByIDArgsForCall(0)).To(Equal(10))
		})

		testAllCommonFailures(get)
	})

	Describe("Put", func() {
		var task *taskpkg.Task
		BeforeEach(func() {
			task = &taskpkg.Task{Name: "task-a", ID: 10}
			repo.FindTaskByIDReturnsOnCall(0, task, nil)
		})

		It("updates the task", func() {
			newTask := *task
			newTask.Name = "new-task-a"
			rsp, err := put("/api/v1/tasks/10", newTask)
			Expect(err).NotTo(HaveOccurred())
			defer rsp.Body.Close()

			Expect(rsp.StatusCode).To(Equal(http.StatusNoContent))

			Expect(repo.FindTaskByIDCallCount()).To(Equal(1))
			Expect(repo.FindTaskByIDArgsForCall(0)).To(Equal(10))

			Expect(repo.UpdateTaskCallCount()).To(Equal(1))
			Expect(repo.UpdateTaskArgsForCall(0)).To(Equal(&newTask))
		})

		Context("when the request body is invalid", func() {
			It("returns bad request", func() {
				rsp, err := put("/api/v1/tasks/10", "asdf")
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when the repo fails to update the task", func() {
			BeforeEach(func() {
				repo.UpdateTaskReturnsOnCall(0, errors.New("some update failure"))
			})

			It("responds with a 500 and the error", func() {
				rsp, err := put("/api/v1/tasks/10", task)
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusInternalServerError))
				assertError(rsp, "some update failure")
			})
		})

		testAllCommonFailures(func(path string) (*http.Response, error) {
			return put(path, task)
		})
	})

	Describe("Delete", func() {
		var task *taskpkg.Task
		BeforeEach(func() {
			task = &taskpkg.Task{Name: "task-a", ID: 1}
			repo.FindTaskByIDReturnsOnCall(0, task, nil)
		})

		It("deletes the task", func() {
			rsp, err := deletee("/api/v1/tasks/10")
			Expect(err).NotTo(HaveOccurred())
			defer rsp.Body.Close()

			Expect(rsp.StatusCode).To(Equal(http.StatusNoContent))

			Expect(repo.FindTaskByIDCallCount()).To(Equal(1))
			Expect(repo.FindTaskByIDArgsForCall(0)).To(Equal(10))

			Expect(repo.DeleteTaskCallCount()).To(Equal(1))
			Expect(repo.DeleteTaskArgsForCall(0)).To(Equal(task))
		})

		Context("when the repo fails to delete the task", func() {
			BeforeEach(func() {
				repo.DeleteTaskReturnsOnCall(0, errors.New("some delete failure"))
			})

			It("responds with a 500 and the error", func() {
				rsp, err := deletee("/api/v1/tasks/10")
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusInternalServerError))
				assertError(rsp, "some delete failure")
			})
		})

		testAllCommonFailures(deletee)
	})
})
