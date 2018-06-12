package api_test

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

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
		logWriter = gbytes.NewBuffer()
		l := log.New(io.MultiWriter(logWriter, GinkgoWriter), "api_test.go log: ", 0)
		a = api.New(address, factory, l)
	})

	Context("context", func() {
		BeforeEach(func() {
			a.Run(ctx)
		})

		AfterEach(func() {
			cancel()
			Eventually(logWriter).Should(gbytes.Say("API server successfully closed listener socket"))
		})

		It("runs", func() {
			_, err := http.Get(fmt.Sprintf("http://%s/", address))
			Expect(err).NotTo(HaveOccurred())
		})

	})
})
