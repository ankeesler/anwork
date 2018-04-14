package task

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

// This is a utility method to run tests with an object conforming to the
// Journal interface.
func RunJournalTests(createJournal func() Journal) {
	var (
		journal Journal
	)
	BeforeEach(func() {
		journal = createJournal()
	})

	Context("when there have been no events added", func() {
		It("returns no events", func() {
			Expect(journal.Events()).To(BeEmpty())
		})
	})

	Context("when adding an event", func() {
		BeforeEach(func() {
			journal.Add("1", EventTypeCreate, 1)
			journal.Add("2", EventTypeDelete, 2)
			journal.Add("3", EventTypeSetState, 3)
		})

		It("returns the added events", func() {
			events := journal.Events()
			Expect(events).To(HaveLen(3))
			Expect(events[0].Title).To(Equal("1"))
			Expect(events[1].Title).To(Equal("2"))
			Expect(events[2].Title).To(Equal("3"))
		})
	})
}
