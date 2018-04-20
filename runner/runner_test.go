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

var _ = Describe("AnworkRunner", func() {
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

	Context("when a valid command is issued and returns successfully", func() {
		It("tells the factory to save the manager", func() {
			err := r.Run([]string{"create", "task-a"})
			Expect(err).NotTo(HaveOccurred())
			Expect(factory.SaveCallCount()).To(Equal(1))
			Expect(factory.SaveArgsForCall(0)).To(Equal(manager))
		})

		It("writes a helpful debug message", func() {
			err := r.Run([]string{"create", "task-a"})
			Expect(err).NotTo(HaveOccurred())
			Expect(debugWriter).To(gbytes.Say("Manager is"))
		})

		Context("when the factory fails to save the manager", func() {
			BeforeEach(func() {
				factory.SaveReturnsOnCall(0, errors.New("some error"))
			})

			It("returns a helpful error message", func() {
				err := r.Run([]string{"create", "task-a"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Could not save manager: some error"))
			})
		})
	})

	Context("when incorrect arguments are passed to a command", func() {
		It("returns a helpful error", func() {
			err := r.Run([]string{"create"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Invalid argument passed to command 'create'"))
			Expect(err.Error()).To(ContainSubstring("Got: []"))
			Expect(err.Error()).To(ContainSubstring("Expected: [task-name]"))
		})
	})

	Context("when the command fails", func() {
		BeforeEach(func() {
			manager.CreateReturnsOnCall(0, errors.New("some error"))
		})
		It("returns the command's error", func() {
			err := r.Run([]string{"create", "task-a"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Command 'create' failed: some error"))
		})
	})

	Context("when an invalid command (e.g., 'tuna') is passed", func() {
		It("returns a helpful error", func() {
			err := r.Run([]string{"tuna"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Unknown command: 'tuna'"))
		})
	})

	Context("when the factory fails to create a manager", func() {
		BeforeEach(func() {
			factory.CreateReturnsOnCall(0, nil, errors.New("some error"))
		})

		It("returns a helpful error", func() {
			err := r.Run([]string{"create", "task-a"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("Could not create manager: some error"))
		})
	})

	Context("when the command requires that the manager factory be reset", func() {
		var envVarBefore string

		BeforeEach(func() {
			envVarBefore = os.Getenv("ANWORK_TEST_RESET_ANSWER")
			Expect(os.Setenv("ANWORK_TEST_RESET_ANSWER", "y")).To(Succeed())
		})

		AfterEach(func() {
			Expect(os.Setenv("ANWORK_TEST_RESET_ANSWER", envVarBefore)).To(Succeed())
		})

		It("calls reset on the manager", func() {
			err := r.Run([]string{"reset"})
			Expect(err).NotTo(HaveOccurred())
			Expect(factory.ResetCallCount()).To(Equal(1))
		})

		Context("when the factory fails to reset", func() {
			BeforeEach(func() {
				factory.ResetReturnsOnCall(0, errors.New("some error"))
			})
			It("returns the error", func() {
				err := r.Run([]string{"reset"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Could not reset factory: some error"))
			})
		})
	})
})
