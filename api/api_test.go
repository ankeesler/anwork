package api_test

import (
	"log"

	"github.com/ankeesler/anwork/api"
	"github.com/ankeesler/anwork/task/taskfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("API", func() {
	var (
		manager *taskfakes.FakeManager
		a       *api.Api

		logWriter *gbytes.Buffer
	)

	BeforeEach(func() {
		logWriter = gbytes.NewBuffer()
		l := log.New(logWriter, "api_test.go log: ", 0)
		a = api.New(manager, l)
	})

	Describe("Start", func() {
		AfterEach(func() {
			a.Stop()
		})

		It("logs that the server is starting", func() {
			errChan := make(chan error)
			Expect(a.Start(address, errChan)).To(Succeed())
			Eventually(logWriter).Should(gbytes.Say("API server starting on %s", address))
			Expect(errChan).To(BeEmpty())
		})

		Context("when the listen address is totally bad", func() {
			It("returns an error message", func() {
				errChan := make(chan error)
				Expect(a.Start("tuna", errChan)).NotTo(Succeed())
			})
		})
	})

	Describe("Stop", func() {
	})
})
