package integration

import (
	"fmt"
	"io"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("anwork", func() {
	var (
		outBuf, errBuf *gbytes.Buffer
	)
	BeforeEach(func() {
		outBuf = gbytes.NewBuffer()
		errBuf = gbytes.NewBuffer()
	})
	AfterEach(func() {
		run(nil, nil, "reset", "y")
	})

	Context("when creating a task", func() {
		BeforeEach(func() {
			run(outBuf, errBuf, "create", "task-a")
		})
		It("shows the task as waiting", func() {
			run(outBuf, errBuf, "show")
			Expect(outBuf).To(gbytes.Say("WAITING tasks:\n  task-a \\(0\\)\nFINISHED tasks"))
		})
		It("shows the correct task details", func() {
			run(outBuf, errBuf, "show", "task-a")
			Expect(outBuf).To(gbytes.Say("Name: task-a\nID: 0\nCreated: .*\nPriority: 10\nState: WAITING"))
		})
		It("records the event in the task's journal", func() {
			run(outBuf, errBuf, "journal", "task-a")
			Expect(outBuf).To(gbytes.Say("\\[.*\\]: Created task task-a"))
		})
		It("records the event in the globas journal", func() {
			run(outBuf, errBuf, "journal")
			Expect(outBuf).To(gbytes.Say("\\[.*\\]: Created task task-a"))
		})
	})

	Context("when creating multiple tasks", func() {
		BeforeEach(func() {
			run(outBuf, errBuf, "create", "task-a")
			run(outBuf, errBuf, "create", "task-b")
			run(outBuf, errBuf, "create", "task-c")
		})
		It("shows the tasks as waiting in the order in which they were created", func() {
			run(outBuf, errBuf, "show")
			Expect(outBuf).To(gbytes.Say("WAITING tasks:\n  task-a \\(0\\)\n  task-b \\(1\\)\n  task-c \\(2\\)\nFINISHED tasks"))
		})
		It("shows the correct task details", func() {
			run(outBuf, errBuf, "show", "task-a")
			Expect(outBuf).To(gbytes.Say("Name: task-a\nID: 0\nCreated: .*\nPriority: 10\nState: WAITING"))
			run(outBuf, errBuf, "show", "task-b")
			Expect(outBuf).To(gbytes.Say("Name: task-b\nID: 1\nCreated: .*\nPriority: 10\nState: WAITING"))
			run(outBuf, errBuf, "show", "task-c")
			Expect(outBuf).To(gbytes.Say("Name: task-c\nID: 2\nCreated: .*\nPriority: 10\nState: WAITING"))
		})
		It("records the events in each task's journal", func() {
			run(outBuf, errBuf, "journal", "task-a")
			Expect(outBuf).To(gbytes.Say("\\[.*\\]: Created task task-a"))
			run(outBuf, errBuf, "journal", "task-b")
			Expect(outBuf).To(gbytes.Say("\\[.*\\]: Created task task-b"))
			run(outBuf, errBuf, "journal", "task-c")
			Expect(outBuf).To(gbytes.Say("\\[.*\\]: Created task task-c"))
		})
		It("records the events in the global journal in most recent to oldest order", func() {
			run(outBuf, errBuf, "journal")
			Expect(outBuf).To(gbytes.Say("\\[.*\\]: Created task task-c\n\\[.*\\]: Created task task-b\n\\[.*\\]: Created task task-a"))
		})
		Context("when changing the priority on tasks", func() {
			BeforeEach(func() {
				run(nil, nil, "set-priority", "task-a", "15")
				run(nil, nil, "set-priority", "task-b", "10")
				run(nil, nil, "set-priority", "task-c", "20")
			})
			It("properly shows the tasks in order of priority", func() {
				run(outBuf, errBuf, "show")
				Expect(outBuf).To(gbytes.Say("WAITING tasks:"))
				Expect(outBuf).To(gbytes.Say("  task-b"))
				Expect(outBuf).To(gbytes.Say("  task-a"))
				Expect(outBuf).To(gbytes.Say("  task-c"))

			})
			It("properly records the priorities", func() {
				run(outBuf, errBuf, "show", "task-a")
				Expect(outBuf).To(gbytes.Say("Name: task-a\nID: 0\nCreated: .*\nPriority: 15\nState: WAITING"))
				run(outBuf, errBuf, "show", "task-b")
				Expect(outBuf).To(gbytes.Say("Name: task-b\nID: 1\nCreated: .*\nPriority: 10\nState: WAITING"))
				run(outBuf, errBuf, "show", "task-c")
				Expect(outBuf).To(gbytes.Say("Name: task-c\nID: 2\nCreated: .*\nPriority: 20\nState: WAITING"))
			})
			It("records the events in each of the task's journals", func() {
				run(outBuf, errBuf, "journal", "task-a")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set priority on task task-a from 10 to 15"))
				run(outBuf, errBuf, "journal", "task-b")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set priority on task task-b from 10 to 10"))
				run(outBuf, errBuf, "journal", "task-c")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set priority on task task-c from 10 to 20"))
			})
			It("records the events in the global journal", func() {
				run(outBuf, errBuf, "journal")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set priority on task task-c from 10 to 20"))
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set priority on task task-b from 10 to 10"))
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set priority on task task-a from 10 to 15"))
			})
		})
		Context("when changing the state on tasks", func() {
			BeforeEach(func() {
				run(nil, nil, "set-running", "task-a")
				run(nil, nil, "set-finished", "task-b")
				run(nil, nil, "set-blocked", "task-c")
			})
			It("properly displays the states", func() {
				run(outBuf, errBuf, "show")
				Expect(outBuf).To(gbytes.Say("RUNNING tasks:\n  task-a"))
				Expect(outBuf).To(gbytes.Say("BLOCKED tasks:\n  task-c"))
				Expect(outBuf).To(gbytes.Say("FINISHED tasks:\n  task-b"))
			})
			It("properly records the states", func() {
				run(outBuf, errBuf, "show", "task-a")
				Expect(outBuf).To(gbytes.Say("Name: task-a\nID: 0\nCreated: .*\nPriority: 10\nState: RUNNING"))
				run(outBuf, errBuf, "show", "task-b")
				Expect(outBuf).To(gbytes.Say("Name: task-b\nID: 1\nCreated: .*\nPriority: 10\nState: FINISHED"))
				run(outBuf, errBuf, "show", "task-c")
				Expect(outBuf).To(gbytes.Say("Name: task-c\nID: 2\nCreated: .*\nPriority: 10\nState: BLOCKED"))
			})
			It("records the events in each of the task's journals", func() {
				run(outBuf, errBuf, "journal", "task-a")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set state on task task-a from Waiting to Running"))
				run(outBuf, errBuf, "journal", "task-b")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set state on task task-b from Waiting to Finished"))
				run(outBuf, errBuf, "journal", "task-c")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set state on task task-c from Waiting to Blocked"))
			})
			It("records the events in the global journal in order from newest to oldest", func() {
				run(outBuf, errBuf, "journal")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set state on task task-c from Waiting to Blocked"))
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set state on task task-b from Waiting to Finished"))
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set state on task task-a from Waiting to Running"))
			})
		})
		Context("when adding a note to tasks", func() {
			BeforeEach(func() {
				run(nil, nil, "note", "task-a", "Here is a note")
			})
			It("records the note in task-a's journal", func() {
				run(outBuf, errBuf, "journal", "task-a")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Note added to task task-a: Here is a note"))
			})
			It("records the note in the global journal", func() {
				run(outBuf, errBuf, "journal")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Note added to task task-a: Here is a note"))
			})
		})
		Context("when deleting a task", func() {
			BeforeEach(func() {
				run(nil, nil, "delete", "task-b")
			})
			It("properly displays the remaining tasks", func() {
				run(outBuf, errBuf, "show")
				Expect(outBuf).To(gbytes.Say("WAITING tasks:\n  task-a.*\n  task-c"))
			})
			It("fails when we try to show the deleted task", func() {
				runWithStatus(1, outBuf, errBuf, "show", "task-b")
			})
			It("fails to show the journal for the deleted task", func() {
				runWithStatus(1, outBuf, errBuf, "journal", "task-b")
			})
			It("records the events in the global journal in order from newest to oldest", func() {
				run(outBuf, errBuf, "journal")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Deleted task task-b"))
			})
		})
		Context("when deleting all tasks", func() {
			BeforeEach(func() {
				run(nil, nil, "delete-all")
			})
			It("no longer shows any tasks", func() {
				run(outBuf, errBuf, "show")
				Expect(outBuf).To(gbytes.Say("RUNNING tasks:\nBLOCKED tasks:\nWAITING tasks:\nFINISHED tasks"))
			})
			It("fails when we try to show the deleted tasks", func() {
				runWithStatus(1, outBuf, errBuf, "show", "task-a")
				runWithStatus(1, outBuf, errBuf, "show", "task-b")
				runWithStatus(1, outBuf, errBuf, "show", "task-c")
			})
			It("fails to show the journal for the deleted tasks", func() {
				runWithStatus(1, outBuf, errBuf, "journal", "task-a")
				runWithStatus(1, outBuf, errBuf, "journal", "task-b")
				runWithStatus(1, outBuf, errBuf, "journal", "task-c")
			})
			It("records the events in the global journal in order from newest to oldest", func() {
				run(outBuf, errBuf, "journal")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Deleted task task-c"))
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Deleted task task-b"))
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Deleted task task-a"))
			})
		})
		Context("when doing a reset", func() {
			var (
				s     *gexec.Session
				stdin io.Writer
			)
			BeforeEach(func() {
				var err error
				cmd := exec.Command(anworkBin, "reset")
				stdin, err = cmd.StdinPipe()
				Expect(err).ToNot(HaveOccurred())

				s, err = gexec.Start(cmd, outBuf, errBuf)
				Expect(err).ToNot(HaveOccurred())

				Eventually(outBuf).Should(gbytes.Say("Are you sure you want to delete all data \\[y/n\\]: "))
			})
			AfterEach(func() {
				s.Interrupt()
			})
			It("does not delete all data if the user says no", func() {
				stdin.Write([]byte("n\n"))
				Eventually(outBuf).Should(gbytes.Say("NOT deleting all data"))

				run(outBuf, errBuf, "show")
				Expect(outBuf).To(gbytes.Say("task-a"))
			})
			It("does delete all data if the user says yes", func() {
				stdin.Write([]byte("y\n"))
				Eventually(outBuf).Should(gbytes.Say("OK, deleting all data"))

				run(outBuf, errBuf, "show")
				Expect(outBuf).ToNot(gbytes.Say("task-a"))
			})
		})
		Context("when performing a summary", func() {
			BeforeEach(func() {
				run(nil, nil, "set-finished", "task-b")
			})
			It("reports the finished tasks", func() {
				run(outBuf, errBuf, "summary", "1")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set state on task task-b from Waiting to Finished"))
				Expect(outBuf).To(gbytes.Say("  took "))
			})
		})
	})

	Context("when creating a task, deleting it, and creating it again", func() {
		BeforeEach(func() {
			run(nil, nil, "create", "task-a")
			run(outBuf, errBuf, "show", "task-a")
			Expect(outBuf).To(gbytes.Say("ID: 0"))

			run(nil, nil, "delete", "task-a")
			runWithStatus(1, outBuf, errBuf, "show", "task-a")

			run(nil, nil, "create", "task-a")
		})
		It("should give unique IDs to both tasks", func() {
			run(outBuf, errBuf, "show", "task-a")
			Expect(outBuf).ToNot(gbytes.Say("ID: 0"))
		})
	})

	Measure("creating 10 tasks", func(b Benchmarker) {
		runtime := b.Time("runtime", func() {
			for i := 0; i < 10; i++ {
				run(nil, nil, "create", fmt.Sprintf("task-%d", i))
			}
		})
		Expect(runtime.Seconds()).To(BeNumerically("<", 1))
	}, 5)

	Measure("CRUD'ing 10 tasks", func(b Benchmarker) {
		runtime := b.Time("runtime", func() {
			for i := 0; i < 10; i++ {
				name := fmt.Sprintf("task-%d", i)
				run(nil, nil, "create", name)
				run(nil, nil, "show", name)
				run(nil, nil, "set-finished", name)
				run(nil, nil, "delete", name)
			}
		})
		Expect(runtime.Seconds()).To(BeNumerically("<", 2))
	}, 5)
})
