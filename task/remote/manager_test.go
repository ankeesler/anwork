// TODO: run the generic manager tests against the API!
package remote_test

import (
	"net/http"

	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/task"
	"github.com/ankeesler/anwork/task/remote"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Manager", func() {
	var (
		server  *ghttp.Server
		manager task.Manager
	)

	BeforeEach(func() {
		var err error
		server = ghttp.NewServer()
		server.Writer = GinkgoWriter

		manager, err = remote.NewManagerFactory(server.URL()).Create()
		Expect(err).NotTo(HaveOccurred())
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Create", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/v1/tasks"),
				ghttp.VerifyJSONRepresenting(api.CreateRequest{Name: "a"}),
				ghttp.VerifyHeaderKV("Content-Type", "application/json"),
				ghttp.RespondWith(http.StatusCreated, nil),
			))
		})

		It("successfully can create tasks via a POST to /api/v1/tasks", func() {
			Expect(manager.Create("a")).To(Succeed())
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
				Expect(manager.Create("a")).To(MatchError("failed to create task"))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("...panics, I guess?", func() {
				Expect(func() { manager.Create("a") }).To(Panic())
			})
		})
	})

	Describe("Delete", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/tasks"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.RespondWithJSONEncoded(
					http.StatusOK,
					[]*task.Task{
						&task.Task{Name: "a", ID: 1},
						&task.Task{Name: "b", ID: 2},
						&task.Task{Name: "c", ID: 3},
					}),
			))

			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/api/v1/tasks/1"),
				ghttp.RespondWith(http.StatusNoContent, nil),
			))
		})

		It("successfully can delete tasks via a DELETE to /api/v1/tasks/:id", func() {
			Expect(manager.Delete("a")).To(Succeed())
			Expect(server.ReceivedRequests()).To(HaveLen(2))
		})

		Context("when the request returns a failure", func() {
			BeforeEach(func() {
				server.SetHandler(1, ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/api/v1/tasks/1"),
					ghttp.RespondWithJSONEncoded(
						http.StatusBadRequest,
						api.ErrorResponse{Message: "failed to delete task"},
					),
				))
			})

			It("prints the failure message", func() {
				Expect(manager.Delete("a")).To(MatchError("failed to delete task"))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("returns an error", func() {
				err := manager.Delete("a")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("connection refused"))
			})
		})
	})

	Describe("Tasks", func() {
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
			actualTasks := manager.Tasks()
			Expect(actualTasks).ToNot(BeNil())
			Expect(actualTasks).To(Equal(tasks))

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("...panics, I guess?", func() {
				Expect(func() { manager.Tasks() }).To(Panic())
			})
		})
	})

	Describe("FindByName", func() {
		var expectedTask *task.Task

		BeforeEach(func() {
			tasks := []*task.Task{
				&task.Task{Name: "task-a", ID: 1},
				&task.Task{Name: "task-b", ID: 2},
				&task.Task{Name: "task-c", ID: 3},
			}
			expectedTask = tasks[0]
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/tasks"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.RespondWithJSONEncoded(http.StatusOK, tasks),
			))
		})

		It("returns the task via a GET to /api/v1/tasks endpoint", func() {
			actualTask := manager.FindByName("task-a")
			Expect(actualTask).ToNot(BeNil())
			Expect(actualTask).To(Equal(expectedTask))

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				tasks := []*task.Task{
					&task.Task{Name: "task-b", ID: 2},
					&task.Task{Name: "task-c", ID: 3},
				}
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/tasks"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.RespondWithJSONEncoded(http.StatusOK, tasks),
				))
			})

			It("returns nil after hitting the /api/v1/tasks endpoint", func() {
				actualTask := manager.FindByName("task-a")
				Expect(actualTask).To(BeNil())
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("...panics, I guess?", func() {
				Expect(func() { manager.FindByName("tuna") }).To(Panic())
			})
		})
	})

	Describe("FindByID", func() {
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
			actualTask := manager.FindByID(1)
			Expect(actualTask).ToNot(BeNil())
			Expect(actualTask).To(Equal(expectedTask))

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/tasks/1"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.RespondWith(http.StatusNotFound, nil),
				))
			})

			It("returns nil after hitting the /api/v1/tasks/1 endpoint", func() {
				actualTask := manager.FindByID(1)
				Expect(actualTask).To(BeNil())
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("...panics, I guess?", func() {
				Expect(func() { manager.FindByName("tuna") }).To(Panic())
			})
		})
	})

	Describe("SetPriority", func() {
		BeforeEach(func() {
			tasks := []*task.Task{
				&task.Task{Name: "task-a", ID: 1},
				&task.Task{Name: "task-b", ID: 2},
				&task.Task{Name: "task-c", ID: 3},
			}
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/tasks"),
				ghttp.RespondWithJSONEncoded(http.StatusOK, tasks),
			))
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", "/api/v1/tasks/1"),
				ghttp.VerifyJSONRepresenting(api.UpdateTaskRequest{Priority: 10}),
				ghttp.RespondWith(http.StatusNoContent, nil),
			))
		})

		It("updates the task via a PUT to /api/v1/tasks/:id endpoint", func() {
			Expect(manager.SetPriority("task-a", 10)).To(Succeed())

			Expect(server.ReceivedRequests()).To(HaveLen(2))
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				tasks := []*task.Task{
					&task.Task{Name: "task-b", ID: 2},
					&task.Task{Name: "task-c", ID: 3},
				}
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/tasks"),
					ghttp.RespondWithJSONEncoded(http.StatusOK, tasks),
				))
			})

			It("returns an error after hitting the /api/v1/tasks/:id endpoint", func() {
				err := manager.SetPriority("task-a", 10)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unknown task"))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("...panics, I guess?", func() {
				Expect(func() { manager.SetPriority("task-a", 1234) }).To(Panic())
			})
		})
	})

	Describe("SetState", func() {
		BeforeEach(func() {
			tasks := []*task.Task{
				&task.Task{Name: "task-a", ID: 1},
				&task.Task{Name: "task-b", ID: 2},
				&task.Task{Name: "task-c", ID: 3},
			}
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/tasks"),
				ghttp.RespondWithJSONEncoded(http.StatusOK, tasks),
			))
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", "/api/v1/tasks/1"),
				ghttp.VerifyJSONRepresenting(api.UpdateTaskRequest{State: task.StateRunning}),
				ghttp.RespondWith(http.StatusNoContent, nil),
			))
		})

		It("updates the task via a PUT to /api/v1/tasks/:id endpoint", func() {
			Expect(manager.SetState("task-a", task.StateRunning)).To(Succeed())

			Expect(server.ReceivedRequests()).To(HaveLen(2))
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				tasks := []*task.Task{
					&task.Task{Name: "task-b", ID: 2},
					&task.Task{Name: "task-c", ID: 3},
				}
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/tasks"),
					ghttp.RespondWithJSONEncoded(http.StatusOK, tasks),
				))
			})

			It("returns an error after hitting the /api/v1/tasks/:id endpoint", func() {
				err := manager.SetState("task-a", task.StateRunning)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unknown task"))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("...panics, I guess?", func() {
				Expect(func() { manager.SetPriority("task-a", 1234) }).To(Panic())
			})
		})
	})

	Describe("Note", func() {
		BeforeEach(func() {
			tasks := []*task.Task{
				&task.Task{Name: "task-a", ID: 1},
				&task.Task{Name: "task-b", ID: 2},
				&task.Task{Name: "task-c", ID: 3},
			}
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/tasks"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.RespondWithJSONEncoded(http.StatusOK, tasks),
			))
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/v1/events"),
				ghttp.VerifyHeaderKV("Content-Type", "application/json"),
				ghttp.RespondWith(http.StatusNoContent, nil),
			))
		})

		It("adds a note via a POST to /api/v1/events", func() {
			Expect(manager.Note("task-a", "here is a note")).To(Succeed())

			Expect(server.ReceivedRequests()).To(HaveLen(2))

			// TODO: how do we test that the body has the right stuff???
			// Is this a sign that we should be using an interface for time.Now()...
			//var payload api.AddEventRequest
			//body := server.ReceivedRequests()[1].Body
			//decoder := json.NewDecoder(body)
			//Expect(decoder.Decode(&payload)).To(Succeed())
			//Expect(payload.Title).To(Equal("Note added to task task-a: here is a note"))
			//Expect(payload.Date).To(BeNumerically("<=", time.Now().Unix()))
			//Expect(payload.Type).To(Equal(task.EventTypeNote))
			//Expect(payload.TaskID).To(Equal(1))
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				tasks := []*task.Task{
					&task.Task{Name: "task-b", ID: 2},
					&task.Task{Name: "task-c", ID: 3},
				}
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest("GET", "/api/v1/tasks"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.RespondWithJSONEncoded(http.StatusOK, tasks),
				))
			})

			It("returns a helpful error", func() {
				err := manager.Note("task-a", "here is a note")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unknown task task-a"))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("...panics, I guess?", func() {
				Expect(func() { manager.Events() }).To(Panic())
			})
		})
	})

	Describe("Events", func() {
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

		It("updates the task via a GET to the /api/v1/events endpoint", func() {
			actualEvents := manager.Events()
			Expect(actualEvents).To(Equal(events))

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("...panics, I guess?", func() {
				Expect(func() { manager.Events() }).To(Panic())
			})
		})
	})

	Describe("Reset", func() {
		var tasks []*task.Task
		var events []*task.Event
		BeforeEach(func() {
			tasks = []*task.Task{
				&task.Task{Name: "task-a", ID: 1},
				&task.Task{Name: "task-b", ID: 2},
				&task.Task{Name: "task-c", ID: 3},
			}
			events = []*task.Event{
				&task.Event{Title: "event-a", TaskID: 1},
				&task.Event{Title: "event-b", TaskID: 2},
				&task.Event{Title: "event-c", TaskID: 3},
			}
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/tasks"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.RespondWithJSONEncoded(http.StatusOK, tasks),
			))
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/events"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.RespondWithJSONEncoded(http.StatusOK, events),
			))
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/api/v1/tasks/1"),
				ghttp.RespondWith(http.StatusNoContent, nil),
			))
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/api/v1/tasks/2"),
				ghttp.RespondWith(http.StatusNoContent, nil),
			))
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/api/v1/tasks/3"),
				ghttp.RespondWith(http.StatusNoContent, nil),
			))
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/api/v1/events/1"),
				ghttp.RespondWith(http.StatusNoContent, nil),
			))
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/api/v1/events/2"),
				ghttp.RespondWith(http.StatusNoContent, nil),
			))
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/api/v1/events/3"),
				ghttp.RespondWith(http.StatusNoContent, nil),
			))
		})

		It("DELETE's all of the tasks and events", func() {
			Expect(manager.Reset()).To(Succeed())
			Expect(server.ReceivedRequests()).To(HaveLen(8))
		})

		Context("when the request returns a failure", func() {
			BeforeEach(func() {
				server.SetHandler(2, ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/api/v1/tasks/1"),
					ghttp.RespondWithJSONEncoded(
						http.StatusBadRequest,
						api.ErrorResponse{Message: "failed to delete task"},
					),
				))
			})

			It("prints the failure message", func() {
				err := manager.Reset()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Encountered errors during reset:\n"))
				Expect(err.Error()).To(ContainSubstring("  delete task 1: failed to delete task"))
			})
		})

		Context("when the request returns an unexpected response with no payload", func() {
			BeforeEach(func() {
				server.SetHandler(5, ghttp.CombineHandlers(
					ghttp.VerifyRequest("DELETE", "/api/v1/events/1"),
					ghttp.RespondWith(http.StatusMethodNotAllowed, nil),
				))
			})

			It("prints a legitimate failure message", func() {
				err := manager.Reset()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Encountered errors during reset:\n"))
				Expect(err.Error()).To(ContainSubstring("  delete event 1: unexpected response: 405 Method Not Allowed"))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("returns an error", func() {
				err := manager.Reset()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("connection refused"))
			})
		})
	})
})
