package handlers_test

import (
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
		manager   *taskfakes.FakeManager
		logWriter *gbytes.Buffer
		h         http.Handler
	)

	BeforeEach(func() {
		manager = &taskfakes.FakeManager{}
		logWriter = gbytes.NewBuffer()
		log := log.New(logWriter, "tasks_handler_test.go log: ", 0)
		h = handlers.NewTasksHandler(manager, log)
	})

	It("logs that handling is happening", func() {
		handleGet(h, "/tasks")
		Eventually(logWriter).Should(gbytes.Say("Handling /api/v1/tasks..."))
	})

	Describe("GET", func() {
		BeforeEach(func() {
			manager.TasksReturnsOnCall(0, []*task.Task{
				&task.Task{Name: "task-a", ID: 1},
				&task.Task{Name: "task-b", ID: 2},
				&task.Task{Name: "task-c", ID: 3},
			})
		})

		It("responds with the tasks from the manager", func() {
			rsp := handleGet(h, "/api/v1/tasks")

			Expect(manager.TasksCallCount()).To(Equal(1))

			Expect(rsp.Code).To(Equal(http.StatusOK))
			// TODO: Content-Type is application/json
			// TODO: Body is tasks from above
		})
	})

})
