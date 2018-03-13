package task_test

import (
	"encoding/json"

	"github.com/ankeesler/anwork/storage"
	"github.com/ankeesler/anwork/task"
	pb "github.com/ankeesler/anwork/task/proto"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	goodProtobufEventContext   = "good-protobuf-event-context"
	goodEventContext           = "good-event-context"
	goodProtobufJournalContext = "good-protobuf-journal-context"
	goodJournalContext         = "good-journal-context"
)

var _ = Describe("EventType's", func() {
	It("lines up with the protocol buffer definitions", func() {
		Expect(task.EventTypeCreate).To(Equal(task.EventType(pb.EventType_CREATE)))
		Expect(task.EventTypeDelete).To(Equal(task.EventType(pb.EventType_DELETE)))
		Expect(task.EventTypeSetState).To(Equal(task.EventType(pb.EventType_SET_STATE)))
		Expect(task.EventTypeNote).To(Equal(task.EventType(pb.EventType_NOTE)))
		Expect(task.EventTypeSetPriority).To(Equal(task.EventType(pb.EventType_SET_PRIORITY)))
	})
})
var _ = Describe("Event's", func() {

	Describe("Persistence", func() {
		var (
			eventA = task.Event{
				Title: "Here is Event-A's title",
				Date:  12345,
				Type:  task.EventTypeNote}
			eventB = task.Event{
				Title:  "Here is Event-B's title",
				Date:   54321,
				Type:   task.EventTypeSetPriority,
				TaskID: 57}
			p storage.FilePersister = storage.FilePersister{Root: root}
		)

		It("are persistable", func() {
			Expect(p.Exists(tmpContext)).To(BeFalse(),
				"Cannot run this test when context (%s) already exists", tmpContext)
			defer p.Delete(tmpContext)

			var e task.Event
			Expect(p.Persist(tmpContext, &eventA)).To(Succeed())
			Expect(p.Unpersist(tmpContext, &e)).To(Succeed())
			Expect(eventA).To(Equal(e))
		})

		It("are unpersistable via protocol buffers (legacy)", func() {
			Expect(p.Exists(goodProtobufEventContext)).To(BeTrue(),
				"Cannot run this test when context (%s) does not exist", goodProtobufEventContext)

			var e task.Event
			Expect(p.Unpersist(goodProtobufEventContext, &e)).To(Succeed())
			Expect(eventB).To(Equal(e))
		})

		It("are serialized via json", func() {
			e := task.Event{
				Title:  "title a",
				Date:   1000,
				Type:   task.EventTypeNote,
				TaskID: 15,
			}
			bytes, err := e.Serialize()
			Expect(err).NotTo(HaveOccurred())

			var jsonE task.Event
			Expect(json.Unmarshal(bytes, &jsonE)).To(Succeed())
			Expect(jsonE).To(Equal(e))
		})

		It("are unpersistable via json", func() {
			Expect(p.Exists(goodEventContext)).To(BeTrue(),
				"Cannot run this test when context (%s) does not exist", goodEventContext)

			var e task.Event
			Expect(p.Unpersist(goodEventContext, &e)).To(Succeed())
			Expect(eventB).To(Equal(e))
		})

		It("fails gracefully when loaded from a bad context", func() {
			Expect(p.Exists(badContext)).To(BeTrue(),
				"Cannot run this test when context (%s) does not exist", badContext)
			Expect(p.Unpersist(badContext, &task.Event{})).ToNot(Succeed())
		})
	})
})
var _ = Describe("Journal", func() {
	var (
		j, tmpJ *task.Journal
		p       storage.FilePersister = storage.FilePersister{Root: root}
	)
	BeforeEach(func() {
		j = task.NewJournal()
		tmpJ = task.NewJournal()
	})
	It("holds no events to start", func() {
		Expect(j.Events).To(BeEmpty())
	})
	Context("when adding an event", func() {
		e0 := &task.Event{Title: "event 0", Type: task.EventTypeSetPriority, TaskID: 0}
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
			e1 := &task.Event{Title: "event 1", Type: task.EventTypeCreate, TaskID: 0}
			e2 := &task.Event{Title: "event 2", Type: task.EventTypeNote, TaskID: 0}
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

	It("serializes to json", func() {
		j.Events = append(j.Events, &task.Event{
			Title: "event a",
			Date:  10000,
			Type:  task.EventTypeNote,
		})
		j.Events = append(j.Events, &task.Event{
			Title: "event b",
			Date:  20000,
			Type:  task.EventTypeSetState,
		})
		j.Events = append(j.Events, &task.Event{
			Title: "event c",
			Date:  30000,
			Type:  task.EventTypeSetPriority,
		})
		j.Events = append(j.Events, &task.Event{
			Title: "event b",
			Date:  40000,
			Type:  task.EventTypeDelete,
		})
		bytes, err := j.Serialize()
		Expect(err).NotTo(HaveOccurred())

		var jsonJ task.Journal
		Expect(json.Unmarshal(bytes, &jsonJ)).To(Succeed())
		Expect(&jsonJ).To(Equal(j))
	})

	It("unpersists correctly from protocol buffers (legacy)", func() {
		Expect(p.Exists(goodProtobufJournalContext)).To(BeTrue(),
			"Cannot run this test when context (%s) does not exist", goodProtobufJournalContext)

		j.Events = append(j.Events, &task.Event{
			Title: "event a",
			Date:  10000,
			Type:  task.EventTypeNote,
		})
		j.Events = append(j.Events, &task.Event{
			Title: "event b",
			Date:  20000,
			Type:  task.EventTypeSetState,
		})
		j.Events = append(j.Events, &task.Event{
			Title: "event c",
			Date:  30000,
			Type:  task.EventTypeSetPriority,
		})
		j.Events = append(j.Events, &task.Event{
			Title: "event b",
			Date:  40000,
			Type:  task.EventTypeDelete,
		})
		Expect(p.Unpersist(goodProtobufJournalContext, tmpJ)).To(Succeed())
		Expect(j).To(Equal(tmpJ))
	})

	It("unpersists correctly from json", func() {
		Expect(p.Exists(goodJournalContext)).To(BeTrue(),
			"Cannot run this test when context (%s) does not exist", goodJournalContext)

		j.Events = append(j.Events, &task.Event{
			Title:  "event a",
			Date:   10000,
			Type:   task.EventTypeNote,
			TaskID: 3,
		})
		j.Events = append(j.Events, &task.Event{
			Title:  "event b",
			Date:   20000,
			Type:   task.EventTypeSetState,
			TaskID: 2,
		})
		j.Events = append(j.Events, &task.Event{
			Title:  "event c",
			Date:   30000,
			Type:   task.EventTypeSetPriority,
			TaskID: 1,
		})
		j.Events = append(j.Events, &task.Event{
			Title:  "event b",
			Date:   40000,
			Type:   task.EventTypeDelete,
			TaskID: 0,
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
