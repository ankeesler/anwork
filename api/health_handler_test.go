package api_test

import (
	"io"
	"log"
	"net/http"

	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/task/taskfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("HealthHandler", func() {
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

		handler = api.NewHealthHandler(manager, l)
	})

	It("logs that handling is happening", func() {
		handleGet(handler, "/api/v1/health")
		Eventually(logWriter).Should(gbytes.Say("Handling GET /api/v1/health..."))
	})

	Describe("GET", func() {
		It("returns a no content success status", func() {
			rsp := handleGet(handler, "/api/v1/health")

			Expect(rsp.Code).To(Equal(http.StatusNoContent))
		})

		It("logs the response", func() {
			handleGet(handler, "/api/v1/health")

			logContents := string(logWriter.Contents())
			Expect(logContents).To(ContainSubstring("Returning healthy..."))
		})
	})

	Describe("POST", func() {
		It("responds with method not allowed", func() {
			rsp := handlePost(handler, "/api/v1/health", nil)
			Expect(manager.EventsCallCount()).To(Equal(0))
			Expect(rsp.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})

	Describe("PUT", func() {
		It("responds with method not allowed", func() {
			rsp := handlePut(handler, "/api/v1/health", nil)
			Expect(manager.EventsCallCount()).To(Equal(0))
			Expect(rsp.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})

	Describe("DELETE", func() {
		It("responds with method not allowed", func() {
			rsp := handleDelete(handler, "/api/v1/health")
			Expect(manager.EventsCallCount()).To(Equal(0))
			Expect(rsp.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})

})