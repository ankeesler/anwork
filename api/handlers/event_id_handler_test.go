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

var _ = Describe("EventIDHandler", func() {
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

		handler = handlers.NewEventIDHandler(manager, l)
	})

	It("logs that handling is happening", func() {
		handleGet(handler, "/api/v1/events/10")
		Eventually(logWriter).Should(gbytes.Say("Handling GET /api/v1/events.."))
		Eventually(logWriter).Should(gbytes.Say("Getting eventID 10"))
	})

	Context("when the last path segment is bunk", func() {
		It("logs an error", func() {
			handleGet(handler, "/api/v1/events/tuna")
			Eventually(logWriter).Should(gbytes.Say("Unable to parse last path segment"))
		})

		It("returns bad request ", func() {
			rsp := handleGet(handler, "/api/v1/events/tuna")
			Expect(rsp.Code).To(Equal(http.StatusBadRequest))
		})
	})

	XDescribe("GET", func() {
		var t *task.Event
		BeforeEach(func() {
			t = &task.Event{Title: "event-a", TaskID: 5}
			// TODO: setup mock.
		})

		It("returns a JSON object representing the event", func() {
			rsp := handleGet(handler, "/api/v1/events/5")

			Expect(manager.FindByIDArgsForCall(0)).To(Equal(5))

			Expect(rsp.Code).To(Equal(http.StatusOK))
			Expect(rsp.HeaderMap["Content-Type"]).To(Equal([]string{"application/json"}))

			expectedEventJson, err := json.Marshal(t)
			Expect(err).NotTo(HaveOccurred())
			actualEventJson, err := ioutil.ReadAll(rsp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(actualEventJson).To(Equal(expectedEventJson))
		})

		It("logs the response", func() {
			handleGet(handler, "/api/v1/events/5")

			expectedEventJson, err := json.Marshal(t)
			Expect(err).NotTo(HaveOccurred())

			logContents := string(logWriter.Contents())
			Expect(logContents).To(ContainSubstring(fmt.Sprintf("Returning event: %s", expectedEventJson)))
		})

		Context("when there is no event associated with the provided ID", func() {
			BeforeEach(func() {
				manager.FindByIDReturnsOnCall(0, nil)
			})

			It("returns a not found", func() {
				rsp := handleGet(handler, "/api/v1/events/10")

				Expect(rsp.Code).To(Equal(http.StatusNotFound))
			})

			It("logs that it was not able to find the event", func() {
				handleGet(handler, "/api/v1/events/10")
				Eventually(logWriter).Should(gbytes.Say("No event with ID 10"))
			})
		})
	})

	Describe("POST", func() {
		It("responds with method not allowed", func() {
			rsp := handlePost(handler, "/api/v1/events/5", nil)
			Expect(manager.EventsCallCount()).To(Equal(0))
			Expect(rsp.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})

	Describe("PUT", func() {
		It("responds with method not allowed", func() {
			rsp := handlePut(handler, "/api/v1/events/5", nil)
			Expect(manager.EventsCallCount()).To(Equal(0))
			Expect(rsp.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})

	Describe("DELETE", func() {
		It("responds with method not allowed", func() {
			rsp := handleDelete(handler, "/api/v1/events/5")
			Expect(manager.EventsCallCount()).To(Equal(0))
			Expect(rsp.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})

})