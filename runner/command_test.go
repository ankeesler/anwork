package runner_test

import (
	"errors"
	"fmt"
	"os"

	"github.com/ankeesler/anwork/runner"
	"github.com/ankeesler/anwork/task"
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
		factory.CreateReturns(manager, nil)

		stdoutWriter = gbytes.NewBuffer()
		debugWriter = gbytes.NewBuffer()

		r = runner.New(factory, stdoutWriter, debugWriter)
	})

	Describe("version", func() {
		It("prints out the version", func() {
			Expect(r.Run([]string{"version"})).To(Succeed())
			Expect(stdoutWriter).To(gbytes.Say("ANWORK Version = 3"))
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

	Describe("create", func() {
		It("calls the manager to create a task", func() {
			Expect(r.Run([]string{"create", "task-a"})).To(Succeed())
			Expect(manager.CreateCallCount()).To(Equal(1))
			Expect(manager.CreateArgsForCall(0)).To(Equal("task-a"))
		})

		Context("when the manager fails to create a task", func() {
			BeforeEach(func() {
				manager.CreateReturnsOnCall(0, errors.New("failed to create task"))
			})

			It("returns a helpful error message", func() {
				err := r.Run([]string{"create", "task-a"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to create task"))
			})
		})
	})

	Describe("delete", func() {
		Context("when the manager successfully deletes the task", func() {
			BeforeEach(func() {
				manager.DeleteReturnsOnCall(0, true)
			})

			It("successfully deletes existing tasks", func() {
				Expect(r.Run([]string{"delete", "task-a"})).To(Succeed())
				Expect(manager.DeleteCallCount()).To(Equal(1))
				Expect(manager.DeleteArgsForCall(0)).To(Equal("task-a"))
			})
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				manager.DeleteReturnsOnCall(0, false)
			})

			It("prints out a helpful message saying that the task was unknown", func() {
				err := r.Run([]string{"delete", "task-a"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unknown task: task-a"))
			})
		})
	})

	Describe("delete-all", func() {
		Context("when there are no tasks", func() {
			It("does not tell the manager to delete anything", func() {
				Expect(r.Run([]string{"delete-all"})).To(Succeed())
				Expect(manager.DeleteCallCount()).To(Equal(0))
			})
		})

		Context("when there are multiple tasks", func() {
			BeforeEach(func() {
				tasks := []*task.Task{
					&task.Task{Name: "task-a"},
					&task.Task{Name: "task-b"},
					&task.Task{Name: "task-c"},
				}
				manager.TasksReturnsOnCall(0, tasks)
			})

			Context("when the manager successfully deletes each task", func() {
				BeforeEach(func() {
					manager.DeleteReturnsOnCall(0, true)
					manager.DeleteReturnsOnCall(1, true)
					manager.DeleteReturnsOnCall(2, true)
				})

				It("calls delete on each task", func() {
					Expect(r.Run([]string{"delete-all"})).To(Succeed())
					Expect(manager.DeleteCallCount()).To(Equal(3))
					Expect(manager.DeleteArgsForCall(0)).To(Equal("task-a"))
					Expect(manager.DeleteArgsForCall(1)).To(Equal("task-b"))
					Expect(manager.DeleteArgsForCall(2)).To(Equal("task-c"))
				})
			})

			Context("when the manager fails to delete a task", func() {
				BeforeEach(func() {
					manager.DeleteReturnsOnCall(0, false)
					manager.DeleteReturnsOnCall(1, true)
					manager.DeleteReturnsOnCall(2, false)
				})

				It("notifies the user of which task was not able to be deleted", func() {
					err := r.Run([]string{"delete-all"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("unable to delete task task-a"))
					Expect(err.Error()).To(ContainSubstring("unable to delete task task-c"))

					Expect(manager.DeleteCallCount()).To(Equal(3))
					Expect(manager.DeleteArgsForCall(0)).To(Equal("task-a"))
					Expect(manager.DeleteArgsForCall(1)).To(Equal("task-b"))
					Expect(manager.DeleteArgsForCall(2)).To(Equal("task-c"))
				})
			})

		})
	})

	Describe("show", func() {
		Context("when there are no tasks", func() {
			It("just prints out the task states", func() {
				Expect(r.Run([]string{"show"})).To(Succeed())
				Expect(stdoutWriter).To(gbytes.Say("RUNNING tasks:\nBLOCKED tasks:\nWAITING tasks:\nFINISHED tasks:\n"))
			})

			Context("when a task name argument is passed", func() {
				It("prints an error an unknown task", func() {
					err := r.Run([]string{"show", "task-a"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("unknown task: task-a"))
				})
			})
		})

		Context("when there are multiple tasks", func() {
			BeforeEach(func() {
				tasks := []*task.Task{
					&task.Task{
						Name:  "task-a",
						ID:    10,
						State: task.StateRunning,
					},
					&task.Task{
						Name:     "task-b",
						Priority: 3,
						ID:       20,
						State:    task.StateWaiting,
					},
					&task.Task{
						Name:  "task-c",
						ID:    30,
						State: task.StateWaiting,
					},
					&task.Task{
						Name:  "task-d",
						ID:    40,
						State: task.StateFinished,
					},
				}
				manager.TasksReturnsOnCall(0, tasks)
				manager.FindByNameStub = func(name string) *task.Task {
					for _, t := range tasks {
						if t.Name == name {
							return t
						}
					}
					return nil
				}
			})
			It("just prints out the task underneath their states in order", func() {
				Expect(r.Run([]string{"show"})).To(Succeed())
				expectedOutput := `RUNNING tasks:
  task-a \(10\)
BLOCKED tasks:
WAITING tasks:
  task-b \(20\)
  task-c \(30\)
FINISHED tasks:
  task-d \(40\)`
				Expect(stdoutWriter).To(gbytes.Say(expectedOutput))
			})

			Context("when a task name argument is passed", func() {
				It("prints the details about a task", func() {
					Expect(r.Run([]string{"show", "task-b"})).To(Succeed())
					expectedOutput := `Name: task-b
ID: 20
Created: \w+ \w+ \d\d? \d\d:\d\d
Priority: 3
State: WAITING`
					Expect(stdoutWriter).To(gbytes.Say(expectedOutput))
				})
			})
		})
	})

	Describe("note", func() {
		Context("when the task exists", func() {
			It("adds a note to the task", func() {
				Expect(r.Run([]string{"note", "task-a", "tuna"})).To(Succeed())

				name, note := manager.NoteArgsForCall(0)
				Expect(name).To(Equal("task-a"))
				Expect(note).To(Equal("tuna"))
			})
		})

		Context("when the manager fails to add a note", func() {
			BeforeEach(func() {
				manager.NoteReturnsOnCall(0, errors.New("task does not exist"))
			})

			It("displays the error to the user", func() {
				err := r.Run([]string{"note", "task-a", "tuna"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot add note: task does not exist"))
			})
		})
	})

	Describe("set-priority", func() {
		Context("when the task exists", func() {
			It("sets the priority on the task", func() {
				Expect(r.Run([]string{"set-priority", "task-a", "10"})).To(Succeed())

				name, prio := manager.SetPriorityArgsForCall(0)
				Expect(name).To(Equal("task-a"))
				Expect(prio).To(Equal(10))
			})
		})

		Context("when the manager fails to set the priority", func() {
			BeforeEach(func() {
				manager.SetPriorityReturnsOnCall(0, errors.New("task does not exist"))
			})

			It("displays the error to the user", func() {
				err := r.Run([]string{"set-priority", "task-a", "10"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot set priority: task does not exist"))
			})
		})

		Context("when the second argument is not a number", func() {
			It("displays the error to the user", func() {
				err := r.Run([]string{"set-priority", "task-a", "tuna"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot set priority: invalid priority: 'tuna'"))
			})
		})
	})

	Describe("set-<state>", func() {
		BeforeEach(func() {
			manager.FindByNameReturns(&task.Task{Name: "task-a"})
		})

		It("sets the state correctly for all valid states", func() {
			states := []task.State{
				task.StateRunning,
				task.StateBlocked,
				task.StateWaiting,
				task.StateFinished,
			}
			stateStrings := []string{"running", "blocked", "waiting", "finished"}
			for i, state := range states {
				stateString := stateStrings[i]
				cmd := fmt.Sprintf("set-%s", stateString)
				Expect(r.Run([]string{cmd, "task-a"})).To(Succeed(), "Command = %s", cmd)
				name, stateArg := manager.SetStateArgsForCall(i)
				Expect(name).To(Equal("task-a"))
				Expect(stateArg).To(Equal(state), "Command = %s", cmd)
			}
		})

		Context("when manager.SetState returns an error", func() {
			BeforeEach(func() {
				manager.SetStateReturns(errors.New("failed to set state"))
			})

			It("prints the error to the user", func() {
				err := r.Run([]string{"set-running", "task-a"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot set state: failed to set state"))
			})
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				manager.FindByNameReturns(nil)
			})

			It("prints the error to the user", func() {
				err := r.Run([]string{"set-running", "task-a"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unknown task: task-a"))
			})
		})
	})

	Describe("journal", func() {
		BeforeEach(func() {
			manager.EventsReturnsOnCall(0, []*task.Event{
				&task.Event{
					TaskID: 1,
					Title:  "event-a",
				},
				&task.Event{
					TaskID: 5,
					Title:  "event-b",
				},
				&task.Event{
					TaskID: 1,
					Title:  "event-c",
				},
				&task.Event{
					TaskID: 5,
					Title:  "event-d",
				},
			})
		})

		Context("when no task name is passed", func() {
			It("prints all journal entries in correct order", func() {
				Expect(r.Run([]string{"journal"})).To(Succeed())
				Expect(stdoutWriter).To(gbytes.Say("\\[.*\\]: event-d"))
				Expect(stdoutWriter).To(gbytes.Say("\\[.*\\]: event-c"))
				Expect(stdoutWriter).To(gbytes.Say("\\[.*\\]: event-b"))
				Expect(stdoutWriter).To(gbytes.Say("\\[.*\\]: event-a"))
			})
		})

		Context("when task name is passed", func() {

			Context("when the task does exist", func() {
				BeforeEach(func() {
					manager.FindByNameReturnsOnCall(0, &task.Task{
						ID: 1,
					})
				})

				It("prints the journal entries associated with that task", func() {
					Expect(r.Run([]string{"journal", "task-a"})).To(Succeed())
					Expect(stdoutWriter).To(gbytes.Say("\\[.*\\]: event-c"))
					Expect(stdoutWriter).To(gbytes.Say("\\[.*\\]: event-a"))
				})
			})

			Context("when the task does not exist", func() {
				It("prints an error", func() {
					err := r.Run([]string{"journal", "not-a-real-task"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("unknown task: not-a-real-task"))
				})
			})
		})
	})
})
