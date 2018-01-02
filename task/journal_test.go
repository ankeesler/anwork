package task

import (
	"time"

	"github.com/ankeesler/anwork/storage"
	pb "github.com/ankeesler/anwork/task/proto"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	eventContext = "event-context"
)

var _ = Describe("EventType's", func() {
	It("lines up with the protocol buffer definitions", func() {
		Expect(EventTypeCreate).To(Equal(EventType(pb.EventType_CREATE)))
		Expect(EventTypeDelete).To(Equal(EventType(pb.EventType_DELETE)))
		Expect(EventTypeSetState).To(Equal(EventType(pb.EventType_SET_STATE)))
		Expect(EventTypeNote).To(Equal(EventType(pb.EventType_NOTE)))
		Expect(EventTypeSetPriority).To(Equal(EventType(pb.EventType_SET_PRIORITY)))
	})
})
var _ = Describe("Event's", func() {
	var (
		eventA = Event{
			Title: "Here is Event-A's title",
			Date:  time.Unix(12345, 0),
			Type:  EventTypeNote}
		eventB = Event{
			Title: "Here is Event-B's title",
			Date:  time.Unix(54321, 0),
			Type:  EventTypeSetPriority}
		tmpEvent Event
		p        storage.Persister = storage.Persister{root}
	)
	It("are persistable", func() {
		Expect(p.Exists(tmpContext)).To(BeFalse(),
			"Cannot run this test when context (%s) already exists", tmpContext)
		defer p.Delete(tmpContext)

		Expect(p.Persist(tmpContext, &eventA)).To(Succeed())
		Expect(p.Unpersist(tmpContext, &tmpEvent)).To(Succeed())
		Expect(eventA).To(Equal(tmpEvent))
	})
	It("are unpersistable", func() {
		Expect(p.Exists(eventContext)).To(BeTrue(),
			"Cannot run this test when context (%s) does not exist", eventContext)

		Expect(p.Unpersist(eventContext, &tmpEvent)).To(Succeed())
		Expect(eventB).To(Equal(tmpEvent))
	})
})
var _ = Describe("Journal", func() {
	var j *Journal
	BeforeEach(func() {
		j = &Journal{}
	})
	It("holds no events to start", func() {
		Expect(j.Events).To(BeEmpty())
	})
	Context("when adding an event", func() {
		e0 := &Event{Title: "event 0", Date: time.Now()}
		BeforeEach(func() {
			j.Events = append(j.Events, e0)
		})
		It("holds one event", func() {
			Expect(j.Events).To(HaveLen(1))
			actualE0 := j.Events[0]
			Expect(actualE0.Title).To(Equal("event 0"))
		})
		Context("when adding more events", func() {
			e1 := &Event{Title: "event 1", Date: time.Now()}
			e2 := &Event{Title: "event 2", Date: time.Now()}
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
