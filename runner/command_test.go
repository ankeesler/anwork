package runner_test

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/ankeesler/anwork/manager/managerfakes"
	"github.com/ankeesler/anwork/runner"
	"github.com/ankeesler/anwork/task"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Command", func() {
	var (
		r                         *runner.Runner
		stdoutWriter, debugWriter *gbytes.Buffer
		manager                   *managerfakes.FakeManager
	)

	BeforeEach(func() {
		manager = &managerfakes.FakeManager{}

		stdoutWriter = gbytes.NewBuffer()
		debugWriter = gbytes.NewBuffer()

		bi := &runner.BuildInfo{
			Hash: "abc123",
			Date: "February 22, 1992",
		}
		r = runner.New(bi, manager, stdoutWriter, debugWriter)
	})

	Describe("version", func() {
		It("prints out the version, git hash, and date", func() {
			Expect(r.Run([]string{"version"})).To(Succeed())
			Eventually(stdoutWriter).Should(gbytes.Say(fmt.Sprintf("ANWORK Version = %d\n", runner.Version)))
			Eventually(stdoutWriter).Should(gbytes.Say(fmt.Sprintf("ANWORK Build Hash = abc123\n")))
			Eventually(stdoutWriter).Should(gbytes.Say(fmt.Sprintf("ANWORK Build Date = February 22, 1992")))
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
				Expect(r.Run([]string{"reset"})).To(Succeed())
				Eventually(stdoutWriter).Should(gbytes.Say("Are you sure you want to delete all data \\[y/n\\]: "))
				Eventually(stdoutWriter).Should(gbytes.Say("OK, deleting all data"))
			})

			It("tells the manager to reset", func() {
				Expect(r.Run([]string{"reset"})).To(Succeed())
				Expect(manager.ResetCallCount()).To(Equal(1))
			})

			Context("when the manager fails to reset", func() {
				BeforeEach(func() {
					manager.ResetReturnsOnCall(0, errors.New("some reset error"))
				})

				It("returns the error", func() {
					err := r.Run([]string{"reset"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("some reset error"))
				})
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
				Expect(manager.ResetCallCount()).To(Equal(0))
			})
		})
	})

	Describe("summary", func() {
		BeforeEach(func() {
			twoDaysAgo := time.Now().Add(-1 * (time.Hour * 24 * 2))
			tenDaysAgo := time.Now().Add(-1 * (time.Hour * 24 * 10))
			manager.EventsReturns([]*task.Event{
				&task.Event{
					Type:   task.EventTypeSetState,
					Title:  "foo",
					TaskID: 1,
				},
				&task.Event{
					Type:   task.EventTypeCreate,
					Title:  "task-a created",
					Date:   twoDaysAgo.Unix(),
					TaskID: 5,
				},
				&task.Event{
					Type:   task.EventTypeSetState,
					Title:  "task-a changed to Finished",
					Date:   twoDaysAgo.Unix(),
					TaskID: 5,
				},
				&task.Event{
					Type:   task.EventTypeSetState,
					Title:  "task-b changed to Finished",
					Date:   tenDaysAgo.Unix(),
					TaskID: 10,
				},
			}, nil)
		})

		It("shows the tasks that have been completed in the provided number of days", func() {
			Expect(r.Run([]string{"summary", "5"})).To(Succeed())

			Eventually(stdoutWriter).Should(gbytes.Say("\\[.*\\]: task-a changed to Finished"))
			Eventually(stdoutWriter).Should(gbytes.Say("  took \\d+\\w"))
		})

		It("does not show the tasks that haven't been completed in the provided number of days", func() {
			Expect(r.Run([]string{"summary", "5"})).To(Succeed())

			Eventually(stdoutWriter).ShouldNot(gbytes.Say("foo"))
			Eventually(stdoutWriter).ShouldNot(gbytes.Say("task-b changed to Finished"))
		})

		Context("when the number of days is invalid", func() {
			It("doesn't display anything", func() {
				err := r.Run([]string{"summary", "tuna"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Cannot convert days tuna to number"))
			})
		})

		Context("when the number of days is 0", func() {
			It("doesn't display anything", func() {
				Expect(r.Run([]string{"summary", "0"})).To(Succeed())
				Expect(stdoutWriter.Contents()).To(BeEmpty())
			})
		})

		Context("when the number of days is negative", func() {
			It("doesn't display anything", func() {
				Expect(r.Run([]string{"summary", "-3"})).To(Succeed())
				Expect(stdoutWriter.Contents()).To(BeEmpty())
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
				manager.DeleteReturnsOnCall(0, nil)
				manager.FindByNameReturnsOnCall(0, &task.Task{Name: "task-a"}, nil)
			})

			It("successfully deletes existing tasks", func() {
				Expect(r.Run([]string{"delete", "task-a"})).To(Succeed())
				Expect(manager.DeleteCallCount()).To(Equal(1))
				Expect(manager.DeleteArgsForCall(0)).To(Equal("task-a"))
			})
		})

		Context("when the manager fails to find the task", func() {
			BeforeEach(func() {
				manager.DeleteReturnsOnCall(0, nil)
				manager.FindByNameReturnsOnCall(0, nil, errors.New("some find error"))
			})

			It("returns the error", func() {
				err := r.Run([]string{"delete", "task-a"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("some find error"))
			})
		})

		Context("when a task spec is passed", func() {
			BeforeEach(func() {
				manager.DeleteReturnsOnCall(0, nil)
				manager.FindByIDReturnsOnCall(0, &task.Task{Name: "task-a"}, nil)
			})

			It("parses the task spec and deletes the correct task", func() {
				Expect(r.Run([]string{"delete", "@1"})).To(Succeed())
				Expect(manager.FindByIDArgsForCall(0)).To(Equal(1))
				Expect(manager.DeleteArgsForCall(0)).To(Equal("task-a"))
			})

			Context("when the task spec is totally bogus", func() {
				It("returns a helpful error", func() {
					err := r.Run([]string{"delete", "@tuna"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("cannot parse task ID"))
				})
			})

			Context("when the task spec is not a valid task ID", func() {
				BeforeEach(func() {
					manager.FindByIDReturnsOnCall(0, nil, nil)
				})

				It("returns a helpful error", func() {
					err := r.Run([]string{"delete", "@1"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("unknown task ID in task spec: 1"))
				})
			})
		})

		Context("when the manager fails to delete the task", func() {
			BeforeEach(func() {
				manager.DeleteReturnsOnCall(0, errors.New("some delete error"))
				manager.FindByNameReturnsOnCall(0, &task.Task{Name: "task-a"}, nil)
			})

			It("returns the error that manager returned", func() {
				err := r.Run([]string{"delete", "task-a"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("some delete error"))
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
				manager.TasksReturnsOnCall(0, tasks, nil)
			})

			Context("when the manager successfully deletes each task", func() {
				BeforeEach(func() {
					manager.DeleteReturnsOnCall(0, nil)
					manager.DeleteReturnsOnCall(1, nil)
					manager.DeleteReturnsOnCall(2, nil)
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
					manager.DeleteReturnsOnCall(0, errors.New("delete failure 0"))
					manager.DeleteReturnsOnCall(1, nil)
					manager.DeleteReturnsOnCall(2, errors.New("delete failure 2"))
				})

				It("notifies the user of which task was not able to be deleted", func() {
					err := r.Run([]string{"delete-all"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("unable to delete task task-a: delete failure 0"))
					Expect(err.Error()).To(ContainSubstring("unable to delete task task-c: delete failure 2"))

					Expect(manager.DeleteCallCount()).To(Equal(3))
					Expect(manager.DeleteArgsForCall(0)).To(Equal("task-a"))
					Expect(manager.DeleteArgsForCall(1)).To(Equal("task-b"))
					Expect(manager.DeleteArgsForCall(2)).To(Equal("task-c"))
				})
			})
		})

		Context("when the manager fails to get the tasks", func() {
			BeforeEach(func() {
				manager.TasksReturnsOnCall(0, nil, errors.New("some error"))
			})

			It("returns the error", func() {
				err := r.Run([]string{"delete-all"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("some error"))
			})
		})
	})

	Describe("show", func() {
		Context("when there are no tasks", func() {
			It("just prints out the task states", func() {
				Expect(r.Run([]string{"show"})).To(Succeed())
				Expect(stdoutWriter).To(gbytes.Say("RUNNING tasks:\nBLOCKED tasks:\nREADY tasks:\nFINISHED tasks:\n"))
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
						State:    task.StateReady,
					},
					&task.Task{
						Name:  "task-c",
						ID:    30,
						State: task.StateReady,
					},
					&task.Task{
						Name:  "task-d",
						ID:    40,
						State: task.StateFinished,
					},
				}
				manager.TasksReturnsOnCall(0, tasks, nil)
				manager.FindByNameStub = func(name string) (*task.Task, error) {
					for _, t := range tasks {
						if t.Name == name {
							return t, nil
						}
					}
					return nil, nil
				}
			})
			It("just prints out the task underneath their states in order", func() {
				Expect(r.Run([]string{"show"})).To(Succeed())
				expectedOutput := `RUNNING tasks:
  task-a \(10\)
BLOCKED tasks:
READY tasks:
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
State: READY`
					Expect(stdoutWriter).To(gbytes.Say(expectedOutput))
				})
			})
		})

		Context("when a task spec is passed", func() {
			BeforeEach(func() {
				manager.FindByIDReturnsOnCall(0,
					&task.Task{
						Name:     "task-a",
						ID:       10,
						State:    task.StateReady,
						Priority: 3,
					},
					nil,
				)
			})

			It("parses the task spec and prints out information about the task", func() {
				Expect(r.Run([]string{"show", "@1"})).To(Succeed())
				Expect(manager.FindByIDArgsForCall(0)).To(Equal(1))
				expectedOutput := `Name: task-a
ID: 10
Created: \w+ \w+ \d\d? \d\d:\d\d
Priority: 3
State: READY`
				Expect(stdoutWriter).To(gbytes.Say(expectedOutput))
			})

			Context("when the task spec is totally bogus", func() {
				It("returns a helpful error", func() {
					err := r.Run([]string{"show", "@tuna"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("cannot parse task ID"))
				})
			})

			Context("when the task spec is not a valid task ID", func() {
				BeforeEach(func() {
					manager.FindByIDReturnsOnCall(0, nil, nil)
				})

				It("returns a helpful error", func() {
					err := r.Run([]string{"show", "@1"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("unknown task ID in task spec: 1"))
				})
			})

			Context("when we fail to get the task", func() {
				BeforeEach(func() {
					manager.FindByIDReturnsOnCall(0, nil, errors.New("some find task"))
				})

				It("returns a helpful error", func() {
					err := r.Run([]string{"show", "@1"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("some find task"))
				})
			})
		})

		Context("when tasks fails", func() {
			BeforeEach(func() {
				manager.TasksReturnsOnCall(0, nil, errors.New("some error"))
			})

			It("it returns the error", func() {
				err := r.Run([]string{"show"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("some error"))
			})
		})
	})

	Describe("note", func() {
		Context("when the task exists", func() {
			BeforeEach(func() {
				manager.FindByNameReturnsOnCall(0, &task.Task{Name: "task-a"}, nil)
			})

			It("adds a note to the task", func() {
				Expect(r.Run([]string{"note", "task-a", "tuna"})).To(Succeed())

				name, note := manager.NoteArgsForCall(0)
				Expect(name).To(Equal("task-a"))
				Expect(note).To(Equal("tuna"))
			})
		})

		Context("when we fail to find the task", func() {
			BeforeEach(func() {
				manager.FindByNameReturnsOnCall(0, nil, errors.New("some find error"))
			})

			It("returns the error", func() {
				err := r.Run([]string{"note", "task-a", "tuna"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("some find error"))
			})
		})

		Context("when the manager fails to add a note", func() {
			BeforeEach(func() {
				manager.FindByNameReturnsOnCall(0, &task.Task{Name: "task-a"}, nil)
				manager.NoteReturnsOnCall(0, errors.New("task does not exist"))
			})

			It("displays the error to the user", func() {
				err := r.Run([]string{"note", "task-a", "tuna"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot add note: task does not exist"))
			})
		})

		Context("when a task spec is passed", func() {
			BeforeEach(func() {
				manager.FindByIDReturnsOnCall(0, &task.Task{Name: "task-a"}, nil)
			})

			It("parses the task spec and adds a note to the correct task", func() {
				Expect(r.Run([]string{"note", "@1", "tuna"})).To(Succeed())

				Expect(manager.FindByIDArgsForCall(0)).To(Equal(1))

				name, note := manager.NoteArgsForCall(0)
				Expect(name).To(Equal("task-a"))
				Expect(note).To(Equal("tuna"))
			})

			Context("when the task spec is totally bogus", func() {
				It("returns a helpful error", func() {
					err := r.Run([]string{"note", "@tuna", "fish"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("cannot parse task ID"))
				})
			})

			Context("when the task spec is not a valid task ID", func() {
				BeforeEach(func() {
					manager.FindByIDReturnsOnCall(0, nil, nil)
				})

				It("returns a helpful error", func() {
					err := r.Run([]string{"note", "@1", "tuna"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("unknown task ID in task spec: 1"))
				})
			})
		})
	})

	Describe("set-priority", func() {
		Context("when the task exists", func() {
			BeforeEach(func() {
				manager.FindByNameReturnsOnCall(0, &task.Task{Name: "task-a"}, nil)
			})

			It("sets the priority on the task", func() {
				Expect(r.Run([]string{"set-priority", "task-a", "10"})).To(Succeed())

				name, prio := manager.SetPriorityArgsForCall(0)
				Expect(name).To(Equal("task-a"))
				Expect(prio).To(Equal(10))
			})
		})

		Context("when we fail to find the task", func() {
			BeforeEach(func() {
				manager.FindByNameReturnsOnCall(0, nil, errors.New("some find error"))
			})

			It("returns the error", func() {
				err := r.Run([]string{"set-priority", "task-a", "10"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("some find error"))
			})
		})

		Context("when the manager fails to set the priority", func() {
			BeforeEach(func() {
				manager.FindByNameReturnsOnCall(0, &task.Task{Name: "task-a"}, nil)
				manager.SetPriorityReturnsOnCall(0, errors.New("task does not exist"))
			})

			It("displays the error to the user", func() {
				err := r.Run([]string{"set-priority", "task-a", "10"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot set priority: task does not exist"))
			})
		})

		Context("when the second argument is not a number", func() {
			BeforeEach(func() {
				manager.FindByNameReturnsOnCall(0, &task.Task{Name: "task-a"}, nil)
			})

			It("displays the error to the user", func() {
				err := r.Run([]string{"set-priority", "task-a", "tuna"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("cannot set priority: invalid priority: 'tuna'"))
			})
		})

		Context("when a task spec is passed", func() {
			BeforeEach(func() {
				manager.FindByIDReturnsOnCall(0, &task.Task{Name: "task-a"}, nil)
			})

			It("parses the task spec and sets the priority on the task", func() {
				Expect(r.Run([]string{"set-priority", "@1", "5"})).To(Succeed())

				Expect(manager.FindByIDArgsForCall(0)).To(Equal(1))

				name, prio := manager.SetPriorityArgsForCall(0)
				Expect(name).To(Equal("task-a"))
				Expect(prio).To(Equal(5))
			})

			Context("when the task spec is totally bogus", func() {
				It("returns a helpful error", func() {
					err := r.Run([]string{"set-priority", "@tuna", "5"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("cannot parse task ID"))
				})
			})

			Context("when the task spec is not a valid task ID", func() {
				BeforeEach(func() {
					manager.FindByIDReturnsOnCall(0, nil, nil)
				})

				It("returns a helpful error", func() {
					err := r.Run([]string{"set-priority", "@1", "5"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("unknown task ID in task spec: 1"))
				})
			})
		})
	})

	Describe("set-<state>", func() {
		BeforeEach(func() {
			manager.FindByNameReturns(&task.Task{Name: "task-a"}, nil)
		})

		It("sets the state correctly for all valid states", func() {
			states := []task.State{
				task.StateRunning,
				task.StateBlocked,
				task.StateReady,
				task.StateFinished,
			}
			stateStrings := []string{"running", "blocked", "ready", "finished"}
			for i, state := range states {
				stateString := stateStrings[i]
				cmd := fmt.Sprintf("set-%s", stateString)
				Expect(r.Run([]string{cmd, "task-a"})).To(Succeed(), "Command = %s", cmd)
				name, stateArg := manager.SetStateArgsForCall(i)
				Expect(name).To(Equal("task-a"))
				Expect(stateArg).To(Equal(state), "Command = %s", cmd)
			}
		})

		Context("when we fail to find the task", func() {
			BeforeEach(func() {
				manager.FindByNameReturnsOnCall(0, nil, errors.New("some find error"))
			})

			It("returns the error", func() {
				err := r.Run([]string{"set-running", "task-a"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("some find error"))
			})
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
				manager.FindByNameReturns(nil, nil)
			})

			It("prints the error to the user", func() {
				err := r.Run([]string{"set-running", "task-a"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unknown task: task-a"))
			})
		})

		Context("when a task spec is passed", func() {
			BeforeEach(func() {
				manager.FindByIDReturnsOnCall(0, &task.Task{Name: "task-a"}, nil)
			})

			It("parses the task spec and sets the state on the task", func() {
				Expect(r.Run([]string{"set-running", "@1"})).To(Succeed())

				Expect(manager.FindByIDArgsForCall(0)).To(Equal(1))

				name, state := manager.SetStateArgsForCall(0)
				Expect(name).To(Equal("task-a"))
				Expect(state).To(Equal(task.State(task.StateRunning)))
			})

			Context("when the task spec is totally bogus", func() {
				It("returns a helpful error", func() {
					err := r.Run([]string{"set-running", "@tuna"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("cannot parse task ID"))
				})
			})

			Context("when the task spec is not a valid task ID", func() {
				BeforeEach(func() {
					manager.FindByIDReturnsOnCall(0, nil, nil)
				})

				It("returns a helpful error", func() {
					err := r.Run([]string{"set-running", "@1"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("unknown task ID in task spec: 1"))
				})
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
			},
				nil,
			)
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

		Context("when events fails", func() {
			BeforeEach(func() {
				manager.EventsReturnsOnCall(0, nil, errors.New("some error"))
			})

			It("returns the error", func() {
				err := r.Run([]string{"journal"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("some error"))
			})
		})

		Context("when task name is passed", func() {

			Context("when the task does exist", func() {
				BeforeEach(func() {
					manager.FindByNameReturnsOnCall(0, &task.Task{
						ID: 1,
					},
						nil,
					)
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

			Context("when we fail to find the task", func() {
				BeforeEach(func() {
					manager.FindByNameReturnsOnCall(0, nil, errors.New("some find error"))
				})

				It("prints an error", func() {
					err := r.Run([]string{"journal", "not-a-real-task"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("some find error"))
				})
			})
		})

		Context("when a task spec is passed", func() {
			BeforeEach(func() {
				manager.FindByIDReturnsOnCall(0, &task.Task{Name: "task-a", ID: 1}, nil)
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
				},
					nil,
				)
			})

			It("parses the task spec and displays the journal", func() {
				Expect(r.Run([]string{"journal", "@1"})).To(Succeed())

				Expect(manager.FindByIDArgsForCall(0)).To(Equal(1))

				Expect(stdoutWriter).To(gbytes.Say("\\[.*\\]: event-c"))
			})

			Context("when the task spec is totally bogus", func() {
				It("returns a helpful error", func() {
					err := r.Run([]string{"journal", "@tuna"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("cannot parse task ID"))
				})
			})

			Context("when the task spec is not a valid task ID", func() {
				BeforeEach(func() {
					manager.FindByIDReturnsOnCall(0, nil, nil)
				})

				It("returns a helpful error", func() {
					err := r.Run([]string{"journal", "@1"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("unknown task ID in task spec: 1"))
				})
			})
		})
	})

	Describe("archive", func() {
		Context("when there are no tasks", func() {
			It("does not tell the manager to delete anything", func() {
				Expect(r.Run([]string{"archive"})).To(Succeed())
				Expect(manager.DeleteCallCount()).To(Equal(0))
			})
		})

		Context("when there are multiple tasks", func() {
			BeforeEach(func() {
				tasks := []*task.Task{
					&task.Task{Name: "task-a", State: task.StateFinished},
					&task.Task{Name: "task-b", State: task.StateBlocked},
					&task.Task{Name: "task-c", State: task.StateFinished},
				}
				manager.TasksReturnsOnCall(0, tasks, nil)
			})

			Context("when the manager deletes stuff happily", func() {
				BeforeEach(func() {
					manager.DeleteReturnsOnCall(0, nil)
					manager.DeleteReturnsOnCall(1, nil)
				})

				It("calls delete on each finished task", func() {
					Expect(r.Run([]string{"archive"})).To(Succeed())
					Expect(manager.DeleteCallCount()).To(Equal(2))
					Expect(manager.DeleteArgsForCall(0)).To(Equal("task-a"))
					Expect(manager.DeleteArgsForCall(1)).To(Equal("task-c"))
				})
			})

			Context("when the manager fails to delete a task", func() {
				BeforeEach(func() {
					manager.DeleteReturnsOnCall(0, errors.New("delete failure 0"))
					manager.DeleteReturnsOnCall(1, nil)
				})

				It("notifies the user of which task was not able to be deleted", func() {
					err := r.Run([]string{"archive"})
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(ContainSubstring("unable to delete task task-a: delete failure 0"))

					Expect(manager.DeleteCallCount()).To(Equal(2))
					Expect(manager.DeleteArgsForCall(0)).To(Equal("task-a"))
					Expect(manager.DeleteArgsForCall(1)).To(Equal("task-c"))
				})
			})
		})
	})

	Describe("rename", func() {
		Context("the manager succeeds", func() {
			BeforeEach(func() {
				manager.FindByNameReturnsOnCall(0, &task.Task{Name: "task-a", State: task.StateFinished}, nil)
				manager.RenameReturnsOnCall(0, nil)
			})

			It("calls the manager and adds a note and succeeds", func() {
				Expect(r.Run([]string{"rename", "task-a", "task-d"})).To(Succeed())

				from, to := manager.RenameArgsForCall(0)
				Expect(from).To(Equal("task-a"))
				Expect(to).To(Equal("task-d"))

				name, note := manager.NoteArgsForCall(0)
				Expect(name).To(Equal("task-d"))
				Expect(note).To(Equal("Renamed task 'task-a' to 'task-d'"))
			})
		})

		Context("when the manager fails", func() {
			BeforeEach(func() {
				manager.FindByNameReturnsOnCall(0, &task.Task{Name: "task-a", State: task.StateFinished}, nil)
				manager.RenameReturnsOnCall(0, errors.New("some rename error"))
			})

			It("calls delete on each finished task", func() {
				err := r.Run([]string{"rename", "task-a", "task-d"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unable to rename task task-a to task-d: some rename error"))
			})
		})

		Context("when we fail to find the task", func() {
			BeforeEach(func() {
				manager.FindByNameReturnsOnCall(0, nil, errors.New("some find error"))
			})

			It("returns an error", func() {
				err := r.Run([]string{"rename", "task-a", "task-d"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("some find error"))
			})
		})

		Context("when the from task is invalid", func() {
			BeforeEach(func() {
				manager.FindByNameReturnsOnCall(0, nil, nil)
			})

			It("calls delete on each finished task", func() {
				err := r.Run([]string{"rename", "task-z", "task-d"})
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("unknown task: task-z"))
			})
		})
	})
})
