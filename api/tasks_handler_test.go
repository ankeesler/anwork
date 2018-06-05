package api_test

import (
	"log"

	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/task/taskfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("/api/v1/tasks", func() {
	var (
		manager *taskfakes.FakeManager
		a       *api.Api

		logWriter *gbytes.Buffer
		errChan   chan error
	)

	BeforeEach(func() {
		logWriter = gbytes.NewBuffer()
		l := log.New(logWriter, "tasks_handler_test.go log: ", 0)
		a = api.New(manager, l)
		errChan = make(chan error)
		a.Start(address, errChan)
	})

	AfterEach(func() {
		Expect(errChan).To(BeEmpty())
		a.Stop()
	})

	It("logs that handling is happening", func() {
		_, err := get("/api/v1/tasks")
		Expect(err).NotTo(HaveOccurred())
		Eventually(logWriter).Should(gbytes.Say("Handling /api/v1/tasks..."))
	})

	Describe("GET", func() {
		It("responds with the tasks from the manager", func() {
			_, err := get("/api/v1/tasks")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
