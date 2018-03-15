package integration

import (
	"fmt"

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

	runTests := func(version int) {
		Context(fmt.Sprintf("version %d", version), func() {
			context := fmt.Sprintf("v%d-context", version)
			It("shows the correct tasks", func() {
				run(outBuf, errBuf, "-o", "data", "-c", context, "show")
				Expect(outBuf).To(gbytes.Say("RUNNING tasks:"))
				Expect(outBuf).To(gbytes.Say("  task-c \\(2\\)"))
				Expect(outBuf).To(gbytes.Say("BLOCKED tasks:"))
				Expect(outBuf).To(gbytes.Say("  task-b \\(1\\)"))
				Expect(outBuf).To(gbytes.Say("WAITING tasks:"))
				Expect(outBuf).To(gbytes.Say("FINISHED tasks:"))
				Expect(outBuf).To(gbytes.Say("  task-a \\(0\\)"))
			})
			It("shows the correct task details", func() {
				run(outBuf, errBuf, "-o", "data", "-c", context, "show", "task-a")
				Expect(outBuf).To(gbytes.Say("Name: task-a\nID: 0\nCreated: \\w+ \\w+ \\d\\d? \\d\\d?:\\d\\d?\nPriority: 10\nState: FINISHED"))
				run(outBuf, errBuf, "-o", "data", "-c", context, "show", "task-b")
				Expect(outBuf).To(gbytes.Say("Name: task-b\nID: 1\nCreated: \\w+ \\w+ \\d\\d? \\d\\d?:\\d\\d?\nPriority: 10\nState: BLOCKED"))
				run(outBuf, errBuf, "-o", "data", "-c", context, "show", "task-c")
				Expect(outBuf).To(gbytes.Say("Name: task-c\nID: 2\nCreated: \\w+ \\w+ \\d\\d? \\d\\d?:\\d\\d?\nPriority: 10\nState: RUNNING"))
			})
			It("shows the correct journal", func() {
				run(outBuf, errBuf, "-o", "data", "-c", context, "journal")
				Expect(outBuf).To(gbytes.Say("\\[\\w+ \\w+ \\d\\d? \\d\\d?:\\d\\d?\\]: Set state on task task-c from Waiting to Running"))
				Expect(outBuf).To(gbytes.Say("\\[\\w+ \\w+ \\d\\d? \\d\\d?:\\d\\d?\\]: Set state on task task-a from Running to Finished"))
				Expect(outBuf).To(gbytes.Say("\\[\\w+ \\w+ \\d\\d? \\d\\d?:\\d\\d?\\]: Set state on task task-b from Waiting to Blocked"))
				Expect(outBuf).To(gbytes.Say("\\[\\w+ \\w+ \\d\\d? \\d\\d?:\\d\\d?\\]: Set state on task task-a from Waiting to Running"))
				Expect(outBuf).To(gbytes.Say("\\[\\w+ \\w+ \\d\\d? \\d\\d?:\\d\\d?\\]: Created task task-c"))
				Expect(outBuf).To(gbytes.Say("\\[\\w+ \\w+ \\d\\d? \\d\\d?:\\d\\d?\\]: Created task task-b"))
				Expect(outBuf).To(gbytes.Say("\\[\\w+ \\w+ \\d\\d? \\d\\d?:\\d\\d?\\]: Created task task-a"))
			})
		})
	}

	runTests(2)
	runTests(3)

})
