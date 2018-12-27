package client_test

import (
	"errors"
	"net/http"

	"github.com/ankeesler/anwork/api"
	clientpkg "github.com/ankeesler/anwork/api/client"
	"github.com/ankeesler/anwork/api/client/clientfakes"
	"github.com/ankeesler/anwork/task"
	taskpkg "github.com/ankeesler/anwork/task"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/ghttp"
)

var _ = Describe("Client (auth)", func() {
	var (
		authenticator *clientfakes.FakeAuthenticator
		cache         *clientfakes.FakeCache

		client taskpkg.Repo
		server *ghttp.Server
	)

	BeforeEach(func() {
		authenticator = &clientfakes.FakeAuthenticator{}
		cache = &clientfakes.FakeCache{}

		server = ghttp.NewServer()

		client = clientpkg.New(
			makeLogger(),
			server.Addr(),
			authenticator,
			cache,
		)
	})

	AfterEach(func() {
		server.Close()
	})

	Context("when the cache is empty", func() {
		BeforeEach(func() {
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodPost, "/api/v1/auth"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.RespondWithJSONEncoded(
					http.StatusOK,
					api.Auth{Token: "some-encrypted-token"},
					http.Header{"Content-Type": {"application/json"}},
				),
			))
			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodGet, "/api/v1/tasks"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer some-token"),
				ghttp.RespondWithJSONEncoded(
					http.StatusOK,
					[]*task.Task{},
					http.Header{"Content-Type": {"application/json"}},
				),
			))

			authenticator.ValidateReturnsOnCall(0, "some-token", nil)
		})

		It("reaches out to the /api/v1/auth endpoint and caches the token", func() {
			_, err := client.Tasks()
			Expect(err).NotTo(HaveOccurred())

			Expect(server.ReceivedRequests()).To(HaveLen(2))

			Expect(authenticator.ValidateCallCount()).To(Equal(1))
			Expect(authenticator.ValidateArgsForCall(0)).To(Equal("some-encrypted-token"))

			Expect(cache.SetCallCount()).To(Equal(1))
			Expect(cache.SetArgsForCall(0)).To(Equal("some-encrypted-token"))
		})

		Context("on bad URL", func() {
			It("returns an error", func() {
				_, err := clientpkg.New(
					makeLogger(),
					"here is a bad url i mean i am sure this is bad/://aff;a;f/a;'a';sd",
					authenticator,
					cache,
				).Tasks()
				Expect(err).To(HaveOccurred())

				Expect(server.ReceivedRequests()).To(HaveLen(0))

				Expect(authenticator.ValidateCallCount()).To(Equal(0))

				Expect(cache.GetCallCount()).To(Equal(0))
			})
		})

		Context("on a failed request", func() {
			It("returns an error", func() {
				_, err := clientpkg.New(
					makeLogger(),
					"asdf",
					authenticator,
					cache,
				).Tasks()
				Expect(err).To(HaveOccurred())

				Expect(server.ReceivedRequests()).To(HaveLen(0))

				Expect(authenticator.ValidateCallCount()).To(Equal(0))

				Expect(cache.GetCallCount()).To(Equal(1))
			})
		})

		Context("on 4xx response", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodPost, "/api/v1/auth"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.RespondWithJSONEncoded(
						http.StatusBadRequest,
						api.Error{Message: "some error message"},
						http.Header{"Content-Type": {"application/json"}},
					),
				))
			})

			It("returns an error", func() {
				_, err := client.Tasks()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("400 Bad Request"))

				Expect(server.ReceivedRequests()).To(HaveLen(1))

				Expect(authenticator.ValidateCallCount()).To(Equal(0))

				Expect(cache.GetCallCount()).To(Equal(1))
			})
		})

		Context("on 5xx response", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodPost, "/api/v1/auth"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.RespondWithJSONEncoded(
						http.StatusInternalServerError,
						api.Error{Message: "some error message"},
						http.Header{"Content-Type": {"application/json"}},
					),
				))
			})

			It("returns an error", func() {
				_, err := client.Tasks()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("500 Internal Server Error"))
				Expect(err.Error()).To(ContainSubstring("some error message"))

				Expect(server.ReceivedRequests()).To(HaveLen(1))

				Expect(authenticator.ValidateCallCount()).To(Equal(0))

				Expect(cache.GetCallCount()).To(Equal(1))
			})
		})

		Context("on bogus response payload", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodPost, "/api/v1/auth"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
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

				Expect(authenticator.ValidateCallCount()).To(Equal(0))

				Expect(cache.GetCallCount()).To(Equal(1))
			})
		})

		Context("when validation of token fails", func() {
			BeforeEach(func() {
				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodPost, "/api/v1/auth"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.RespondWithJSONEncoded(
						http.StatusOK,
						api.Auth{Token: "some-token"},
						http.Header{"Content-Type": {"application/json"}},
					),
				))

				authenticator.ValidateReturnsOnCall(0, "", errors.New("some validate error"))
			})

			It("returns an error", func() {
				_, err := client.Tasks()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("some validate error"))

				Expect(server.ReceivedRequests()).To(HaveLen(1))

				Expect(authenticator.ValidateCallCount()).To(Equal(1))

				Expect(cache.GetCallCount()).To(Equal(1))
			})
		})
	})

	Context("when the cache is not empty", func() {
		BeforeEach(func() {
			cache.GetReturnsOnCall(0, "some-encrypted-token", true)
			authenticator.ValidateReturnsOnCall(0, "some-token", nil)

			server.AppendHandlers(ghttp.CombineHandlers(
				ghttp.VerifyRequest(http.MethodGet, "/api/v1/tasks"),
				ghttp.VerifyHeaderKV("Accept", "application/json"),
				ghttp.VerifyHeaderKV("Authorization", "bearer some-token"),
				ghttp.RespondWithJSONEncoded(
					http.StatusOK,
					[]*task.Task{},
					http.Header{"Content-Type": {"application/json"}},
				),
			))
		})

		It("does not reach out to the /api/v1/auth endpoint and validates the token", func() {
			_, err := client.Tasks()
			Expect(err).NotTo(HaveOccurred())

			Expect(server.ReceivedRequests()).To(HaveLen(1))

			Expect(authenticator.ValidateCallCount()).To(Equal(1))
			Expect(authenticator.ValidateArgsForCall(0)).To(Equal("some-encrypted-token"))

			Expect(cache.SetCallCount()).To(Equal(0))
		})

		Context("when the authenticator fails to validate the token", func() {
			BeforeEach(func() {
				authenticator.ValidateReturnsOnCall(0, "", errors.New("some validate error"))
				authenticator.ValidateReturnsOnCall(1, "some-new-token", nil)

				server.SetHandler(0, ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodPost, "/api/v1/auth"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.RespondWithJSONEncoded(
						http.StatusOK,
						api.Auth{Token: "some-new-encrypted-token"},
						http.Header{"Content-Type": {"application/json"}},
					),
				))
				server.AppendHandlers(ghttp.CombineHandlers(
					ghttp.VerifyRequest(http.MethodGet, "/api/v1/tasks"),
					ghttp.VerifyHeaderKV("Accept", "application/json"),
					ghttp.VerifyHeaderKV("Authorization", "bearer some-new-token"),
					ghttp.RespondWithJSONEncoded(
						http.StatusOK,
						[]*task.Task{},
						http.Header{"Content-Type": {"application/json"}},
					),
				))
			})

			It("reaches out to /api/v1/auth to get a new token", func() {
				_, err := client.Tasks()
				Expect(err).NotTo(HaveOccurred())

				Expect(server.ReceivedRequests()).To(HaveLen(2))

				Expect(authenticator.ValidateCallCount()).To(Equal(2))
				Expect(authenticator.ValidateArgsForCall(0)).To(Equal("some-encrypted-token"))
				Expect(authenticator.ValidateArgsForCall(1)).To(Equal("some-new-encrypted-token"))

				Expect(cache.SetCallCount()).To(Equal(1))
				Expect(cache.SetArgsForCall(0)).To(Equal("some-new-encrypted-token"))
			})
		})
	})
})
