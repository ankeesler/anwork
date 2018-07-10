package api_test

import (
	"context"
	"errors"
	"io"
	"log"
	"net"

	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/task/taskfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("API", func() {
	var (
		factory *taskfakes.FakeManagerFactory
		a       *api.Api

		logWriter *gbytes.Buffer

		ctx, cancel = context.WithCancel(context.Background())
	)

	BeforeEach(func() {
		factory = &taskfakes.FakeManagerFactory{}
		manager := &taskfakes.FakeManager{}
		factory.CreateReturnsOnCall(0, manager, nil)

		logWriter = gbytes.NewBuffer()
		l := log.New(io.MultiWriter(logWriter, GinkgoWriter), "api_test.go log: ", log.Ldate|log.Ltime|log.Lshortfile)

		a = api.New(address, factory, l)
	})

	Context("when creating a manager fails", func() {
		BeforeEach(func() {
			factory.CreateReturnsOnCall(0, nil, errors.New("some factory error"))
		})

		It("returns the error", func() {
			Expect(a.Run(ctx)).To(MatchError("some factory error"))
		})

		It("logs an error", func() {
			a.Run(ctx)
			Eventually(logWriter).Should(gbytes.Say("failed to make server:"))
		})
	})

	Context("when listening on the address fails", func() {
		var listener net.Listener

		BeforeEach(func() {
			var err error
			listener, err = net.Listen("tcp", address)
			Expect(err).NotTo(HaveOccurred())
		})

		AfterEach(func() {
			Expect(listener.Close()).To(Succeed())
		})

		It("returns an error", func() {
			err := a.Run(ctx)
			Expect(err).To(HaveOccurred())
		})

		It("logs an error", func() {
			a.Run(ctx)
			Eventually(logWriter).Should(gbytes.Say("failed to listen on address %s:", address))
		})
	})

	Describe("Run", func() {
		BeforeEach(func() {
			Expect(a.Run(ctx)).To(Succeed())
		})

		AfterEach(func() {
			cancel()
			Eventually(logWriter).Should(gbytes.Say("listener closed"))
		})

		It("starts a server on the provided address", func() {
			_, err := get("/")
			Expect(err).NotTo(HaveOccurred())
		})

		It("logs that it started", func() {
			Eventually(logWriter).Should(gbytes.Say("API server starting on %s", address))
		})
	})
})
