package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/api/apifakes"
	taskpkg "github.com/ankeesler/anwork/task"
	"github.com/ankeesler/anwork/task/taskfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/http_server"
)

var _ = Describe("API", func() {
	var (
		repo          *taskfakes.FakeRepo
		authenticator *apifakes.FakeAuthenticator

		process ifrit.Process
	)

	BeforeEach(func() {
		repo = &taskfakes.FakeRepo{}
		authenticator = &apifakes.FakeAuthenticator{}

		a := api.New(lagertest.NewTestLogger("api"), repo, authenticator)
		runner := http_server.New("127.0.0.1:12345", a)
		process = ifrit.Invoke(runner)
	})

	AfterEach(func() {
		process.Signal(os.Kill)
		Eventually(process.Wait()).Should(Receive())
	})

	Context("authentication", func() {
		It("doesn't call authenticate() on the /api/v1/auth endpoint", func() {
			rsp, err := post("/api/v1/auth", nil)
			Expect(err).NotTo(HaveOccurred())
			defer rsp.Body.Close()

			Expect(authenticator.AuthenticateCallCount()).To(Equal(0))
		})

		It("doesn't call authenticate() on the /api/v1/health endpoint", func() {
			rsp, err := post("/api/v1/health", nil)
			Expect(err).NotTo(HaveOccurred())
			defer rsp.Body.Close()

			Expect(authenticator.AuthenticateCallCount()).To(Equal(0))
		})

		It("passes the bearer token to the authenticator", func() {
			rsp, err := get("")
			Expect(err).NotTo(HaveOccurred())
			defer rsp.Body.Close()

			Expect(authenticator.AuthenticateCallCount()).To(Equal(1))
			Expect(authenticator.AuthenticateArgsForCall(0)).To(Equal("some-token"))
		})

		Context("no Authorization header is included in the request", func() {
			It("fails with a 401 and an error message", func() {
				url := fmt.Sprintf("http://127.0.0.1:12345")

				req, err := http.NewRequest(http.MethodPut, url, nil)
				Expect(err).NotTo(HaveOccurred())

				rsp, err := http.DefaultClient.Do(req)
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusUnauthorized))
				assertError(rsp, "missing authorization header")

				Expect(authenticator.AuthenticateCallCount()).To(Equal(0))
			})
		})

		Context("the token is incorrectly formatted", func() {
			It("fails with a 400 and an error message", func() {
				url := fmt.Sprintf("http://127.0.0.1:12345")

				req, err := http.NewRequest(http.MethodPut, url, nil)
				Expect(err).NotTo(HaveOccurred())

				req.Header.Set("Authorization", "bearerasdfasdfasdf")

				rsp, err := http.DefaultClient.Do(req)
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusBadRequest))
				assertError(rsp, "invalid authorization data")

				Expect(authenticator.AuthenticateCallCount()).To(Equal(0))
			})
		})

		Context("the token is not of type bearer", func() {
			It("fails with a 400 and an error message", func() {
				url := fmt.Sprintf("http://127.0.0.1:12345")

				req, err := http.NewRequest(http.MethodPut, url, nil)
				Expect(err).NotTo(HaveOccurred())

				req.Header.Set("Authorization", "pancake asdfasdfasdf")

				rsp, err := http.DefaultClient.Do(req)
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusBadRequest))
				assertError(rsp, "invalid authorization data")

				Expect(authenticator.AuthenticateCallCount()).To(Equal(0))
			})
		})

		Context("failed authentication", func() {
			BeforeEach(func() {
				authenticator.AuthenticateReturnsOnCall(0, errors.New("some auth error"))
			})

			It("returns an error and a 403", func() {
				rsp, err := get("")
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusForbidden))
				assertError(rsp, "some auth error")

				Expect(authenticator.AuthenticateCallCount()).To(Equal(1))
			})
		})
	})

	Context("path not found", func() {
		It("returns a 404", func() {
			rsp, err := get("/alskjdnflkajnsdflkajsndf")
			Expect(err).NotTo(HaveOccurred())
			defer rsp.Body.Close()

			Expect(rsp.StatusCode).To(Equal(http.StatusNotFound))
		})
	})

	Context("method not allowed", func() {
		It("returns a 405", func() {
			rsp, err := deletee("/api/v1/tasks")
			Expect(err).NotTo(HaveOccurred())
			defer rsp.Body.Close()

			Expect(rsp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
		})
	})
})

func get(path string) (*http.Response, error) {
	return do(http.MethodGet, path, nil)
}

func put(path string, body interface{}) (*http.Response, error) {
	return do(http.MethodPut, path, body)
}

func post(path string, body interface{}) (*http.Response, error) {
	return do(http.MethodPost, path, body)
}

func deletee(path string) (*http.Response, error) {
	return do(http.MethodDelete, path, nil)
}

func do(method, path string, body interface{}) (*http.Response, error) {
	url := fmt.Sprintf("http://127.0.0.1:12345%s", path)

	var data []byte
	if body != nil {
		var err error
		data, err = json.Marshal(body)
		ExpectWithOffset(1, err).NotTo(HaveOccurred())
	}

	buf := bytes.NewBuffer(data)
	req, err := http.NewRequest(method, url, buf)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	req.Header.Set("Authorization", "bearer some-token")

	return http.DefaultClient.Do(req)
}

func assertError(rsp *http.Response, message string) {
	bytes, err := ioutil.ReadAll(rsp.Body)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	var errMsg api.Error
	ExpectWithOffset(1, json.Unmarshal(bytes, &errMsg)).NotTo(HaveOccurred())

	ExpectWithOffset(1, errMsg.Message).To(Equal(message))
}

func assertTasks(rsp *http.Response, tasks []*taskpkg.Task) {
	bytes, err := ioutil.ReadAll(rsp.Body)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	actualTasks := make([]*taskpkg.Task, 1)
	ExpectWithOffset(1, json.Unmarshal(bytes, &actualTasks)).NotTo(HaveOccurred())

	ExpectWithOffset(1, actualTasks).To(Equal(tasks))
}

func assertTask(rsp *http.Response, task *taskpkg.Task) {
	bytes, err := ioutil.ReadAll(rsp.Body)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	var actualTask taskpkg.Task
	ExpectWithOffset(1, json.Unmarshal(bytes, &actualTask)).NotTo(HaveOccurred())

	ExpectWithOffset(1, actualTask).To(Equal(*task))
}

func assertEvents(rsp *http.Response, tasks []*taskpkg.Event) {
	bytes, err := ioutil.ReadAll(rsp.Body)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	actualEvents := make([]*taskpkg.Event, 1)
	ExpectWithOffset(1, json.Unmarshal(bytes, &actualEvents)).NotTo(HaveOccurred())

	ExpectWithOffset(1, actualEvents).To(Equal(tasks))
}

func assertEvent(rsp *http.Response, task *taskpkg.Event) {
	bytes, err := ioutil.ReadAll(rsp.Body)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	var actualEvent taskpkg.Event
	ExpectWithOffset(1, json.Unmarshal(bytes, &actualEvent)).NotTo(HaveOccurred())

	ExpectWithOffset(1, actualEvent).To(Equal(*task))
}
