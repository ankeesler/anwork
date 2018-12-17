package runner_test

import (
	"errors"

	"github.com/ankeesler/anwork/manager/managerfakes"
	"github.com/ankeesler/anwork/runner"
	"github.com/ankeesler/anwork/task"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("AnworkRunner", func() {
	var (
		r                         *runner.Runner
		manager                   *managerfakes.FakeManager
		stdoutWriter, debugWriter *gbytes.Buffer
	)

	BeforeEach(func() {
		manager = &managerfakes.FakeManager{}

		stdoutWriter = gbytes.NewBuffer()
		debugWriter = gbytes.NewBuffer()

		bi := &runner.BuildInfo{
			Hash: "whatever",
			Date: "don't care",
		}
		r = runner.New(bi, manager, stdoutWriter, debugWriter)
	})

	Context("when a valid command is issued and returns successfully", func() {
		It("writes a helpful debug message", func() {
			err := r.Run([]string{"create", "task-a"})
			Expect(err).NotTo(HaveOccurred())
			Expect(debugWriter).To(gbytes.Say("Manager is"))
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

	Context("when a command takes an optional argument", func() {
		BeforeEach(func() {
			manager.FindByNameReturnsOnCall(0, &task.Task{}, nil)
		})

		It("allows the optional argument not to be passed", func() {
			Expect(r.Run([]string{"show"})).To(Succeed())
		})

		It("allows the optional argument to be passed", func() {
			Expect(r.Run([]string{"show", "tuna"})).To(Succeed())
		})

		It("fails if more than the optional argument is passed", func() {
			err := r.Run([]string{"show", "tuna", "fish"})
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Invalid argument passed to command 'show'"))
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

	Context("when shortcut commands are passed", func() {
		Context("when the shortcut command is invalid", func() {
			It("fails with an unknown command error", func() {
				err := r.Run([]string{"sq"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Unknown command: 'sq'"))
			})
		})
		Context("when the shortcut command is valid", func() {
			It("runs the commands as expected", func() {
				manager.FindByNameReturnsOnCall(0, &task.Task{Name: "some-task"}, nil)
				err := r.Run([]string{"sr", "some-task"})
				Expect(err).NotTo(HaveOccurred())
			})

			Context("when the shortcut command has an argument passed to it", func() {
				It("returns a helpful error with the actual command name", func() {
					err := r.Run([]string{"sr"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("Invalid argument passed to command 'set-running'"))
				})
			})
		})
	})

	Describe("Usage", func() {
		It("prints the usage information for every command in a command line format", func() {
			buffer := gbytes.NewBuffer()
			runner.Usage(buffer)
			Expect(buffer).To(gbytes.Say("  create task-name"))
			Expect(buffer).To(gbytes.Say("   Create a new task \\(alias: c\\)"))
			Expect(buffer).To(gbytes.Say("  show \\[task-name\\]"))
			Expect(buffer).To(gbytes.Say("   Show the current tasks, or the details of a specific task \\(alias: s\\)"))
			Expect(buffer).To(gbytes.Say("  set-running task-name"))
			Expect(buffer).To(gbytes.Say("   Mark a task as running \\(alias: sr\\)"))
			Expect(buffer).To(gbytes.Say("  set-ready task-name"))
			Expect(buffer).To(gbytes.Say("   Mark a task as ready \\(alias: sy\\)"))
		})
	})

	Describe("MarkdownUsage", func() {
		It("prints the usage information for every command in a github markdown format", func() {
			buffer := gbytes.NewBuffer()
			runner.MarkdownUsage(buffer)
			Expect(buffer).To(gbytes.Say("### `anwork create task-name`"))
			Expect(buffer).To(gbytes.Say("\\* Create a new task\n"))
			Expect(buffer).To(gbytes.Say("\\* Alias: `c`"))
			Expect(buffer).To(gbytes.Say("### `anwork show \\[task-name\\]`"))
			Expect(buffer).To(gbytes.Say("\\* Show the current tasks, or the details of a specific task\n"))
			Expect(buffer).To(gbytes.Say("\\* Alias: `s`"))
			Expect(buffer).To(gbytes.Say("### `anwork set-running task-name`"))
			Expect(buffer).To(gbytes.Say("\\* Mark a task as running\n"))
			Expect(buffer).To(gbytes.Say("\\* Alias: `sr`"))
			Expect(buffer).To(gbytes.Say("### `anwork set-ready task-name`"))
			Expect(buffer).To(gbytes.Say("\\* Mark a task as ready\n"))
			Expect(buffer).To(gbytes.Say("\\* Alias: `sy`"))
		})
	})
})
