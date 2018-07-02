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

var _ = Describe("EventsHandler", func() {
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

		handler = handlers.NewEventsHandler(manager, l)
	})

	It("logs that handling is happening", func() {
		handleGet(handler, "/api/v1/events")
		Eventually(logWriter).Should(gbytes.Say("Handling GET /api/v1/events..."))
	})

	Describe("GET", func() {
		var events []*task.Event
		BeforeEach(func() {
			events = []*task.Event{
				&task.Event{Title: "event-a", TaskID: 1},
				&task.Event{Title: "event-b", TaskID: 2},
				&task.Event{Title: "event-c", TaskID: 3},
			}
			manager.EventsReturnsOnCall(0, events)
		})

		It("responds with the events from the manager", func() {
			rsp := handleGet(handler, "/api/v1/events")

			Expect(manager.EventsCallCount()).To(Equal(1))

			Expect(rsp.Code).To(Equal(http.StatusOK))
			Expect(rsp.HeaderMap["Content-Type"]).To(Equal([]string{"application/json"}))

			expectedEventsJson, err := json.Marshal(events)
			Expect(err).NotTo(HaveOccurred())
			actualEventsJson, err := ioutil.ReadAll(rsp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(actualEventsJson).To(Equal(expectedEventsJson))
		})

		It("logs the events that it is returning", func() {
			handleGet(handler, "/api/v1/events")

			expectedEventsJson, err := json.Marshal(events)
			Expect(err).NotTo(HaveOccurred())

			logContents := string(logWriter.Contents())
			Expect(logContents).To(ContainSubstring(fmt.Sprintf("Returning events %s", expectedEventsJson)))
		})
	})

	Describe("POST", func() {
		It("responds with method not allowed", func() {
			rsp := handlePost(handler, "/api/v1/events", nil)
			Expect(manager.EventsCallCount()).To(Equal(0))
			Expect(rsp.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})

	Describe("PUT", func() {
		It("responds with method not allowed", func() {
			rsp := handlePut(handler, "/api/v1/events", nil)
			Expect(manager.EventsCallCount()).To(Equal(0))
			Expect(rsp.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})

	Describe("DELETE", func() {
		It("responds with method not allowed", func() {
			rsp := handleDelete(handler, "/api/v1/events")
			Expect(manager.EventsCallCount()).To(Equal(0))
			Expect(rsp.Code).To(Equal(http.StatusMethodNotAllowed))
		})
	})

})
