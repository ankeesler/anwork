package task

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Journal", func() {
	var j *Journal
	BeforeEach(func() {
		j = &Journal{}
	})
	It("holds no events to start", func() {
		Expect(j.Events).To(BeEmpty())
	})
	Context("when adding an event", func() {
		e0 := &Event{Title: "event 0", T: time.Now()}
		BeforeEach(func() {
			j.Events = append(j.Events, e0)
		})
		It("holds one event", func() {
			Expect(j.Events).To(HaveLen(1))
			actualE0 := j.Events[0]
			Expect(actualE0.Title).To(Equal("event 0"))
		})
		Context("when adding more events", func() {
			e1 := &Event{Title: "event 1", T: time.Now()}
			e2 := &Event{Title: "event 2", T: time.Now()}
			BeforeEach(func() {
				j.Events = append(j.Events, e1)
				j.Events = append(j.Events, e2)
			})
			It("holds three events", func() {
				Expect(j.Events).To(HaveLen(3))
			})
			It("stores the events in order from oldest to newest", func() {
				Expect(j.Events[0].Title).To(Equal("event 0"))
				Expect(j.Events[1].Title).To(Equal("event 1"))
				Expect(j.Events[2].Title).To(Equal("event 2"))
			})
		})
	})
})
