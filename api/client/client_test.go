package client_test

import (
	"net/http"

	"github.com/ankeesler/anwork/api"
	clientpkg "github.com/ankeesler/anwork/api/client"
	"github.com/ankeesler/anwork/api/client/clientfakes"
	taskpkg "github.com/ankeesler/anwork/task"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Client", func() {
	var (
		authenticator *clientfakes.FakeAuthenticator
		cache         *clientfakes.FakeCache

		client taskpkg.Repo
		server *ghttp.Server

		tasks  []*taskpkg.Task
		events []*taskpkg.Event
	)

	testBadURL := func(clientFunc func(c taskpkg.Repo) error) {
		It("returns error on bad URL", func() {
			c := clientpkg.New(
				makeLogger(),
				"here is a bad url i mean i am sure this is bad/://aff;a;f/a;'a';sd",
				authenticator,
				cache,
			)
			Expect(clientFunc(c)).NotTo(Succeed())

			Expect(server.ReceivedRequests()).To(HaveLen(0))

			Expect(cache.GetCallCount()).To(Equal(0))
		})
	}

	testFailedRequest := func(clientFunc func(c taskpkg.Repo) error) {
		It("returns error on failed request", func() {
			c := clientpkg.New(
				makeLogger(),
				"asdf",
				authenticator,
				cache,
			)
			Expect(clientFunc(c)).NotTo(Succeed())

			Expect(server.ReceivedRequests()).To(HaveLen(0))

			Expect(cache.GetCallCount()).To(Equal(1))
		})
	}

	test4xxResponse := func(clientFunc func(c taskpkg.Repo) error) {
		Context("on 4xx response", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusMethodNotAllowed, nil),
				))
			})

			It("returns error on 4xx response containing status", func() {
				err := clientFunc(client)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("405 Method Not Allowed"))

				Expect(server.ReceivedRequests()).To(HaveLen(1))

				Expect(cache.GetCallCount()).To(Equal(1))
			})
		})
	}

	test5xxResponse := func(clientFunc func(c taskpkg.Repo) error) {
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

				Expect(cache.GetCallCount()).To(Equal(1))
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

					Expect(cache.GetCallCount()).To(Equal(1))
				})
			})
		})
	}

	testAllCommonFailures := func(clientFunc func(c taskpkg.Repo) error) {
		testBadURL(clientFunc)
		testFailedRequest(clientFunc)
		test4xxResponse(clientFunc)
		test5xxResponse(clientFunc)
	}

	testBad2xxResponseBody := func(clientFunc func(c taskpkg.Repo) error) {
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

				Expect(cache.GetCallCount()).To(Equal(1))
			})
		})
	}

	BeforeEach(func() {
		authenticator = &clientfakes.FakeAuthenticator{}

		cache = &clientfakes.FakeCache{}
		cache.GetReturnsOnCall(0, "bearer some-token", true)

		server = ghttp.NewServer()
		client = clientpkg.New(
			makeLogger(),
			server.Addr(),
			authenticator,
			cache,
		)

		tasks = []*taskpkg.Task{
			&taskpkg.Task{Name: "task-a", ID: 1},
			&taskpkg.Task{Name: "task-b", ID: 2},
			&taskpkg.Task{Name: "task-c", ID: 3},
		}
		events = []*taskpkg.Event{
			&taskpkg.Event{Title: "event-a", ID: 1},
			&taskpkg.Event{Title: "event-b", ID: 2},
			&taskpkg.Event{Title: "event-c", ID: 3},
		}
	})

	AfterEach(func() {
		server.Close()
	})

	Describe("CreateTask", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodPost, "/api/v1/tasks"),
				ghttp.VerifyHeaderKV("Authorization", "bearer some-token"),
				ghttp.VerifyJSONRepresenting(tasks[0]),
				ghttp.VerifyHeaderKV("Content-Type", "application/json"),
				ghttp.RespondWith(
					http.StatusCreated,
					nil,
					http.Header{"Location": {"/api/v1/tasks/10"}}),
			))
		})

		It("POSTs to /api/v1/tasks", func() {
			task := tasks[0]
			Expect(client.CreateTask(task)).To(Succeed())

			Expect(server.ReceivedRequests()).To(HaveLen(1))

			Expect(cache.GetCallCount()).To(Equal(1))
		})

		It("sets the provided task's ID to the newly allocated ID", func() {
			task := tasks[0]
			Expect(client.CreateTask(task)).To(Succeed())

			Expect(task.ID).To(Equal(10))

			Expect(server.ReceivedRequests()).To(HaveLen(1))

			Expect(cache.GetCallCount()).To(Equal(1))
		})

		Context("when the returned location is invalid", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodPost, "/api/v1/tasks"),
					ghttp.VerifyJSONRepresenting(tasks[0]),
					ghttp.VerifyHeaderKV("Content-Type", "application/json"),
					ghttp.VerifyHeaderKV("Authorization", "bearer some-token"),
					ghttp.RespondWith(
						http.StatusCreated,
						nil,
						http.Header{"Location": {"/api/v1/tasks/tuna"}}),
				))
			})

			It("returns an error", func() {
				task := tasks[0]
				err := client.CreateTask(task)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("could not parse ID from Location response header: /api/v1/tasks/tuna"))

				Expect(server.ReceivedRequests()).To(HaveLen(1))

				Expect(cache.GetCallCount()).To(Equal(1))
			})
		})

		testAllCommonFailures(func(c taskpkg.Repo) error {
			return c.CreateTask(tasks[0])
		})
	})

	Describe("Tasks", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodGet, "/api/v1/tasks"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer some-token"),
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

			Expect(cache.GetCallCount()).To(Equal(1))
		})

		testBad2xxResponseBody(func(c taskpkg.Repo) error {
			_, err := c.Tasks()
			return err
		})

		testAllCommonFailures(func(c taskpkg.Repo) error {
			_, err := c.Tasks()
			return err
		})
	})

	Describe("FindTaskByID", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodGet, "/api/v1/tasks/10"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer some-token"),
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

			Expect(cache.GetCallCount()).To(Equal(1))
		})

		Context("on 404 not found response", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusNotFound, nil),
				))
			})

			It("returns nil, nil", func() {
				task, err := client.FindTaskByID(10)
				Expect(err).NotTo(HaveOccurred())
				Expect(task).To(BeNil())

				Expect(server.ReceivedRequests()).To(HaveLen(1))

				Expect(cache.GetCallCount()).To(Equal(1))
			})
		})

		testBad2xxResponseBody(func(c taskpkg.Repo) error {
			_, err := c.FindTaskByID(10)
			return err
		})

		testAllCommonFailures(func(c taskpkg.Repo) error {
			_, err := c.FindTaskByID(10)
			return err
		})
	})

	Describe("FindTaskByName", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodGet, "/api/v1/tasks", "name=task-a"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer some-token"),
				ghttp.RespondWithJSONEncoded(
					http.StatusOK,
					[]*taskpkg.Task{tasks[0]},
					http.Header{"Content-Type": {"application/json"}},
				),
			))
		})

		It("finds a task by name", func() {
			task, err := client.FindTaskByName("task-a")
			Expect(err).NotTo(HaveOccurred())
			Expect(task).To(Equal(tasks[0]))

			Expect(server.ReceivedRequests()).To(HaveLen(1))

			Expect(cache.GetCallCount()).To(Equal(1))
		})

		Context("when the server responds with an empty array of tasks", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.RespondWithJSONEncoded(
						http.StatusOK,
						[]*taskpkg.Task{},
						http.Header{"Content-Type": {"application/json"}}),
				))
			})

			It("returns nil, nil", func() {
				task, err := client.FindTaskByName("task-a")
				Expect(err).NotTo(HaveOccurred())
				Expect(task).To(BeNil())

				Expect(server.ReceivedRequests()).To(HaveLen(1))

				Expect(cache.GetCallCount()).To(Equal(1))
			})
		})

		testBad2xxResponseBody(func(c taskpkg.Repo) error {
			_, err := c.FindTaskByName("task-a")
			return err
		})

		testAllCommonFailures(func(c taskpkg.Repo) error {
			_, err := c.FindTaskByName("task-a")
			return err
		})
	})

	Describe("UpdateTask", func() {
		var task taskpkg.Task
		BeforeEach(func() {
			task = *tasks[0]
			task.Name = "updated-task-a"

			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodPut, "/api/v1/tasks/1"),
				ghttp.VerifyHeaderKV("Content-Type", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer some-token"),
				ghttp.VerifyJSONRepresenting(task),
				ghttp.RespondWith(http.StatusNoContent, nil),
			))
		})

		It("updates a task", func() {
			err := client.UpdateTask(&task)
			Expect(err).NotTo(HaveOccurred())

			Expect(server.ReceivedRequests()).To(HaveLen(1))

			Expect(cache.GetCallCount()).To(Equal(1))
		})

		testAllCommonFailures(func(c taskpkg.Repo) error {
			return c.UpdateTask(&task)
		})
	})

	Describe("DeleteTask", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodDelete, "/api/v1/tasks/10"),
				ghttp.VerifyHeaderKV("Authorization", "bearer some-token"),
				ghttp.RespondWith(http.StatusNoContent, nil),
			))
		})

		It("deletes a task by ID", func() {
			tasks[0].ID = 10
			Expect(client.DeleteTask(tasks[0])).To(Succeed())

			Expect(server.ReceivedRequests()).To(HaveLen(1))

			Expect(cache.GetCallCount()).To(Equal(1))
		})

		testAllCommonFailures(func(c taskpkg.Repo) error {
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
				ghttp.VerifyHeaderKV("Authorization", "bearer some-token"),
				ghttp.RespondWith(
					http.StatusCreated,
					nil,
					http.Header{"Location": {"/api/v1/events/10"}}),
			))
		})

		It("POSTs to /api/v1/events", func() {
			event := events[0]
			Expect(client.CreateEvent(event)).To(Succeed())

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		It("sets the provided event's ID to the newly allocated ID", func() {
			event := events[0]
			Expect(client.CreateEvent(event)).To(Succeed())

			Expect(event.ID).To(Equal(10))

			Expect(server.ReceivedRequests()).To(HaveLen(1))
		})

		Context("when the returned location is invalid", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodPost, "/api/v1/events"),
					ghttp.VerifyJSONRepresenting(events[0]),
					ghttp.VerifyHeaderKV("Content-Type", "application/json"),
					ghttp.VerifyHeaderKV("Authorization", "bearer some-token"),
					ghttp.RespondWith(
						http.StatusCreated,
						nil,
						http.Header{"Location": {"/api/v1/events/tuna"}}),
				))
			})

			It("returns an error", func() {
				event := events[0]
				err := client.CreateEvent(event)
				Expect(err).To(HaveOccurred())
				Expect(err).To(MatchError("could not parse ID from Location response header: /api/v1/events/tuna"))

				Expect(server.ReceivedRequests()).To(HaveLen(1))
			})
		})

		testAllCommonFailures(func(c taskpkg.Repo) error {
			return c.CreateEvent(events[0])
		})
	})

	Describe("Events", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodGet, "/api/v1/events"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer some-token"),
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

			Expect(cache.GetCallCount()).To(Equal(1))
		})

		testBad2xxResponseBody(func(c taskpkg.Repo) error {
			_, err := c.Events()
			return err
		})

		testAllCommonFailures(func(c taskpkg.Repo) error {
			_, err := c.Events()
			return err
		})
	})

	Describe("FindEventByID", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodGet, "/api/v1/events/10"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer some-token"),
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

			Expect(cache.GetCallCount()).To(Equal(1))
		})

		Context("on 404 not found response", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.RespondWith(http.StatusNotFound, nil),
				))
			})

			It("returns nil, nil", func() {
				event, err := client.FindEventByID(10)
				Expect(err).NotTo(HaveOccurred())
				Expect(event).To(BeNil())

				Expect(server.ReceivedRequests()).To(HaveLen(1))

				Expect(cache.GetCallCount()).To(Equal(1))
			})
		})

		testBad2xxResponseBody(func(c taskpkg.Repo) error {
			_, err := c.FindEventByID(10)
			return err
		})

		testAllCommonFailures(func(c taskpkg.Repo) error {
			_, err := c.FindEventByID(10)
			return err
		})
	})

	Describe("DeleteEvent", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodDelete, "/api/v1/events/10"),
				ghttp.VerifyHeaderKV("Authorization", "bearer some-token"),
				ghttp.RespondWith(http.StatusNoContent, nil),
			))
		})

		It("deletes an event by ID", func() {
			events[0].ID = 10
			Expect(client.DeleteEvent(events[0])).To(Succeed())

			Expect(server.ReceivedRequests()).To(HaveLen(1))

			Expect(cache.GetCallCount()).To(Equal(1))
		})

		testAllCommonFailures(func(c taskpkg.Repo) error {
			events[0].ID = 10
			return c.DeleteEvent(events[0])
		})
	})
})
