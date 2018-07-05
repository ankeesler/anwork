package api_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/task"
	"github.com/ankeesler/anwork/task/taskfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("TaskIDHandler", func() {
	var (
		manager *taskfakes.FakeManager

		logWriter *gbytes.Buffer

		handler http.Handler
	)

	BeforeEach(func() {
		factory := &taskfakes.FakeManagerFactory{}
		manager = &taskfakes.FakeManager{}
		factory.CreateReturnsOnCall(0, manager, nil)

		logWriter = gbytes.NewBuffer()
		l := log.New(io.MultiWriter(logWriter, GinkgoWriter), "api_test.go log: ", log.Ldate|log.Ltime|log.Lshortfile)

		handler = api.NewTaskIDHandler(manager, l)
	})

	It("logs that handling is happening", func() {
		handleGet(handler, "/api/v1/tasks/10")
		Eventually(logWriter).Should(gbytes.Say("Handling GET /api/v1/tasks.."))
		Eventually(logWriter).Should(gbytes.Say("Getting taskID 10"))
	})

	Context("when the last path segment is bunk", func() {
		It("logs an error", func() {
			handleGet(handler, "/api/v1/tasks/tuna")
			Eventually(logWriter).Should(gbytes.Say("Unable to parse last path segment"))
		})

		It("returns bad request ", func() {
			rsp := handleGet(handler, "/api/v1/tasks/tuna")
			Expect(rsp.Code).To(Equal(http.StatusBadRequest))
		})
	})

	Describe("GET", func() {
		var t *task.Task
		BeforeEach(func() {
			t = &task.Task{Name: "task-a", ID: 5}
			manager.FindByIDReturnsOnCall(0, t)
		})

		It("returns a JSON object representing the task", func() {
			rsp := handleGet(handler, "/api/v1/tasks/5")

			Expect(manager.FindByIDArgsForCall(0)).To(Equal(5))

			Expect(rsp.Code).To(Equal(http.StatusOK))
			Expect(rsp.HeaderMap["Content-Type"]).To(Equal([]string{"application/json"}))

			expectedTaskJson, err := json.Marshal(t)
			Expect(err).NotTo(HaveOccurred())
			actualTaskJson, err := ioutil.ReadAll(rsp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(actualTaskJson).To(Equal(expectedTaskJson))
		})

		It("logs the response", func() {
			handleGet(handler, "/api/v1/tasks/5")

			expectedTaskJson, err := json.Marshal(t)
			Expect(err).NotTo(HaveOccurred())

			logContents := string(logWriter.Contents())
			Expect(logContents).To(ContainSubstring(fmt.Sprintf("Returning task: %s", expectedTaskJson)))
		})

		Context("when there is no task associated with the provided ID", func() {
			BeforeEach(func() {
				manager.FindByIDReturnsOnCall(0, nil)
			})

			It("returns a not found", func() {
				rsp := handleGet(handler, "/api/v1/tasks/10")

				Expect(rsp.Code).To(Equal(http.StatusNotFound))
			})

			It("logs that it was not able to find the task", func() {
				handleGet(handler, "/api/v1/tasks/10")
				Eventually(logWriter).Should(gbytes.Say("No task with ID 10"))
			})
		})
	})

	Describe("POST", func() {
		It("responds with method not allowed", func() {
			rsp := handlePost(handler, "/api/v1/tasks/5", nil)
			Expect(manager.TasksCallCount()).To(Equal(0))
			Expect(rsp.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})

	Describe("PUT", func() {
		It("responds with method not allowed", func() {
			rsp := handlePut(handler, "/api/v1/tasks/5", nil)
			Expect(manager.TasksCallCount()).To(Equal(0))
			Expect(rsp.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})

	Describe("DELETE", func() {
		BeforeEach(func() {
			manager.FindByIDReturnsOnCall(0, &task.Task{Name: "task-a"})
			manager.DeleteReturnsOnCall(0, nil)
		})

		It("responds with no content", func() {
			rsp := handleDelete(handler, "/api/v1/tasks/5")

			Expect(manager.FindByIDCallCount()).To(Equal(1))
			Expect(manager.FindByIDArgsForCall(0)).To(Equal(5))
			Expect(manager.DeleteCallCount()).To(Equal(1))
			Expect(manager.DeleteArgsForCall(0)).To(Equal("task-a"))

			Expect(rsp.Code).To(Equal(http.StatusNoContent))
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				manager.FindByIDReturnsOnCall(0, nil)
			})

			It("responds with not found", func() {
				rsp := handleDelete(handler, "/api/v1/tasks/5")

				Expect(rsp.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("when the delete operation fails", func() {
			BeforeEach(func() {
				manager.FindByIDReturnsOnCall(0, &task.Task{Name: "task-a"})
				manager.DeleteReturnsOnCall(0, errors.New("some delete error"))
			})

			PIt("responds with a server error", func() {
				rsp := handleDelete(handler, "/api/v1/tasks/5")

				Expect(rsp.Code).To(Equal(http.StatusInternalServerError))
				Expect(rsp.HeaderMap["Content-Type"]).To(Equal([]string{"application/json"}))

				expectedErrJson, err := json.Marshal(api.ErrorResponse{Message: "some delete error"})
				Expect(err).NotTo(HaveOccurred())
				errJson, err := ioutil.ReadAll(rsp.Body)
				Expect(err).NotTo(HaveOccurred())
				Expect(errJson).To(Equal(expectedErrJson))
			})
		})
	})
})
