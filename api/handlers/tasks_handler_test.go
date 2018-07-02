package handlers_test

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ankeesler/anwork/api/handlers"
	"github.com/ankeesler/anwork/task"
	"github.com/ankeesler/anwork/task/taskfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("TasksHandler", func() {
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

		handler = handlers.NewTasksHandler(manager, l)
	})

	It("logs that handling is happening", func() {
		handleGet(handler, "/api/v1/tasks")
		Eventually(logWriter).Should(gbytes.Say("Handling GET /api/v1/tasks..."))
	})

	Describe("GET", func() {
		var tasks []*task.Task
		BeforeEach(func() {
			tasks = []*task.Task{
				&task.Task{Name: "task-a", ID: 1},
				&task.Task{Name: "task-b", ID: 2},
				&task.Task{Name: "task-c", ID: 3},
			}
			manager.TasksReturnsOnCall(0, tasks)
		})

		It("responds with the tasks from the manager", func() {
			rsp := handleGet(handler, "/api/v1/tasks")

			Expect(manager.TasksCallCount()).To(Equal(1))

			Expect(rsp.Code).To(Equal(http.StatusOK))
			Expect(rsp.HeaderMap["Content-Type"]).To(Equal([]string{"application/json"}))

			expectedTasksJson, err := json.Marshal(tasks)
			Expect(err).NotTo(HaveOccurred())
			actualTasksJson, err := ioutil.ReadAll(rsp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(actualTasksJson).To(Equal(expectedTasksJson))
		})

		It("logs the tasks that it is returning", func() {
			handleGet(handler, "/api/v1/tasks")

			expectedTasksJson, err := json.Marshal(tasks)
			Expect(err).NotTo(HaveOccurred())

			logContents := string(logWriter.Contents())
			Expect(logContents).To(ContainSubstring(fmt.Sprintf("Returning tasks %s", expectedTasksJson)))
		})
	})

	Describe("POST", func() {
		It("responds with method not allowed", func() {
			rsp := handlePost(handler, "/api/v1/tasks", nil)
			Expect(manager.TasksCallCount()).To(Equal(0))
			Expect(rsp.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})

	Describe("PUT", func() {
		It("responds with method not allowed", func() {
			rsp := handlePut(handler, "/api/v1/tasks", nil)
			Expect(manager.TasksCallCount()).To(Equal(0))
			Expect(rsp.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})

	Describe("DELETE", func() {
		It("responds with method not allowed", func() {
			rsp := handleDelete(handler, "/api/v1/tasks")
			Expect(manager.TasksCallCount()).To(Equal(0))
			Expect(rsp.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})

})
