package api_test

import (
	"errors"
	"log"
	"net/http"
	"os"

	api "github.com/ankeesler/anwork/api2"
	"github.com/ankeesler/anwork/api2/api2fakes"
	"github.com/ankeesler/anwork/task2"
	"github.com/ankeesler/anwork/task2/task2fakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/http_server"
)

var _ = Describe("Events", func() {
	var (
		repo          *task2fakes.FakeRepo
		authenticator *api2fakes.FakeAuthenticator

		process ifrit.Process
	)

	BeforeEach(func() {
		repo = &task2fakes.FakeRepo{}
		authenticator = &api2fakes.FakeAuthenticator{}

		a := api.New(log.New(GinkgoWriter, "api-test: ", 0), repo, authenticator)
		runner := http_server.New("127.0.0.1:12345", a)
		process = ifrit.Invoke(runner)
	})

	AfterEach(func() {
		process.Signal(os.Kill)
		Eventually(process.Wait()).Should(Receive())
	})

	Describe("Get", func() {
		var events []*task2.Event
		BeforeEach(func() {
			events = []*task2.Event{
				&task2.Event{Title: "event-a", ID: 1},
				&task2.Event{Title: "event-b", ID: 2},
				&task2.Event{Title: "event-c", ID: 3},
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
		var event *task2.Event
		BeforeEach(func() {
			event = &task2.Event{Title: "event-a", ID: 1}

			repo.CreateEventStub = func(e *task2.Event) error {
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
