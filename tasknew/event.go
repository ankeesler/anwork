package task

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
	// The time that the Event took place, represented by the number of seconds since January 1, 1970.
	Date int64 `json:"date"`
	// The type of Event.
	Type EventType `json:"type"`
	// The ID of the Task to which this Event refers.
	TaskID int `json:"taskid"`
}
