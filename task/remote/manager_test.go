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

		statusCode int
		body       interface{}
	)

	BeforeEach(func() {
		var err error
		server = ghttp.NewServer()
		manager, err = remote.NewManagerFactory(server.URL()).Create()
		Expect(err).NotTo(HaveOccurred())

		server.AppendHandlers(ghttp.CombineHandlers(
			ghttp.RespondWithJSONEncodedPtr(&statusCode, &body),
		))
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("Create", func() {
		BeforeEach(func() {
			statusCode = http.StatusCreated
			server.WrapHandler(0, ghttp.CombineHandlers(
				ghttp.VerifyRequest("POST", "/api/v1/tasks"),
				ghttp.VerifyJSONRepresenting(api.CreateRequest{Name: "a"}),
			))
		})

		It("successfully can create tasks via a POST to /api/v1/tasks", func() {
			Expect(manager.Create("a")).To(Succeed())
		})

		Context("when the request returns a failure", func() {
			BeforeEach(func() {
				statusCode = http.StatusBadRequest
				body = api.ErrorResponse{Message: "failed to create task"}
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
			statusCode = http.StatusNoContent
			server.WrapHandler(0, ghttp.CombineHandlers(
				ghttp.VerifyRequest("DELETE", "/api/v1/tasks/a"),
			))
		})

		It("successfully can delete tasks via a DELETE to /api/v1/tasks/:name", func() {
			Expect(manager.Delete("a")).To(Succeed())
		})

		Context("when the request returns a failure", func() {
			BeforeEach(func() {
				statusCode = http.StatusBadRequest
				body = api.ErrorResponse{Message: "failed to delete task"}
			})

			It("prints the failure message", func() {
				Expect(manager.Delete("a")).To(MatchError("failed to delete task"))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("...panics, I guess?", func() {
				Expect(func() { manager.Delete("a") }).To(Panic())
			})
		})
	})

	Describe("Tasks", func() {
		var tasks []*task.Task

		BeforeEach(func() {
			statusCode = http.StatusOK
			tasks = []*task.Task{
				&task.Task{Name: "task-a", ID: 1},
				&task.Task{Name: "task-b", ID: 2},
				&task.Task{Name: "task-c", ID: 3},
			}
			body = tasks

			server.WrapHandler(0, ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/tasks"),
			))
		})

		It("returns the tasks via a GET to /api/v1/tasks endpoint", func() {
			actualTasks := manager.Tasks()
			Expect(actualTasks).ToNot(BeNil())
			Expect(actualTasks).To(Equal(tasks))
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
			statusCode = http.StatusOK
			expectedTask = &task.Task{Name: "task-a", ID: 1}
			body = expectedTask

			server.WrapHandler(0, ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/tasks/task-a"),
			))
		})

		It("returns the task via a GET to /api/v1/tasks/:name endpoint", func() {
			actualTask := manager.FindByName("task-a")
			Expect(actualTask).ToNot(BeNil())
			Expect(actualTask).To(Equal(expectedTask))
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				statusCode = http.StatusNotFound
			})

			It("returns nil after hitting the /api/v1/tasks/:name endpoint", func() {
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
		var tasks []*task.Task
		var expectedTask *task.Task

		BeforeEach(func() {
			statusCode = http.StatusOK
			tasks = []*task.Task{
				&task.Task{Name: "task-a", ID: 1},
				&task.Task{Name: "task-b", ID: 2},
				&task.Task{Name: "task-c", ID: 3},
			}
			expectedTask = tasks[0]
			body = tasks

			server.WrapHandler(0, ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/tasks"),
			))
		})

		It("returns the task via a GET to /api/v1/tasks endpoint", func() {
			actualTask := manager.FindByID(1)
			Expect(actualTask).ToNot(BeNil())
			Expect(actualTask).To(Equal(expectedTask))
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				tasks = tasks[1:]
				body = tasks
			})

			It("returns nil after hitting the /api/v1/tasks endpoint", func() {
				actualTask := manager.FindByID(1)
				Expect(actualTask).To(BeNil())
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				server.Close()
			})

			It("...panics, I guess?", func() {
				Expect(func() { manager.FindByID(1234) }).To(Panic())
			})
		})
	})

	Describe("SetPriority", func() {
		BeforeEach(func() {
			statusCode = http.StatusNoContent
			server.WrapHandler(0, ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", "/api/v1/tasks/task-a"),
				ghttp.VerifyJSONRepresenting(api.SetRequest{Priority: 10}),
			))
		})

		It("updates the task via a PUT to /api/v1/tasks/:name endpoint", func() {
			Expect(manager.SetPriority("task-a", 10)).To(Succeed())
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				statusCode = http.StatusNotFound
				body = api.ErrorResponse{Message: "unknown task"}
			})

			It("returns an error after hitting the /api/v1/tasks/:name endpoint", func() {
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
			statusCode = http.StatusNoContent
			server.WrapHandler(0, ghttp.CombineHandlers(
				ghttp.VerifyRequest("PUT", "/api/v1/tasks/task-a"),
				ghttp.VerifyJSONRepresenting(api.SetRequest{State: task.StateRunning}),
			))
		})

		It("updates the task via a PUT to /api/v1/tasks/:name endpoint", func() {
			Expect(manager.SetState("task-a", task.StateRunning)).To(Succeed())
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				statusCode = http.StatusNotFound
				body = api.ErrorResponse{Message: "unknown task"}
			})

			It("returns an error after hitting the /api/v1/tasks/:name endpoint", func() {
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
				Expect(func() { manager.SetState("task-a", task.StateRunning) }).To(Panic())
			})
		})
	})

	Describe("Note", func() {
	})

	Describe("Events", func() {
		var events []*task.Event

		BeforeEach(func() {
			statusCode = http.StatusOK
			events = []*task.Event{
				&task.Event{Title: "task-a", TaskID: 1},
				&task.Event{Title: "task-b", TaskID: 2},
				&task.Event{Title: "task-c", TaskID: 3},
			}
			body = events

			server.WrapHandler(0, ghttp.CombineHandlers(
				ghttp.VerifyRequest("GET", "/api/v1/events"),
			))
		})

		It("updates the task via a GET to the /api/v1/events endpoint", func() {
			actualEvents := manager.Events()
			Expect(actualEvents).To(Equal(events))
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
})
