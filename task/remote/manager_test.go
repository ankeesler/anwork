// TODO: run the generic manager tests against the API!
package remote_test

import (
	"errors"
	"time"

	"github.com/ankeesler/anwork/task"
	"github.com/ankeesler/anwork/task/remote"
	"github.com/ankeesler/anwork/task/remote/remotefakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Manager", func() {
	var (
		client  *remotefakes.FakeAPIClient
		manager task.Manager
	)

	BeforeEach(func() {
		client = &remotefakes.FakeAPIClient{}

		var err error
		manager, err = remote.NewManagerFactory(client).Create()
		Expect(err).NotTo(HaveOccurred())
	})

	Describe("Create", func() {
		It("successfully can create tasks via its client", func() {
			Expect(manager.Create("a")).To(Succeed())
			Expect(client.CreateTaskCallCount()).To(Equal(1))
			Expect(client.CreateTaskArgsForCall(0)).To(Equal("a"))
		})

		Context("when the request returns a failure", func() {
			BeforeEach(func() {
				client.CreateTaskReturnsOnCall(0, errors.New("failed to create task"))
			})

			It("prints the failure message", func() {
				Expect(manager.Create("a")).To(MatchError("failed to create task"))
			})
		})
	})

	Describe("Delete", func() {
		BeforeEach(func() {
			tasks := []*task.Task{
				&task.Task{Name: "a", ID: 1},
				&task.Task{Name: "b", ID: 2},
				&task.Task{Name: "c", ID: 3},
			}
			client.GetTasksReturnsOnCall(0, tasks, nil)
		})

		It("successfully can delete tasks via a DELETE to /api/v1/tasks/:id", func() {
			Expect(manager.Delete("a")).To(Succeed())

			Expect(client.GetTasksCallCount()).To(Equal(1))

			Expect(client.DeleteTaskCallCount()).To(Equal(1))
			Expect(client.DeleteTaskArgsForCall(0)).To(Equal(1))
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				tasks := []*task.Task{
					&task.Task{Name: "b", ID: 2},
					&task.Task{Name: "c", ID: 3},
				}
				client.GetTasksReturnsOnCall(0, tasks, nil)
			})

			It("returns a failure", func() {
				Expect(manager.Delete("a")).To(MatchError("Unknown task with name a"))
			})
		})

		Context("when the request returns a failure", func() {
			BeforeEach(func() {
				client.DeleteTaskReturnsOnCall(0, errors.New("failed to delete task"))
			})

			It("prints the failure message", func() {
				Expect(manager.Delete("a")).To(MatchError("failed to delete task"))
			})
		})
	})

	Describe("Tasks", func() {
		var tasks []*task.Task

		BeforeEach(func() {
			tasks = []*task.Task{
				&task.Task{Name: "task-a", ID: 1},
				&task.Task{Name: "task-b", ID: 2},
				&task.Task{Name: "task-c", ID: 3},
			}
			client.GetTasksReturnsOnCall(0, tasks, nil)
		})

		It("returns the tasks via a call to the client", func() {
			actualTasks := manager.Tasks()
			Expect(actualTasks).To(Equal(tasks))

			Expect(client.GetTasksCallCount()).To(Equal(1))
		})

		Context("when the client fails", func() {
			BeforeEach(func() {
				client.GetTasksReturnsOnCall(0, nil, errors.New("failed to get tasks"))
			})

			It("panics...yeesh", func() {
				Expect(func() { manager.Tasks() }).To(Panic())
			})
		})
	})

	Describe("FindByName", func() {
		var expectedTask *task.Task

		BeforeEach(func() {
			tasks := []*task.Task{
				&task.Task{Name: "task-a", ID: 1},
				&task.Task{Name: "task-b", ID: 2},
				&task.Task{Name: "task-c", ID: 3},
			}
			expectedTask = tasks[0]
			client.GetTasksReturnsOnCall(0, tasks, nil)
		})

		It("returns the task via a call to the client", func() {
			actualTask := manager.FindByName("task-a")
			Expect(actualTask).To(Equal(expectedTask))

			Expect(client.GetTasksCallCount()).To(Equal(1))
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				tasks := []*task.Task{
					&task.Task{Name: "task-b", ID: 2},
					&task.Task{Name: "task-c", ID: 3},
				}
				client.GetTasksReturnsOnCall(0, tasks, nil)
			})

			It("returns nil", func() {
				actualTask := manager.FindByName("task-a")
				Expect(actualTask).To(BeNil())
			})
		})

		Context("when the client fails", func() {
			BeforeEach(func() {
				client.GetTasksReturnsOnCall(0, nil, errors.New("failed to get tasks"))
			})

			It("panics...yeesh", func() {
				Expect(func() { manager.Tasks() }).To(Panic())
			})
		})
	})

	Describe("FindByID", func() {
		var expectedTask *task.Task

		BeforeEach(func() {
			expectedTask = &task.Task{Name: "task-a", ID: 1}
			client.GetTaskReturnsOnCall(0, expectedTask, nil)
		})

		It("returns the task via a call to the client", func() {
			actualTask := manager.FindByID(1)
			Expect(actualTask).To(Equal(expectedTask))

			Expect(client.GetTaskCallCount()).To(Equal(1))
			Expect(client.GetTaskArgsForCall(0)).To(Equal(1))
		})

		Context("when the client fails", func() {
			BeforeEach(func() {
				client.GetTasksReturnsOnCall(0, nil, errors.New("failed to get tasks"))
			})

			It("panics...yeesh", func() {
				Expect(func() { manager.Tasks() }).To(Panic())
			})
		})
	})

	Describe("SetPriority", func() {
		BeforeEach(func() {
			tasks := []*task.Task{
				&task.Task{Name: "task-a", ID: 1},
				&task.Task{Name: "task-b", ID: 2},
				&task.Task{Name: "task-c", ID: 3},
			}
			client.GetTasksReturnsOnCall(0, tasks, nil)
		})

		It("updates the task via a call to the client", func() {
			Expect(manager.SetPriority("task-a", 10)).To(Succeed())

			Expect(client.GetTasksCallCount()).To(Equal(1))

			Expect(client.UpdatePriorityCallCount()).To(Equal(1))
			id, prio := client.UpdatePriorityArgsForCall(0)
			Expect(id).To(Equal(1))
			Expect(prio).To(Equal(10))
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				tasks := []*task.Task{
					&task.Task{Name: "task-b", ID: 2},
					&task.Task{Name: "task-c", ID: 3},
				}
				client.GetTasksReturnsOnCall(0, tasks, nil)
			})

			It("returns an error", func() {
				err := manager.SetPriority("task-a", 10)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unknown task with name task-a"))
			})
		})

		Context("when the client fails", func() {
			BeforeEach(func() {
				client.UpdatePriorityReturnsOnCall(0, errors.New("failed to update prio"))
			})

			It("returns the error", func() {
				Expect(manager.SetPriority("task-a", 10)).To(MatchError("failed to update prio"))
			})
		})
	})

	Describe("SetState", func() {
		BeforeEach(func() {
			tasks := []*task.Task{
				&task.Task{Name: "task-a", ID: 1},
				&task.Task{Name: "task-b", ID: 2},
				&task.Task{Name: "task-c", ID: 3},
			}
			client.GetTasksReturnsOnCall(0, tasks, nil)
		})

		It("updates the task via a call to the client", func() {
			Expect(manager.SetState("task-a", task.StateRunning)).To(Succeed())

			Expect(client.GetTasksCallCount()).To(Equal(1))

			Expect(client.UpdateStateCallCount()).To(Equal(1))
			id, state := client.UpdateStateArgsForCall(0)
			Expect(id).To(Equal(1))
			Expect(state).To(Equal(task.StateRunning))
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				tasks := []*task.Task{
					&task.Task{Name: "task-b", ID: 2},
					&task.Task{Name: "task-c", ID: 3},
				}
				client.GetTasksReturnsOnCall(0, tasks, nil)
			})

			It("returns an error", func() {
				err := manager.SetState("task-a", task.StateRunning)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unknown task with name task-a"))
			})
		})

		Context("when the client fails", func() {
			BeforeEach(func() {
				client.UpdateStateReturnsOnCall(0, errors.New("failed to update state"))
			})

			It("returns the error", func() {
				Expect(manager.SetState("task-a", task.StateRunning)).To(MatchError("failed to update state"))
			})
		})
	})

	Describe("Note", func() {
		BeforeEach(func() {
			tasks := []*task.Task{
				&task.Task{Name: "task-a", ID: 1},
				&task.Task{Name: "task-b", ID: 2},
				&task.Task{Name: "task-c", ID: 3},
			}
			client.GetTasksReturnsOnCall(0, tasks, nil)
		})

		It("adds a note via a call to its client", func() {
			Expect(manager.Note("task-a", "here is a note")).To(Succeed())

			Expect(client.GetTasksCallCount()).To(Equal(1))
			Expect(client.CreateEventCallCount()).To(Equal(1))
			title, teyep, startTime, id := client.CreateEventArgsForCall(0)
			Expect(title).To(Equal("here is a note"))
			Expect(teyep).To(Equal(task.EventTypeNote))
			Expect(time.Unix(startTime, 0)).To(BeTemporally("~", time.Now(), time.Second))
			Expect(id).To(Equal(1))
		})

		Context("when the task does not exist", func() {
			BeforeEach(func() {
				tasks := []*task.Task{
					&task.Task{Name: "task-b", ID: 2},
					&task.Task{Name: "task-c", ID: 3},
				}
				client.GetTasksReturnsOnCall(0, tasks, nil)
			})

			It("returns a helpful error", func() {
				err := manager.Note("task-a", "here is a note")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Unknown task with name task-a"))
			})
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				client.GetTasksReturnsOnCall(0, nil, errors.New("failed to get tasks"))
			})

			It("panics...yeesh", func() {
				Expect(func() { manager.Note("task-a", "here is a note") }).To(Panic())
			})
		})
	})

	Describe("Events", func() {
		var events []*task.Event

		BeforeEach(func() {
			events = []*task.Event{
				&task.Event{Title: "event-a", TaskID: 1},
				&task.Event{Title: "event-b", TaskID: 2},
				&task.Event{Title: "event-c", TaskID: 3},
			}
			client.GetEventsReturnsOnCall(0, events, nil)
		})

		It("gets the events via a call to the client", func() {
			actualEvents := manager.Events()
			Expect(actualEvents).To(Equal(events))

			Expect(client.GetEventsCallCount()).To(Equal(1))
		})

		Context("when the request fails", func() {
			BeforeEach(func() {
				client.GetEventsReturnsOnCall(0, nil, errors.New("failed request"))
			})

			It("...panics, I guess?", func() {
				Expect(func() { manager.Events() }).To(Panic())
			})
		})
	})

	Describe("DeleteEvent", func() {
		It("deletes an event via a call to the client", func() {
			Expect(manager.DeleteEvent(12345)).To(Succeed())

			Expect(client.DeleteEventCallCount()).To(Equal(1))
			Expect(client.DeleteEventArgsForCall(0)).To(Equal(int64(12345)))
		})

		Context("when the client fails", func() {
			BeforeEach(func() {
				client.DeleteEventReturnsOnCall(0, errors.New("failed to delete event"))
			})
			It("returns the error", func() {
				Expect(manager.DeleteEvent(12345)).To(MatchError("failed to delete event"))
			})
		})
	})

	Describe("Reset", func() {
		var tasks []*task.Task
		var events []*task.Event
		BeforeEach(func() {
			tasks = []*task.Task{
				&task.Task{Name: "task-a", ID: 1},
				&task.Task{Name: "task-b", ID: 2},
				&task.Task{Name: "task-c", ID: 3},
			}
			client.GetTasksReturnsOnCall(0, tasks, nil)
			events = []*task.Event{
				&task.Event{Title: "event-a", Date: int64(1)},
				&task.Event{Title: "event-b", Date: int64(2)},
				&task.Event{Title: "event-c", Date: int64(3)},
			}
			client.GetEventsReturnsOnCall(0, events, nil)
		})

		It("DELETE's all of the tasks and events", func() {
			Expect(manager.Reset()).To(Succeed())

			Expect(client.GetTasksCallCount()).To(Equal(1))
			Expect(client.GetEventsCallCount()).To(Equal(1))

			Expect(client.DeleteTaskCallCount()).To(Equal(3))
			Expect(client.DeleteTaskArgsForCall(0)).To(Equal(1))
			Expect(client.DeleteTaskArgsForCall(1)).To(Equal(2))
			Expect(client.DeleteTaskArgsForCall(2)).To(Equal(3))

			Expect(client.DeleteEventCallCount()).To(Equal(3))
			Expect(client.DeleteEventArgsForCall(0)).To(Equal(int64(1)))
			Expect(client.DeleteEventArgsForCall(1)).To(Equal(int64(2)))
			Expect(client.DeleteEventArgsForCall(2)).To(Equal(int64(3)))
		})

		Context("when we fail to get the tasks", func() {
			BeforeEach(func() {
				client.GetTasksReturnsOnCall(0, nil, errors.New("failed to get tasks"))
			})

			It("returns the failure message", func() {
				Expect(manager.Reset()).To(MatchError("failed to get tasks"))
			})
		})

		Context("when we fail to get the events", func() {
			BeforeEach(func() {
				client.GetEventsReturnsOnCall(0, nil, errors.New("failed to get events"))
			})

			It("returns the failure message", func() {
				Expect(manager.Reset()).To(MatchError("failed to get events"))
			})
		})

		Context("when some of the deletes return failure", func() {
			BeforeEach(func() {
				client.DeleteTaskReturnsOnCall(0, errors.New("failed to delete task"))
				client.DeleteEventReturnsOnCall(2, errors.New("failed to delete event"))
			})

			It("returns a formatted error message", func() {
				err := manager.Reset()
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("Encountered errors during reset:\n"))
				Expect(err.Error()).To(ContainSubstring("  delete task 1: failed to delete task"))
				Expect(err.Error()).To(ContainSubstring("  delete event 3: failed to delete event"))
			})
		})
	})

	Describe("Rename", func() {
		BeforeEach(func() {
			tasks := []*task.Task{
				&task.Task{Name: "task-a", ID: 1},
				&task.Task{Name: "task-b", ID: 2},
				&task.Task{Name: "task-c", ID: 3},
			}
			client.GetTasksReturnsOnCall(0, tasks, nil)
		})

		It("renames a task via a call to the client", func() {
			Expect(manager.Rename("task-a", "task-b")).To(Succeed())

			Expect(client.GetTasksCallCount()).To(Equal(1))

			Expect(client.UpdateNameCallCount()).To(Equal(1))
			id, to := client.UpdateNameArgsForCall(0)
			Expect(id).To(Equal(1))
			Expect(to).To(Equal("task-b"))
		})

		It("returns an error for an unknown task", func() {
			err := manager.Rename("task-z", "task-b")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(ContainSubstring("Unknown task with name task-z"))

			Expect(client.GetTasksCallCount()).To(Equal(1))
		})

		Context("when the client fails", func() {
			BeforeEach(func() {
				client.UpdateNameReturnsOnCall(0, errors.New("failed to rename task"))
			})
			It("returns the error", func() {
				err := manager.Rename("task-a", "task-b")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("failed to rename task"))
			})
		})
	})
})
