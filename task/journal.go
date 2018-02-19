package task

import (
	"encoding/json"
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

// An Event is something that took place. It is stored in a Journal. Each Event refers to only one
// Task.
type Event struct {
	// A string description of the Event.
	Title string `json:"title"`
	// The time that the Event took place.
	Date time.Time `json:"date"`
	// The type of Event.
	Type EventType `json:"type"`
	// The ID of the Task to which this Event refers.
	TaskID int `json:"taskid"`
}

func newEvent(title string, teyep EventType, taskID int) *Event {
	e := &Event{
		Title:  title,
		Date:   time.Now(),
		Type:   teyep,
		TaskID: taskID,
	}

	// Truncate the start time at the seconds since we only persist the seconds amount.
	e.Date = e.Date.Truncate(time.Second)

	return e
}

func (e *Event) Serialize() ([]byte, error) {
	return json.Marshal(e)
}

func (e *Event) toProtobuf(eProtobuf *pb.Event) {
	eProtobuf.Title = e.Title
	eProtobuf.Date = e.Date.Unix()
	eProtobuf.Type = pb.EventType(e.Type)
	eProtobuf.TaskID = int32(e.TaskID)
}

func (e *Event) Unserialize(bytes []byte) error {
	eProtobuf := pb.Event{}
	err := proto.Unmarshal(bytes, &eProtobuf)
	if err == nil {
		e.fromProtobuf(&eProtobuf)
		return nil
	}

	if err := json.Unmarshal(bytes, e); err != nil {
		return err
	}
	return nil
}

func (e *Event) fromProtobuf(eProtobuf *pb.Event) {
	e.Title = eProtobuf.Title
	e.Date = time.Unix(eProtobuf.Date, 0) // sec, nsec
	e.Type = EventType(eProtobuf.Type)
	e.TaskID = int(eProtobuf.TaskID)

	noteTaskID(e.TaskID)
}

// A Journal is a sequence of Event's.
type Journal struct {
	Events []*Event `json:"events"`
}

func NewJournal() *Journal {
	return &Journal{Events: make([]*Event, 0)}
}

func (j *Journal) Serialize() ([]byte, error) {
	return json.Marshal(j)
}

func (j *Journal) toProtobuf(jProtobuf *pb.Journal) {
	var esProtobuf []*pb.Event = make([]*pb.Event, len(j.Events))
	for index, event := range j.Events {
		esProtobuf[index] = &pb.Event{}
		event.toProtobuf(esProtobuf[index])
	}

	jProtobuf.Events = esProtobuf
}

func (j *Journal) Unserialize(bytes []byte) error {
	jProtobuf := pb.Journal{}
	err := proto.Unmarshal(bytes, &jProtobuf)
	if err == nil {
		j.fromProtobuf(&jProtobuf)
		return nil
	}

	return json.Unmarshal(bytes, j)
}

func (j *Journal) fromProtobuf(jProtobuf *pb.Journal) {
	esProtobuf := jProtobuf.GetEvents()
	j.Events = make([]*Event, len(esProtobuf))
	for index, eProtobuf := range esProtobuf {
		j.Events[index] = &Event{}
		j.Events[index].fromProtobuf(eProtobuf)
	}
}
