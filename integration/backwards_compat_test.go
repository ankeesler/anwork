package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Backwards compatibility", func() {
	var (
		outBuf, errBuf *gbytes.Buffer
	)
	BeforeEach(func() {
		outBuf = gbytes.NewBuffer()
		errBuf = gbytes.NewBuffer()
	})

	Context("version 2", func() {
		It("shows the correct tasks", func() {
			run(outBuf, errBuf, "-o", "data", "show")
			Expect(outBuf).To(gbytes.Say("RUNNING tasks:"))
			Expect(outBuf).To(gbytes.Say("  task-c \\(2\\)"))
			Expect(outBuf).To(gbytes.Say("BLOCKED tasks:"))
			Expect(outBuf).To(gbytes.Say("  task-b \\(1\\)"))
			Expect(outBuf).To(gbytes.Say("WAITING tasks:"))
			Expect(outBuf).To(gbytes.Say("FINISHED tasks:"))
			Expect(outBuf).To(gbytes.Say("  task-a \\(0\\)"))
		})
		It("shows the correct task details", func() {
			if runningOnTravis {
				Skip("I don't think the date/time is set properly when we run in Travis CI...")
			}
			run(outBuf, errBuf, "-o", "data", "show", "task-a")
			Expect(outBuf).To(gbytes.Say("Name: task-a\nID: 0\nCreated: Monday January 15 18:59\nPriority: 10\nState: FINISHED"))
			run(outBuf, errBuf, "-o", "data", "show", "task-b")
			Expect(outBuf).To(gbytes.Say("Name: task-b\nID: 1\nCreated: Monday January 15 18:59\nPriority: 10\nState: BLOCKED"))
			run(outBuf, errBuf, "-o", "data", "show", "task-c")
			Expect(outBuf).To(gbytes.Say("Name: task-c\nID: 2\nCreated: Monday January 15 18:59\nPriority: 10\nState: RUNNING"))
		})
		It("shows the correct journal", func() {
			if runningOnTravis {
				Skip("I don't think the date/time is set properly when we run in Travis CI...")
			}
			run(outBuf, errBuf, "-o", "data", "journal")
			Expect(outBuf).To(gbytes.Say("\\[Monday January 15 19:00\\]: Set state on task task-c from Waiting to Running"))
			Expect(outBuf).To(gbytes.Say("\\[Monday January 15 19:00\\]: Set state on task task-a from Running to Finished"))
			Expect(outBuf).To(gbytes.Say("\\[Monday January 15 19:00\\]: Set state on task task-b from Waiting to Blocked"))
			Expect(outBuf).To(gbytes.Say("\\[Monday January 15 18:59\\]: Set state on task task-a from Waiting to Running"))
			Expect(outBuf).To(gbytes.Say("\\[Monday January 15 18:59\\]: Created task task-c"))
			Expect(outBuf).To(gbytes.Say("\\[Monday January 15 18:59\\]: Created task task-b"))
			Expect(outBuf).To(gbytes.Say("\\[Monday January 15 18:59\\]: Created task task-a"))
		})
	})
})
