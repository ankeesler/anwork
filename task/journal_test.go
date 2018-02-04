package task

import (
	"time"

	"github.com/ankeesler/anwork/storage"
	pb "github.com/ankeesler/anwork/task/proto"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	goodEventContext   = "good-event-context"
	goodJournalContext = "good-journal-context"
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
		p        storage.FilePersister = storage.FilePersister{Root: root}
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
		Expect(p.Exists(goodEventContext)).To(BeTrue(),
			"Cannot run this test when context (%s) does not exist", goodEventContext)

		Expect(p.Unpersist(goodEventContext, &tmpEvent)).To(Succeed())
		Expect(eventB).To(Equal(tmpEvent))
	})
	It("fails gracefully when loaded from a bad context", func() {
		Expect(p.Exists(badContext)).To(BeTrue(),
			"Cannot run this test when context (%s) does not exist", badContext)
		Expect(p.Unpersist(badContext, &Event{})).ToNot(Succeed())
	})
})
var _ = Describe("Journal", func() {
	var (
		j, tmpJ *Journal
		p       storage.FilePersister = storage.FilePersister{Root: root}
	)
	BeforeEach(func() {
		j = newJournal()
		tmpJ = newJournal()
	})
	It("holds no events to start", func() {
		Expect(j.Events).To(BeEmpty())
	})
	Context("when adding an event", func() {
		e0 := newEvent("event 0", EventTypeSetPriority, 0)
		BeforeEach(func() {
			j.Events = append(j.Events, e0)
		})
		It("holds one event", func() {
			Expect(j.Events).To(HaveLen(1))
			actualE0 := j.Events[0]
			Expect(actualE0.Title).To(Equal("event 0"))
		})
		It("persists that event correctly", func() {
			Expect(p.Exists(tmpContext)).To(BeFalse(),
				"Cannot run this test when context (%s) already exists", tmpContext)
			defer p.Delete(tmpContext)

			Expect(p.Persist(tmpContext, j)).To(Succeed())
			Expect(p.Unpersist(tmpContext, tmpJ)).To(Succeed())
			Expect(j).To(Equal(tmpJ))
		})
		Context("when adding more events", func() {
			e1 := newEvent("event 1", EventTypeCreate, 0)
			e2 := newEvent("event 2", EventTypeNote, 0)
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
			It("persists those events correctly", func() {
				Expect(p.Exists(tmpContext)).To(BeFalse(),
					"Cannot run this test when context (%s) already exists", tmpContext)
				defer p.Delete(tmpContext)

				Expect(p.Persist(tmpContext, j)).To(Succeed())
				Expect(p.Unpersist(tmpContext, tmpJ)).To(Succeed())
				Expect(j).To(Equal(tmpJ))
			})
		})
	})

	It("unpersists correctly", func() {
		Expect(p.Exists(goodJournalContext)).To(BeTrue(),
			"Cannot run this test when context (%s) does not exist", goodJournalContext)

		j.Events = append(j.Events, &Event{
			Title: "event a",
			Date:  time.Unix(10000, 0),
			Type:  EventTypeNote,
		})
		j.Events = append(j.Events, &Event{
			Title: "event b",
			Date:  time.Unix(20000, 0),
			Type:  EventTypeSetState,
		})
		j.Events = append(j.Events, &Event{
			Title: "event c",
			Date:  time.Unix(30000, 0),
			Type:  EventTypeSetPriority,
		})
		j.Events = append(j.Events, &Event{
			Title: "event b",
			Date:  time.Unix(40000, 0),
			Type:  EventTypeDelete,
		})
		Expect(p.Unpersist(goodJournalContext, tmpJ)).To(Succeed())
		Expect(j).To(Equal(tmpJ))
	})

	It("gracefully fails when the context is bad", func() {
		Expect(p.Exists(badContext)).To(BeTrue(),
			"Cannot run this test when context (%s) does not exist", badContext)
		Expect(p.Unpersist(badContext, tmpJ)).ToNot(Succeed())
	})
})
