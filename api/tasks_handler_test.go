package api_test

import (
	"context"
	"io"
	"log"
	"net/http"

	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/task"
	"github.com/ankeesler/anwork/task/taskfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = XDescribe("TasksHandler", func() {
	var (
		manager *taskfakes.FakeManager

		logWriter *gbytes.Buffer

		ctx, cancel = context.WithCancel(context.Background())
	)

	BeforeEach(func() {
		factory := &taskfakes.FakeManagerFactory{}
		manager = &taskfakes.FakeManager{}
		factory.CreateReturnsOnCall(0, manager, nil)

		logWriter = gbytes.NewBuffer()
		l := log.New(io.MultiWriter(logWriter, GinkgoWriter), "api_test.go log: ", log.Ldate|log.Ltime|log.Lshortfile)

		a := api.New(address, factory, l)
		Expect(a.Run(ctx)).To(Succeed())
	})

	AfterEach(func() {
		cancel()
		Eventually(logWriter).Should(gbytes.Say("listener closed"))
	})

	It("logs that handling is happening", func() {
		_, err := get("/api/v1/tasks")
		Expect(err).NotTo(HaveOccurred())
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
			rsp, err := get("/api/v1/tasks")
			Expect(err).NotTo(HaveOccurred())

			Expect(manager.TasksCallCount()).To(Equal(1))

			Expect(rsp.StatusCode).To(Equal(http.StatusOK))
			// TODO: Content-Type is application/json
			// TODO: Body is tasks from above
		})
	})

})
