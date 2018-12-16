package api_test

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

var _ = Describe("API", func() {
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
		XIt("returns a 405", func() {
			rsp, err := post("/api/v1/tasks/1")
			Expect(err).NotTo(HaveOccurred())
			defer rsp.Body.Close()

			Expect(rsp.StatusCode).To(Equal(http.StatusMethodNotAllowed))
		})
	})
})

func get(path string) (*http.Response, error) {
	return http.Get(fmt.Sprintf("http://127.0.0.1:12345%s", path))
}

func post(path string) (*http.Response, error) {
	return http.Post(fmt.Sprintf("http://127.0.0.1:12345%s", path), "", nil)
}

func assertError(rsp *http.Response, message string) {
	bytes, err := ioutil.ReadAll(rsp.Body)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	var errMsg api.Error
	ExpectWithOffset(1, json.Unmarshal(bytes, &errMsg)).NotTo(HaveOccurred())

	ExpectWithOffset(1, errMsg.Message).To(Equal(message))
}

func assertTasks(rsp *http.Response, tasks []*task2.Task) {
	bytes, err := ioutil.ReadAll(rsp.Body)
	ExpectWithOffset(1, err).NotTo(HaveOccurred())

	actualTasks := make([]*task2.Task, 1)
	ExpectWithOffset(1, json.Unmarshal(bytes, &tasks)).NotTo(HaveOccurred())

	ExpectWithOffset(1, tasks).To(Equal(actualTasks))
}
