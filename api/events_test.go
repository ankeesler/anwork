package api_test

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/api/apifakes"
	"github.com/ankeesler/anwork/lag"
	"github.com/ankeesler/anwork/task"
	"github.com/ankeesler/anwork/task/taskfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/http_server"
)

var _ = Describe("Events", func() {
	var (
		repo          *taskfakes.FakeRepo
		authenticator *apifakes.FakeAuthenticator

		process ifrit.Process
	)

	BeforeEach(func() {
		repo = &taskfakes.FakeRepo{}
		authenticator = &apifakes.FakeAuthenticator{}

		log := log.New(GinkgoWriter, "api-test: ", 0)
		a := api.New(lag.New(log, lag.D), repo, authenticator)
		runner := http_server.New("127.0.0.1:12345", a)
		process = ifrit.Invoke(runner)
	})

	AfterEach(func() {
		process.Signal(os.Kill)
		Eventually(process.Wait()).Should(Receive())
	})

	Describe("Get", func() {
		var events []*task.Event
		BeforeEach(func() {
			events = []*task.Event{
				&task.Event{Title: "event-a", ID: 1},
				&task.Event{Title: "event-b", ID: 2},
				&task.Event{Title: "event-c", ID: 3},
			}
			repo.EventsReturnsOnCall(0, events, nil)
		})

		It("responds with the events that the repo returns", func() {
			rsp, err := get("/api/v1/events")
			Expect(err).NotTo(HaveOccurred())
			defer rsp.Body.Close()

			Expect(rsp.StatusCode).To(Equal(http.StatusOK))
			assertEvents(rsp, events)

			Expect(repo.EventsCallCount()).To(Equal(1))
		})

		Context("when getting the events fails", func() {
			BeforeEach(func() {
				repo.EventsReturnsOnCall(0, nil, errors.New("some events error"))
			})

			It("returns a 500 with an error", func() {
				rsp, err := get("/api/v1/events")
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusInternalServerError))
				assertError(rsp, "some events error")
			})
		})
	})

	Describe("Post", func() {
		var event *task.Event
		BeforeEach(func() {
			event = &task.Event{Title: "event-a", ID: 1}

			repo.CreateEventStub = func(e *task.Event) error {
				e.ID = 10
				return nil
			}
		})

		It("creates a event and responds with the location", func() {
			rsp, err := post("/api/v1/events", event)
			Expect(err).NotTo(HaveOccurred())
			defer rsp.Body.Close()

			Expect(rsp.StatusCode).To(Equal(http.StatusCreated))
			Expect(rsp.Header.Get("Location")).To(Equal("/api/v1/events/10"))

			event.ID = 10
			Expect(repo.CreateEventCallCount()).To(Equal(1))
			Expect(repo.CreateEventArgsForCall(0)).To(Equal(event))
		})

		Context("when the request payload is invalid", func() {
			It("responds with a 400 bad request", func() {
				rsp, err := post("/api/v1/events", "askjdnflkajnsfd")
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when we fail to create the event", func() {
			BeforeEach(func() {
				repo.CreateEventReturnsOnCall(0, errors.New("some create error"))
			})

			It("responds with a 500 internal server error", func() {

				rsp, err := post("/api/v1/events", event)
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusInternalServerError))
				assertError(rsp, "some create error")
			})
		})
	})
})
