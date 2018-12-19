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

var _ = Describe("Event", func() {
	var (
		repo          *taskfakes.FakeRepo
		authenticator *apifakes.FakeAuthenticator

		process ifrit.Process
	)

	testAllCommonFailures := func(doFunc func(path string) (*http.Response, error)) {
		Context("when the id in the path is invalid", func() {
			It("returns with 400 bad request", func() {
				rsp, err := doFunc("/api/v1/events/tuna")
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusBadRequest))
			})
		})

		Context("when the repo fails to get the event", func() {
			BeforeEach(func() {
				repo.FindEventByIDReturnsOnCall(0, nil, errors.New("some find error"))
			})

			It("responds with a 500 internal server error plus the error", func() {
				rsp, err := doFunc("/api/v1/events/10")
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusInternalServerError))
				assertError(rsp, "some find error")
			})
		})

		Context("when the event does not exist", func() {
			BeforeEach(func() {
				repo.FindEventByIDReturnsOnCall(0, nil, nil)
			})

			It("responds with a 404 not found", func() {
				rsp, err := doFunc("/api/v1/events/10")
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusNotFound))
			})
		})
	}

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
		var event *task.Event
		BeforeEach(func() {
			event = &task.Event{Title: "event-a", ID: 1}
			repo.FindEventByIDReturnsOnCall(0, event, nil)
		})

		It("responds with the event", func() {
			rsp, err := get("/api/v1/events/10")
			Expect(err).NotTo(HaveOccurred())
			defer rsp.Body.Close()

			Expect(rsp.StatusCode).To(Equal(http.StatusOK))
			assertEvent(rsp, event)

			Expect(repo.FindEventByIDCallCount()).To(Equal(1))
			Expect(repo.FindEventByIDArgsForCall(0)).To(Equal(10))
		})

		testAllCommonFailures(get)
	})

	Describe("Delete", func() {
		var event *task.Event
		BeforeEach(func() {
			event = &task.Event{Title: "event-a", ID: 1}
			repo.FindEventByIDReturnsOnCall(0, event, nil)
		})

		It("deletes the event", func() {
			rsp, err := deletee("/api/v1/events/10")
			Expect(err).NotTo(HaveOccurred())
			defer rsp.Body.Close()

			Expect(rsp.StatusCode).To(Equal(http.StatusNoContent))

			Expect(repo.FindEventByIDCallCount()).To(Equal(1))
			Expect(repo.FindEventByIDArgsForCall(0)).To(Equal(10))

			Expect(repo.DeleteEventCallCount()).To(Equal(1))
			Expect(repo.DeleteEventArgsForCall(0)).To(Equal(event))
		})

		Context("when the repo fails to delete the event", func() {
			BeforeEach(func() {
				repo.DeleteEventReturnsOnCall(0, errors.New("some delete failure"))
			})

			It("responds with a 500 and the error", func() {
				rsp, err := deletee("/api/v1/events/10")
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusInternalServerError))
				assertError(rsp, "some delete failure")
			})
		})

		testAllCommonFailures(deletee)
	})
})
