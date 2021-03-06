package api_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"code.cloudfoundry.org/lager/lagertest"
	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/api/apifakes"
	"github.com/ankeesler/anwork/task/taskfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/http_server"
)

var _ = Describe("Tasks", func() {
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

	Describe("Post", func() {
		BeforeEach(func() {
			authenticator.TokenReturnsOnCall(0, "here is a token", nil)
		})

		It("returns a token from the authenticator", func() {
			rsp, err := post("/api/v1/auth", nil)
			Expect(err).NotTo(HaveOccurred())
			defer rsp.Body.Close()

			Expect(rsp.StatusCode).To(Equal(http.StatusOK))

			var token api.Auth
			Expect(json.NewDecoder(rsp.Body).Decode(&token)).To(Succeed())
			Expect(token).To(Equal(api.Auth{Token: "here is a token"}))

			Expect(authenticator.TokenCallCount()).To(Equal(1))
		})

		Context("when the authenticator fails", func() {
			BeforeEach(func() {
				authenticator.TokenReturnsOnCall(0, "", errors.New("some auth error"))
			})

			It("returns an error", func() {
				rsp, err := post("/api/v1/auth", nil)
				Expect(err).NotTo(HaveOccurred())
				defer rsp.Body.Close()

				Expect(rsp.StatusCode).To(Equal(http.StatusInternalServerError))
				assertError(rsp, "some auth error")

				Expect(authenticator.TokenCallCount()).To(Equal(1))
			})
		})
	})
})
