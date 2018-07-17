package api_test

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/ankeesler/anwork/api"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("HealthHandler", func() {
	var (
		logWriter *gbytes.Buffer

		handler http.Handler
	)

	BeforeEach(func() {
		logWriter = gbytes.NewBuffer()
		l := log.New(io.MultiWriter(logWriter, GinkgoWriter), "api_test.go log: ", log.Ldate|log.Ltime|log.Lshortfile)

		handler = api.NewNavHandler(l)
	})

	It("logs that handling is happening", func() {
		handleGet(handler, "/api")
		Eventually(logWriter).Should(gbytes.Say("Handling GET /api..."))
	})

	Describe("GET", func() {
		It("returns an OK status with a json map of links", func() {
			rsp := handleGet(handler, "/api")

			Expect(rsp.Code).To(Equal(http.StatusOK))

			v, ok := rsp.HeaderMap["Content-Type"]
			Expect(ok).To(BeTrue(), "Did not find 'Content-Type' header")
			Expect(v).To(Equal([]string{"application/json"}))

			var payload map[string]map[string]string
			decoder := json.NewDecoder(rsp.Body)
			Expect(decoder.Decode(&payload)).To(Succeed())

			links, ok := payload["links"]
			Expect(ok).To(BeTrue(), "Expected to find 'links' key in returned JSON")

			tasks, ok := links["tasks"]
			Expect(ok).To(BeTrue(), "Expected to find 'tasks' key in returned JSON")
			Expect(tasks).To(Equal("/api/v1/tasks"))
			events, ok := links["events"]
			Expect(ok).To(BeTrue(), "Expected to find 'events' key in returned JSON")
			Expect(events).To(Equal("/api/v1/events"))
			health, ok := links["health"]
			Expect(ok).To(BeTrue(), "Expected to find 'health' key in returned JSON")
			Expect(health).To(Equal("/api/v1/health"))
		})
	})

	Describe("POST", func() {
		It("responds with method not allowed", func() {
			rsp := handlePost(handler, "/api", nil)
			Expect(rsp.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})

	Describe("PUT", func() {
		It("responds with method not allowed", func() {
			rsp := handlePut(handler, "/api", nil)
			Expect(rsp.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})

	Describe("DELETE", func() {
		It("responds with method not allowed", func() {
			rsp := handleDelete(handler, "/api")
			Expect(rsp.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})

})
