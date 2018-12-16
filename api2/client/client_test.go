package client_test

import (
	"net/http"

	api "github.com/ankeesler/anwork/api2"
	apiclient "github.com/ankeesler/anwork/api2/client"
	"github.com/ankeesler/anwork/task2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Client", func() {
	var (
		client task2.Repo
		server *ghttp.Server

		tasks  []*task2.Task
		events []*task2.Event
	)

	testBadURL := func(clientFunc func(c task2.Repo) error) {
		It("returns error on bad URL", func() {
			c := apiclient.New("here is a bad url i mean i am sure this is bad/://aff;a;f/a;'a';sd")
			Expect(clientFunc(c)).NotTo(Succeed())

			Expect(server.ReceivedRequests()).To(HaveLen(0))
		})
	}

	testFailedRequest := func(clientFunc func(c task2.Repo) error) {
		It("returns error on failed request", func() {
			c := apiclient.New("asdf")
			Expect(clientFunc(c)).NotTo(Succeed())

			Expect(server.ReceivedRequests()).To(HaveLen(0))
		})
	}

	test4xxResponse := func(clientFunc func(c task2.Repo) error) {
		Context("on 4xx response", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusNotFound, nil),
				))
			})

			It("returns error on 4xx response containing status", func() {
				err := clientFunc(client)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("404 Not Found"))

				Expect(server.ReceivedRequests()).To(HaveLen(1))
			})
		})
	}

	test5xxResponse := func(clientFunc func(c task2.Repo) error) {
		Context("on 5xx response", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.RespondWithJSONEncoded(
						http.StatusInternalServerError,
						api.Error{Message: "some message"},
						http.Header{"Content-Type": {"application/json"}},
					),
				))
			})

			It("returns error on 5xx response containing status and message", func() {
				err := clientFunc(client)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("500 Internal Server Error"))
				Expect(err.Error()).To(ContainSubstring("some message"))

				Expect(server.ReceivedRequests()).To(HaveLen(1))
			})

			Context("when payload is not api.Error", func() {
				BeforeEach(func() {
					server.SetHandler(0, ghttp.CombineHandlers(
						ghttp.RespondWithJSONEncoded(
							http.StatusInternalServerError,
							"here is something",
							http.Header{"Content-Type": {"application/json"}},
						),
					))
				})

				It("returns error on 5xx response containing status and ???", func() {
					err := clientFunc(client)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("500 Internal Server Error"))
					Expect(err.Error()).To(ContainSubstring("???"))

					Expect(server.ReceivedRequests()).To(HaveLen(1))
				})
			})
		})
	}

	testAllCommonFailures := func(clientFunc func(c task2.Repo) error) {
		testBadURL(clientFunc)
		testFailedRequest(clientFunc)
		test4xxResponse(clientFunc)
		test5xxResponse(clientFunc)
	}

	testBad2xxResponseBody := func(clientFunc func(c task2.Repo) error) {
		Context("when the response payload is invalid", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.RespondWithJSONEncoded(
						http.StatusOK,
						"bad payload",
						http.Header{"Content-Type": {"application/json"}},
					),
				))
			})

			It("returns an error", func() {
				_, err := client.Tasks()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("bad payload"))

				Expect(server.ReceivedRequests()).To(HaveLen(1))
			})
		})
	}

	BeforeEach(func() {
		server = ghttp.NewServer()
		client = apiclient.New(server.Addr())

		tasks = []*task2.Task{
			&task2.Task{Name: "task-a", ID: 1},
			&task2.Task{Name: "task-b", ID: 2},
			&task2.Task{Name: "task-c", ID: 3},
		}
		events = []*task2.Event{
			&task2.Event{Title: "event-a", ID: 1},
			&task2.Event{Title: "event-b", ID: 2},
			&task2.Event{Title: "event-c", ID: 3},
		}
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("CreateTask", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodPost, "/api/v1/tasks"),
				ghttp.VerifyJSONRepresenting(tasks[0]),
				ghttp.VerifyHeaderKV("Content-Type", "application/json"),
				ghttp.RespondWith(
					http.StatusCreated,
					nil,
					http.Header{"Location": {"/api/v1/tasks/1"}}),
			))
		})

		It("POSTs to /api/v1/tasks", func() {
			task := tasks[0]
			Expect(client.CreateTask(task)).To(Succeed())

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		testAllCommonFailures(func(c task2.Repo) error {
			return c.CreateTask(tasks[0])
		})
	})

	Describe("Tasks", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodGet, "/api/v1/tasks"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.RespondWithJSONEncoded(
					http.StatusOK,
					tasks,
					http.Header{"Content-Type": {"application/json"}},
				),
			))
		})

		It("gets the tasks from the server", func() {
			rspTasks, err := client.Tasks()
			Expect(err).NotTo(HaveOccurred())
			Expect(rspTasks).To(Equal(tasks))

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		testBad2xxResponseBody(func(c task2.Repo) error {
			_, err := c.Tasks()
			return err
		})

		testAllCommonFailures(func(c task2.Repo) error {
			_, err := c.Tasks()
			return err
		})
	})

	Describe("FindTaskByID", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodGet, "/api/v1/tasks/10"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.RespondWithJSONEncoded(
					http.StatusOK,
					tasks[0],
					http.Header{"Content-Type": {"application/json"}},
				),
			))
		})

		It("finds a task by ID", func() {
			task, err := client.FindTaskByID(10)
			Expect(err).NotTo(HaveOccurred())
			Expect(task).To(Equal(tasks[0]))

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		testBad2xxResponseBody(func(c task2.Repo) error {
			_, err := c.FindTaskByID(10)
			return err
		})

		testAllCommonFailures(func(c task2.Repo) error {
			_, err := c.FindTaskByID(10)
			return err
		})
	})

	Describe("FindTaskByName", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodGet, "/api/v1/tasks", "name=task-a"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.RespondWithJSONEncoded(
					http.StatusOK,
					tasks[0],
					http.Header{"Content-Type": {"application/json"}},
				),
			))
		})

		It("finds a task by name", func() {
			task, err := client.FindTaskByName("task-a")
			Expect(err).NotTo(HaveOccurred())
			Expect(task).To(Equal(tasks[0]))

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		testBad2xxResponseBody(func(c task2.Repo) error {
			_, err := c.FindTaskByName("task-a")
			return err
		})

		testAllCommonFailures(func(c task2.Repo) error {
			_, err := c.FindTaskByName("task-a")
			return err
		})
	})

	Describe("UpdateTask", func() {
		var task task2.Task
		BeforeEach(func() {
			task = *tasks[0]
			task.Name = "updated-task-a"

			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodPut, "/api/v1/tasks/1"),
				ghttp.VerifyHeaderKV("Content-Type", "application/json"),
				ghttp.VerifyJSONRepresenting(task),
				ghttp.RespondWith(http.StatusNoContent, nil),
			))
		})

		It("updates a task", func() {
			err := client.UpdateTask(&task)
			Expect(err).NotTo(HaveOccurred())

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		testAllCommonFailures(func(c task2.Repo) error {
			return c.UpdateTask(&task)
		})
	})

	Describe("DeleteTask", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodDelete, "/api/v1/tasks/10"),
				ghttp.RespondWith(http.StatusNoContent, nil),
			))
		})

		It("deletes a task by ID", func() {
			tasks[0].ID = 10
			Expect(client.DeleteTask(tasks[0])).To(Succeed())

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		testAllCommonFailures(func(c task2.Repo) error {
			tasks[0].ID = 10
			return c.DeleteTask(tasks[0])
		})
	})

	Describe("CreateEvent", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodPost, "/api/v1/events"),
				ghttp.VerifyJSONRepresenting(events[0]),
				ghttp.VerifyHeaderKV("Content-Type", "application/json"),
				ghttp.RespondWith(
					http.StatusCreated,
					nil,
					http.Header{"Location": {"/api/v1/events/1"}}),
			))
		})

		It("POSTs to /api/v1/events", func() {
			event := events[0]
			Expect(client.CreateEvent(event)).To(Succeed())

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		testAllCommonFailures(func(c task2.Repo) error {
			return c.CreateEvent(events[0])
		})
	})

	Describe("Events", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodGet, "/api/v1/events"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.RespondWithJSONEncoded(
					http.StatusOK,
					events,
					http.Header{"Content-Type": {"application/json"}},
				),
			))
		})

		It("gets the events from the server", func() {
			rspEvents, err := client.Events()
			Expect(err).NotTo(HaveOccurred())
			Expect(rspEvents).To(Equal(events))

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		testBad2xxResponseBody(func(c task2.Repo) error {
			_, err := c.Events()
			return err
		})

		testAllCommonFailures(func(c task2.Repo) error {
			_, err := c.Events()
			return err
		})
	})

	Describe("FindEventByID", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodGet, "/api/v1/events/10"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.RespondWithJSONEncoded(
					http.StatusOK,
					events[0],
					http.Header{"Content-Type": {"application/json"}}),
			))
		})

		It("gets the event by ID", func() {
			event, err := client.FindEventByID(10)
			Expect(err).NotTo(HaveOccurred())
			Expect(event).To(Equal(events[0]))

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		testAllCommonFailures(func(c task2.Repo) error {
			_, err := c.FindEventByID(10)
			return err
		})
	})

	Describe("DeleteEvent", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodDelete, "/api/v1/events/10"),
				ghttp.RespondWith(http.StatusNoContent, nil),
			))
		})

		It("deletes an event by ID", func() {
			events[0].ID = 10
			Expect(client.DeleteEvent(events[0])).To(Succeed())

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		testAllCommonFailures(func(c task2.Repo) error {
			events[0].ID = 10
			return c.DeleteEvent(events[0])
		})
	})
})
