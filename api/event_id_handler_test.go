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

		handler = api.NewEventIDHandler(manager, l)
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

	Describe("GET", func() {
		var e *task.Event
		BeforeEach(func() {
			e = &task.Event{Title: "event-a", Date: 123}
			events := []*task.Event{
				e,
				&task.Event{Title: "event-b", Date: 456},
				&task.Event{Title: "event-c", Date: 789},
			}
			manager.EventsReturnsOnCall(0, events)
		})

		It("returns a JSON object representing the event", func() {
			rsp := handleGet(handler, "/api/v1/events/123")

			Expect(manager.EventsCallCount()).To(Equal(1))

			Expect(rsp.Code).To(Equal(http.StatusOK))
			Expect(rsp.HeaderMap["Content-Type"]).To(Equal([]string{"application/json"}))

			expectedEventJson, err := json.Marshal(e)
			Expect(err).NotTo(HaveOccurred())
			actualEventJson, err := ioutil.ReadAll(rsp.Body)
			Expect(err).NotTo(HaveOccurred())
			Expect(actualEventJson).To(Equal(expectedEventJson))
		})

		It("logs the response", func() {
			handleGet(handler, "/api/v1/events/123")

			expectedEventJson, err := json.Marshal(e)
			Expect(err).NotTo(HaveOccurred())

			logContents := string(logWriter.Contents())
			Expect(logContents).To(ContainSubstring(fmt.Sprintf("Returning event: %s", expectedEventJson)))
		})

		Context("when there is no event associated with the provided ID", func() {
			BeforeEach(func() {
				events := []*task.Event{
					&task.Event{Title: "event-b", Date: 456},
					&task.Event{Title: "event-c", Date: 789},
				}
				manager.EventsReturnsOnCall(0, events)
			})

			It("returns a not found", func() {
				rsp := handleGet(handler, "/api/v1/events/123")

				Expect(rsp.Code).To(Equal(http.StatusNotFound))
			})

			It("logs that it was not able to find the event", func() {
				handleGet(handler, "/api/v1/events/123")
				Eventually(logWriter).Should(gbytes.Say("No event with ID 123"))
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
		BeforeEach(func() {
			manager.EventsReturnsOnCall(0, []*task.Event{
				&task.Event{Date: 1},
				&task.Event{Date: 3},
				&task.Event{Date: 5},
			})
		})

		It("deletes the event via the manager", func() {
			rsp := handleDelete(handler, "/api/v1/events/5")

			Expect(manager.DeleteEventCallCount()).To(Equal(1))
			Expect(manager.DeleteEventArgsForCall(0)).To(Equal(int64(5)))

			Expect(rsp.Code).To(Equal(http.StatusNoContent))
		})

		It("logs a message saying it is deleting an event", func() {
			handleDelete(handler, "/api/v1/events/5")

			Eventually(logWriter).Should(gbytes.Say("Deleting event with start time 5"))
		})

		Context("when the manager fails to delete the event", func() {
			BeforeEach(func() {
				manager.DeleteEventReturnsOnCall(0, errors.New("failed to delete event"))
			})

			It("returns the error via an internal server error status", func() {
				rsp := handleDelete(handler, "/api/v1/events/5")

				Expect(manager.DeleteEventCallCount()).To(Equal(1))
				Expect(manager.DeleteEventArgsForCall(0)).To(Equal(int64(5)))

				Expect(rsp.Code).To(Equal(http.StatusInternalServerError))

				var errRsp api.ErrorResponse
				decoder := json.NewDecoder(rsp.Body)
				Expect(decoder.Decode(&errRsp)).To(Succeed())
				Expect(errRsp).To(Equal(api.ErrorResponse{Message: "failed to delete event"}))
			})

			It("logs the error", func() {
				handleDelete(handler, "/api/v1/events/5")

				Eventually(logWriter).Should(gbytes.Say("failed to delete event"))
			})
		})
	})

})
