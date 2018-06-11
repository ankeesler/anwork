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
		factory *taskfakes.FakeManagerFactory
		a       *api.Api

		logWriter *gbytes.Buffer
	)

	BeforeEach(func() {
		factory = &taskfakes.FakeManagerFactory{}
		logWriter = gbytes.NewBuffer()
		l := log.New(logWriter, "api_test.go log: ", 0)
		a = api.New(address, factory, l)
	})

	It("runs", func() {
		Expect(a.Run()).To(Succeed())
	})
})
