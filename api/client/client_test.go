// TODO: run the client tests against the real API!
package client_test

import (
	"net/http"

	"github.com/ankeesler/anwork/api"
	apiclient "github.com/ankeesler/anwork/api/client"
	"github.com/ankeesler/anwork/task"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Manager", func() {
	var (
		server *ghttp.Server
		client *apiclient.Client
	)

	BeforeEach(func() {
		server = ghttp.NewServer()
		server.Writer = GinkgoWriter

		client = apiclient.New(server.URL())
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("CreateTask", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/v1/tasks"),
				ghttp.VerifyJSONRepresenting(api.CreateRequest{Name: "a"}),
				ghttp.VerifyHeaderKV("Content-Type", "application/json"),
				ghttp.RespondWith(http.StatusCreated, nil),
			))
		})

		It("successfully can create tasks via a POST to /api/v1/tasks", func() {
			Expect(client.CreateTask("a")).To(Succeed())
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when the request returns a failure", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/api/v1/tasks"),
					ghttp.VerifyJSONRepresenting(api.CreateRequest{Name: "a"}),
					ghttp.VerifyHeaderKV("Content-Type", "application/json"),
					ghttp.RespondWithJSONEncoded(
						http.StatusBadRequest,
						api.ErrorResponse{Message: "failed to create task"},
					),
				))
			})

			It("prints the failure message", func() {
				Expect(client.CreateTask("a")).To(MatchError("failed to create task"))
			})
		})

		Context("when the request returns a failure with no payload", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/api/v1/tasks"),
					ghttp.VerifyJSONRepresenting(api.CreateRequest{Name: "a"}),
					ghttp.VerifyHeaderKV("Content-Type", "application/json"),
					ghttp.RespondWith(http.StatusInternalServerError, nil),
				))
			})

			It("prints the failure message", func() {
				Expect(client.CreateTask("a")).To(MatchError("Unexpected response payload: "))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("returns an error", func() {
				Expect(client.CreateTask("a")).NotTo(Succeed())
			})
		})
	})

	Describe("DeleteTask", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/api/v1/tasks/1"),
				ghttp.RespondWith(http.StatusNoContent, nil),
			))
		})

		It("successfully can delete tasks via a DELETE to /api/v1/tasks/:id", func() {
			Expect(client.DeleteTask(1)).To(Succeed())
			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when the request returns no payload", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/api/v1/tasks/1"),
					ghttp.RespondWith(http.StatusBadRequest, nil),
				))
			})

			It("prints the failure message", func() {
				err := client.DeleteTask(1)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unexpected response payload: "))
			})
		})

		Context("when the request returns a failure", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/api/v1/tasks/1"),
					ghttp.RespondWithJSONEncoded(
						http.StatusBadRequest,
						api.ErrorResponse{Message: "failed to delete task"},
					),
				))
			})

			It("prints the failure message", func() {
				Expect(client.DeleteTask(1)).To(MatchError("failed to delete task"))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("returns an error", func() {
				Expect(client.DeleteTask(1)).NotTo(Succeed())
			})
		})
	})

	Describe("GetTasks", func() {
		var tasks []*task.Task

		BeforeEach(func() {
			tasks = []*task.Task{
				&task.Task{Name: "task-a", ID: 1},
				&task.Task{Name: "task-b", ID: 2},
				&task.Task{Name: "task-c", ID: 3},
			}
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/tasks"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.RespondWithJSONEncoded(http.StatusOK, tasks),
			))
		})

		It("returns the tasks via a GET to /api/v1/tasks endpoint", func() {
			actualTasks, err := client.GetTasks()
			Expect(err).NotTo(HaveOccurred())
			Expect(actualTasks).ToNot(BeNil())
			Expect(actualTasks).To(Equal(tasks))

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when the response has a weird payload", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/tasks"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.RespondWith(http.StatusOK, nil),
				))
			})

			It("returns an error", func() {
				_, err := client.GetTasks()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unexpected response: "))
			})
		})

		Context("when the response has a weird status", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/tasks"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.RespondWith(http.StatusInternalServerError, nil),
				))
			})

			It("returns an error", func() {
				_, err := client.GetTasks()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unexpected response status: 500 Internal Server Error"))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("returns an error", func() {
				_, err := client.GetTasks()
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("GetTask", func() {
		var expectedTask *task.Task

		BeforeEach(func() {
			expectedTask = &task.Task{Name: "task-a", ID: 1}
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/tasks/1"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.RespondWithJSONEncoded(http.StatusOK, expectedTask),
			))
		})

		It("returns the task via a GET to /api/v1/tasks/:id endpoint", func() {
			actualTask, err := client.GetTask(1)
			Expect(err).NotTo(HaveOccurred())
			Expect(actualTask).ToNot(BeNil())
			Expect(actualTask).To(Equal(expectedTask))

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when the response has a weird payload", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/tasks/1"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.RespondWith(http.StatusOK, nil),
				))
			})

			It("returns an error", func() {
				_, err := client.GetTask(1)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unexpected response payload: "))
			})
		})

		Context("when the response is a weird status", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/tasks/1"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.RespondWith(http.StatusInternalServerError, nil),
				))
			})

			It("returns an error", func() {
				_, err := client.GetTask(1)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unexpected response status: 500 Internal Server Error"))
			})
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/tasks/1"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.RespondWith(http.StatusNotFound, nil),
				))
			})

			It("returns an error", func() {
				_, err := client.GetTask(1)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unknown task with ID 1"))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("returns an error", func() {
				_, err := client.GetTask(1)
				Expect(err).To(HaveOccurred())
			})
		})
	})

	Describe("UpdatePriority", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", "/api/v1/tasks/1"),
				ghttp.VerifyJSONRepresenting(api.UpdateTaskRequest{Priority: 10}),
				ghttp.RespondWith(http.StatusNoContent, nil),
			))
		})

		It("updates the task via a PUT to /api/v1/tasks/:id endpoint", func() {
			Expect(client.UpdatePriority(1, 10)).To(Succeed())

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/api/v1/tasks/1"),
					ghttp.RespondWith(http.StatusNotFound, nil),
				))
			})

			It("returns an error after hitting the /api/v1/tasks/:id endpoint", func() {
				err := client.UpdatePriority(1, 10)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unknown task with ID 1"))
			})
		})

		Context("when the error payload is weird", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/api/v1/tasks/1"),
					ghttp.RespondWith(http.StatusInternalServerError, nil),
				))
			})

			It("returns an error after hitting the /api/v1/tasks/:id endpoint", func() {
				err := client.UpdatePriority(1, 10)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unexpected response payload: "))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("returns an error", func() {
				Expect(client.UpdatePriority(1, 1234)).NotTo(Succeed())
			})
		})
	})

	Describe("UpdateState", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", "/api/v1/tasks/1"),
				ghttp.VerifyJSONRepresenting(api.UpdateTaskRequest{State: task.StateRunning}),
				ghttp.RespondWith(http.StatusNoContent, nil),
			))
		})

		It("updates the task via a PUT to /api/v1/tasks/:id endpoint", func() {
			Expect(client.UpdateState(1, task.StateRunning)).To(Succeed())

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/api/v1/tasks/1"),
					ghttp.RespondWith(http.StatusNotFound, nil),
				))
			})

			It("returns an error after hitting the /api/v1/tasks/:id endpoint", func() {
				err := client.UpdateState(1, task.StateRunning)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unknown task with ID 1"))
			})
		})

		Context("when the error payload is weird", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/api/v1/tasks/1"),
					ghttp.RespondWith(http.StatusInternalServerError, nil),
				))
			})

			It("returns an error after hitting the /api/v1/tasks/:id endpoint", func() {
				err := client.UpdateState(1, task.StateRunning)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unexpected response payload: "))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("...panics, I guess?", func() {
				Expect(client.UpdateState(1, task.StateReady)).NotTo(Succeed())
			})
		})
	})

	Describe("UpdateName", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", "/api/v1/tasks/1"),
				ghttp.VerifyJSONRepresenting(api.UpdateTaskRequest{Name: "new-name"}),
				ghttp.RespondWith(http.StatusNoContent, nil),
			))
		})

		It("updates the task via a PUT to /api/v1/tasks/:id endpoint", func() {
			Expect(client.UpdateName(1, "new-name")).To(Succeed())

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/api/v1/tasks/1"),
					ghttp.RespondWith(http.StatusNotFound, nil),
				))
			})

			It("returns an error after hitting the /api/v1/tasks/:id endpoint", func() {
				err := client.UpdateName(1, "new-name")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unknown task with ID 1"))
			})
		})

		Context("when the error payload is weird", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("PUT", "/api/v1/tasks/1"),
					ghttp.RespondWith(http.StatusInternalServerError, nil),
				))
			})

			It("returns an error after hitting the /api/v1/tasks/:id endpoint", func() {
				err := client.UpdateName(1, "new-name")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unexpected response payload: "))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("returns an error", func() {
				Expect(client.UpdateName(1, "new-name")).NotTo(Succeed())
			})
		})
	})

	Describe("CreateEvent", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/v1/events"),
				ghttp.VerifyHeaderKV("Content-Type", "application/json"),
				ghttp.VerifyJSONRepresenting(api.AddEventRequest{
					Title:  "event-a",
					Type:   task.EventTypeNote,
					Date:   12345,
					TaskID: 5,
				}),
				ghttp.RespondWithJSONEncoded(http.StatusNoContent, nil),
			))
		})

		It("creates the event with the start time provided", func() {
			Expect(client.CreateEvent("event-a", task.EventTypeNote, 12345, 5)).To(Succeed())

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when the response has a weird status", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/api/v1/events"),
					ghttp.VerifyHeaderKV("Content-Type", "application/json"),
					ghttp.VerifyJSONRepresenting(api.AddEventRequest{
						Title:  "event-a",
						Type:   task.EventTypeNote,
						Date:   12345,
						TaskID: 5,
					}),
					ghttp.RespondWithJSONEncoded(http.StatusInternalServerError, api.ErrorResponse{
						Message: "failed to create event",
					}),
				))
			})

			It("returns the error in the payload", func() {
				err := client.CreateEvent("event-a", task.EventTypeNote, 12345, 5)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("failed to create event"))
			})

		})

		Context("when the response has a weird payload", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("POST", "/api/v1/events"),
					ghttp.VerifyHeaderKV("Content-Type", "application/json"),
					ghttp.VerifyJSONRepresenting(api.AddEventRequest{
						Title:  "event-a",
						Type:   task.EventTypeNote,
						Date:   12345,
						TaskID: 5,
					}),
					ghttp.RespondWith(http.StatusInternalServerError, nil),
				))
			})

			It("returns an error", func() {
				err := client.CreateEvent("event-a", task.EventTypeNote, 12345, 5)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Unexpected response payload: "))
			})
		})
	})

	Describe("DeleteEvent", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/api/v1/events/12345"),
				ghttp.RespondWithJSONEncoded(http.StatusNoContent, nil),
			))
		})

		It("deletes the event with the start time provided", func() {
			Expect(client.DeleteEvent(12345)).To(Succeed())

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when the response has a weird status", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/api/v1/events/12345"),
					ghttp.RespondWithJSONEncoded(http.StatusInternalServerError, api.ErrorResponse{Message: "failed to delete event"}),
				))
			})

			It("returns the error in the payload", func() {
				Expect(client.DeleteEvent(12345)).To(MatchError("failed to delete event"))
			})

		})

		Context("when the response has a weird payload", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/api/v1/events/12345"),
					ghttp.RespondWith(http.StatusInternalServerError, nil),
				))
			})

			It("returns an error", func() {
				err := client.DeleteEvent(12345)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Unexpected response payload: "))
			})

		})
	})

	Describe("GetEvents", func() {
		var events []*task.Event

		BeforeEach(func() {
			events = []*task.Event{
				&task.Event{Title: "event-a", TaskID: 1},
				&task.Event{Title: "event-b", TaskID: 2},
				&task.Event{Title: "event-c", TaskID: 3},
			}
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/events"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.RespondWithJSONEncoded(http.StatusOK, events),
			))
		})

		It("gets all events via a GET to the /api/v1/events endpoint", func() {
			actualEvents, err := client.GetEvents()
			Expect(err).NotTo(HaveOccurred())
			Expect(actualEvents).To(Equal(events))

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when the response payload is bogus", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/events"),
					ghttp.RespondWith(http.StatusOK, nil),
				))
			})

			It("returns an error", func() {
				_, err := client.GetEvents()
				Expect(err).To(MatchError("Unexpected response payload: "))
			})
		})

		Context("when the response status is wrong", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/events"),
					ghttp.RespondWith(http.StatusInternalServerError, nil),
				))
			})

			It("returns an error", func() {
				_, err := client.GetEvents()
				Expect(err).To(MatchError("Unexpected response status: 500 Internal Server Error"))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("returns an error", func() {
				_, err := client.GetEvents()
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
