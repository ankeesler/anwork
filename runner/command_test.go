package runner_test

import (
	"github.com/ankeesler/anwork/runner"
	"github.com/ankeesler/anwork/task/taskfakes"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Command", func() {
	var (
		r                         *runner.Runner
		factory                   *taskfakes.FakeManagerFactory
		stdoutWriter, debugWriter *gbytes.Buffer
		manager                   *taskfakes.FakeManager
	)

	BeforeEach(func() {
		manager = &taskfakes.FakeManager{}

		factory = &taskfakes.FakeManagerFactory{}
		factory.CreateReturnsOnCall(0, manager, nil)

		stdoutWriter = gbytes.NewBuffer()
		debugWriter = gbytes.NewBuffer()

		r = runner.New(factory, stdoutWriter, debugWriter)
	})

	run := func(args ...string) {
		Expect(r.Run(args)).To(Succeed())
	}

	Describe("create", func() {
		It("calls the manager to create a task", func() {
			run("create", "task-a")
			Expect(manager.CreateCallCount()).To(Equal(1))
			Expect(manager.CreateArgsForCall(0)).To(Equal("task-a"))
		})
	})

	Describe("reset", func() {
	})
})
