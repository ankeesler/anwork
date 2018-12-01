package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("SetState", func() {
	var (
		outBuf, errBuf *gbytes.Buffer
	)

	BeforeEach(func() {
		outBuf = gbytes.NewBuffer()
		errBuf = gbytes.NewBuffer()

		run(nil, nil, "create", "task-a")
		run(nil, nil, "create", "task-b")
		run(nil, nil, "create", "task-c")
	})

	AfterEach(func() {
		run(nil, nil, "reset")
	})

	runTests := func(useCommandAliases bool) {
		Context("when changing the state on tasks", func() {
			BeforeEach(func() {
				if useCommandAliases {
					run(nil, nil, "sr", "task-a")
					run(nil, nil, "sf", "task-b")
					run(nil, nil, "sb", "task-c")
				} else {
					run(nil, nil, "set-running", "task-a")
					run(nil, nil, "set-finished", "task-b")
					run(nil, nil, "set-blocked", "task-c")
				}
			})
			It("properly displays the states", func() {
				run(outBuf, errBuf, "show")
				Expect(outBuf).To(gbytes.Say("RUNNING tasks:\n  task-a"))
				Expect(outBuf).To(gbytes.Say("BLOCKED tasks:\n  task-c"))
				Expect(outBuf).To(gbytes.Say("FINISHED tasks:\n  task-b"))
			})
			It("properly records the states", func() {
				run(outBuf, errBuf, "show", "task-a")
				Expect(outBuf).To(gbytes.Say("Name: task-a\nID: \\d+\nCreated: .*\nPriority: 10\nState: RUNNING"))
				run(outBuf, errBuf, "show", "task-b")
				Expect(outBuf).To(gbytes.Say("Name: task-b\nID: \\d+\nCreated: .*\nPriority: 10\nState: FINISHED"))
				run(outBuf, errBuf, "show", "task-c")
				Expect(outBuf).To(gbytes.Say("Name: task-c\nID: \\d+\nCreated: .*\nPriority: 10\nState: BLOCKED"))
			})
			It("records the events in each of the task's journals", func() {
				run(outBuf, errBuf, "journal", "task-a")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set state on task 'task-a' from Ready to Running"))
				run(outBuf, errBuf, "journal", "task-b")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set state on task 'task-b' from Ready to Finished"))
				run(outBuf, errBuf, "journal", "task-c")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set state on task 'task-c' from Ready to Blocked"))
			})
			It("records the events in the global journal in order from newest to oldest", func() {
				run(outBuf, errBuf, "journal")
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set state on task 'task-c' from Ready to Blocked"))
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set state on task 'task-b' from Ready to Finished"))
				Expect(outBuf).To(gbytes.Say("\\[.*\\]: Set state on task 'task-a' from Ready to Running"))
			})
		})
	}

	runTests(false)
	runTests(true)
})
