package runner_test

import (
	"errors"
	"os"

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

	Describe("create", func() {
		It("calls the manager to create a task", func() {
			Expect(r.Run([]string{"create", "task-a"})).To(Succeed())
			Expect(manager.CreateCallCount()).To(Equal(1))
			Expect(manager.CreateArgsForCall(0)).To(Equal("task-a"))
		})

		Context("when the manager fails to create a task", func() {
			BeforeEach(func() {
				factory.CreateReturnsOnCall(1, manager, nil)
				manager.CreateReturnsOnCall(1, errors.New("failed to create task"))
			})

			It("returns a helpful error message", func() {
				Expect(r.Run([]string{"create", "task-a"})).To(Succeed())

				err := r.Run([]string{"create", "task-a"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to create task"))
			})
		})
	})

	Describe("reset", func() {
		Context("when the ANWORK_TEST_RESET_ANSWER environmental variable is set to 'y'", func() {
			var envVarBefore string

			BeforeEach(func() {
				envVarBefore = os.Getenv("ANWORK_TEST_RESET_ANSWER")
				Expect(os.Setenv("ANWORK_TEST_RESET_ANSWER", "y")).To(Succeed())
			})

			AfterEach(func() {
				Expect(os.Setenv("ANWORK_TEST_RESET_ANSWER", envVarBefore)).To(Succeed())
			})

			It("asks the user to confirm and tells them their data is being deleted", func() {
				r.Run([]string{"reset"})
				Eventually(stdoutWriter).Should(gbytes.Say("Are you sure you want to delete all data \\[y/n\\]: "))
				Eventually(stdoutWriter).Should(gbytes.Say("OK, deleting all data"))
			})

			It("tells the factory to reset", func() {
				Expect(r.Run([]string{"reset"})).To(Succeed())
				Expect(factory.ResetCallCount()).To(Equal(1))
			})
		})

		Context("when the ANWORK_TEST_RESET_ANSWER environmental variable is set to something other than 'y'", func() {
			var envVarBefore string

			BeforeEach(func() {
				envVarBefore = os.Getenv("ANWORK_TEST_RESET_ANSWER")
				Expect(os.Setenv("ANWORK_TEST_RESET_ANSWER", "tuna")).To(Succeed())
			})

			AfterEach(func() {
				Expect(os.Setenv("ANWORK_TEST_RESET_ANSWER", envVarBefore)).To(Succeed())
			})

			It("asks the user to confirm and tells them their data is not being deleted", func() {
				Expect(r.Run([]string{"reset"})).To(Succeed())
				Eventually(stdoutWriter).Should(gbytes.Say("Are you sure you want to delete all data \\[y/n\\]: "))
				Eventually(stdoutWriter).Should(gbytes.Say("NOT deleting all data"))
			})

			It("does not tell the factory to reset", func() {
				Expect(r.Run([]string{"reset"})).To(Succeed())
				Expect(factory.ResetCallCount()).To(Equal(0))
			})
		})
	})
})
