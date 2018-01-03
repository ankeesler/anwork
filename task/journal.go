package task

import (
	"time"

	pb "github.com/ankeesler/anwork/task/proto"
	"github.com/golang/protobuf/proto"
)

//go:generate protoc --proto_path=proto --go_out=proto task.proto

// An EventType describes the type of Event that took place in the Manager.
type EventType int

const (
	EventTypeCreate      = EventType(0)
	EventTypeDelete      = EventType(1)
	EventTypeSetState    = EventType(2)
	EventTypeNote        = EventType(3)
	EventTypeSetPriority = EventType(4)
)

// An Event is something that took place. It is stored in a Journal.
type Event struct {
	Title string
	Date  time.Time
	Type  EventType
}

func newEvent(title string, teyep EventType) *Event {
	e := &Event{
		Title: title,
		Date:  time.Now(),
		Type:  teyep,
	}

	// Truncate the start time at the seconds since we only persist the seconds amount.
	e.Date = e.Date.Truncate(time.Second)

	return e
}

func (e *Event) Serialize() ([]byte, error) {
	eProtobuf := pb.Event{
		Title: e.Title,
		Date:  e.Date.Unix(),
		Type:  pb.EventType(e.Type),
	}
	return proto.Marshal(&eProtobuf)
}

func (e *Event) Unserialize(bytes []byte) error {
	eProtobuf := pb.Event{}
	err := proto.Unmarshal(bytes, &eProtobuf)
	if err != nil {
		return err
	}

	e.Title = eProtobuf.Title
	e.Date = time.Unix(eProtobuf.Date, 0) // sec, nsec
	e.Type = EventType(eProtobuf.Type)

	return nil
}

// A Journal is a sequence of Event's.
type Journal struct {
	Events []*Event
}

func (j *Journal) Serialize() ([]byte, error) {
	var esProtobuf []*pb.Event = make([]*pb.Event, len(j.Events))
	for index, event := range j.Events {
		eProtobuf := &pb.Event{
			Title: event.Title,
			Date:  event.Date.Unix(),
			Type:  pb.EventType(event.Type),
		}
		esProtobuf[index] = eProtobuf
	}

	jProtobuf := pb.Journal{
		Events: esProtobuf,
	}
	return proto.Marshal(&jProtobuf)
}

func (j *Journal) Unserialize(bytes []byte) error {
	jProtobuf := pb.Journal{}
	err := proto.Unmarshal(bytes, &jProtobuf)
	if err != nil {
		return err
	}

	esProtobuf := jProtobuf.GetEvents()
	j.Events = make([]*Event, len(esProtobuf))
	for index, eProtobuf := range esProtobuf {
		j.Events[index] = &Event{}
		j.Events[index].Title = eProtobuf.Title
		j.Events[index].Date = time.Unix(eProtobuf.Date, 0) // sec, nsec
		j.Events[index].Type = EventType(eProtobuf.Type)
	}

	return nil
}
