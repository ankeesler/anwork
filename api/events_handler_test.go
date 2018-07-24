package api_test

import (
	"bytes"
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

		handler = api.NewEventsHandler(manager, l)
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
		BeforeEach(func() {
			manager.FindByIDReturnsOnCall(0, &task.Task{
				Name: "task-a",
				ID:   1,
			})
		})

		It("creates an event if the event is a note", func() {
			payload, err := json.Marshal(api.AddEventRequest{
				Title:  "event-a",
				Date:   12345,
				Type:   task.EventTypeNote,
				TaskID: 1,
			})
			Expect(err).NotTo(HaveOccurred())
			payloadBuffer := bytes.NewBuffer(payload)

			rsp := handlePost(handler, "/api/v1/events", payloadBuffer)

			Expect(manager.FindByIDCallCount()).To(Equal(1))
			Expect(manager.NoteCallCount()).To(Equal(1))
			name, note := manager.NoteArgsForCall(0)
			Expect(name).To(Equal("task-a"))
			Expect(note).To(Equal("event-a"))

			Expect(rsp.Code).To(Equal(http.StatusNoContent))
		})

		It("returns bad request if the event is not a note", func() {
			payload, err := json.Marshal(api.AddEventRequest{
				Title:  "event-a",
				Date:   12345,
				Type:   task.EventTypeCreate,
				TaskID: 1,
			})
			Expect(err).NotTo(HaveOccurred())
			payloadBuffer := bytes.NewBuffer(payload)

			rsp := handlePost(handler, "/api/v1/events", payloadBuffer)

			Expect(manager.FindByIDCallCount()).To(Equal(0))
			Expect(manager.NoteCallCount()).To(Equal(0))

			Expect(rsp.Code).To(Equal(http.StatusBadRequest))

			var errRsp api.ErrorResponse
			decoder := json.NewDecoder(rsp.Body)
			Expect(decoder.Decode(&errRsp)).To(Succeed())
			Expect(errRsp.Message).To(Equal("Invalid event type 0, the only supported event type is 3"))
		})

		Context("when the payload is bogus", func() {
			It("returns bad request", func() {
				payloadBuffer := bytes.NewBuffer([]byte("tuna"))

				rsp := handlePost(handler, "/api/v1/events", payloadBuffer)

				Expect(rsp.Code).To(Equal(http.StatusBadRequest))

				var errRsp api.ErrorResponse
				decoder := json.NewDecoder(rsp.Body)
				Expect(decoder.Decode(&errRsp)).To(Succeed())
				Expect(errRsp.Message).To(Equal("Invalid request payload: tuna"))
			})
		})

		Context("when there is no task with the ID in the event", func() {
			BeforeEach(func() {
				manager.FindByIDReturnsOnCall(0, nil)
			})

			It("returns bad request", func() {
				payload, err := json.Marshal(api.AddEventRequest{
					Title:  "event-a",
					Date:   12345,
					Type:   task.EventTypeNote,
					TaskID: 1,
				})
				Expect(err).NotTo(HaveOccurred())
				payloadBuffer := bytes.NewBuffer(payload)

				rsp := handlePost(handler, "/api/v1/events", payloadBuffer)

				Expect(manager.FindByIDCallCount()).To(Equal(1))
				Expect(manager.NoteCallCount()).To(Equal(0))

				Expect(rsp.Code).To(Equal(http.StatusBadRequest))

				var errRsp api.ErrorResponse
				decoder := json.NewDecoder(rsp.Body)
				Expect(decoder.Decode(&errRsp)).To(Succeed())
				Expect(errRsp.Message).To(Equal("Unknown task for ID 1"))
			})
		})

		Context("when the manager fails to write the note", func() {
			BeforeEach(func() {
				manager.NoteReturnsOnCall(0, errors.New("some note error"))
			})

			It("returns internal server error", func() {
				payload, err := json.Marshal(api.AddEventRequest{
					Title:  "event-a",
					Date:   12345,
					Type:   task.EventTypeNote,
					TaskID: 1,
				})
				Expect(err).NotTo(HaveOccurred())
				payloadBuffer := bytes.NewBuffer(payload)

				rsp := handlePost(handler, "/api/v1/events", payloadBuffer)

				Expect(rsp.Code).To(Equal(http.StatusBadRequest))

				var errRsp api.ErrorResponse
				decoder := json.NewDecoder(rsp.Body)
				Expect(decoder.Decode(&errRsp)).To(Succeed())
				Expect(errRsp.Message).To(Equal("Failed to add note: some note error"))
			})
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
