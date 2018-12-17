// Package task contains the domain objects for anwork: Task's and Event's.
//
// A Task is something that someone is working on. It could be something like "mow the lawn" or "buy
// sister a holiday present."
//
// Every Task is in one of a number of different State's: Ready, Blocked, Running, or Finished. A
// Task also has a priority which describes its relative importance to all other Task's.
//
// An Event is something that happened to a Task.
package task

// A State describes the status of some Task.
type State string

// These are the states that a Task could be in.
const (
	StateReady    State = "Ready"
	StateBlocked        = "Blocked"
	StateRunning        = "Running"
	StateFinished       = "Finished"
)

// A Task is something that someone is working on. It could be something like "mow the lawn" or "buy
// sister a holiday present." A Task also has a priority which describes its relative importance to
// all other Task's.
type Task struct {
	// The name of the Task, i.e., "mow the lawn" or "PROJECT-123-fix-infinite-recursion."
	Name string `json:"name"`

	// This is a unique ID. Every Task has a different ID.
	ID int `json:"id"`

	// This is when the Task was created, represented by the number of seconds since January 1, 1970.
	StartDate int64 `json:"startDate"`

	// This is the priority of the Task. The lower the number, the higher the importance.
	Priority int `json:"priority"`

	// This is the State of the Task. See State* for possible values. A Task can go through any
	// number of State changes over the course of its life.
	State State `json:"state"`
}

// An EventType describes the type of Event that took place in the Manager.
type EventType int

// These are the types of Event's that can occur.
const (
	EventTypeCreate = iota
	EventTypeDelete
	EventTypeSetState
	EventTypeNote
	EventTypeSetPriority
)

// An Event is something that took place. Each Event is associated with only one Task.
type Event struct {
	// Unique identifier for the Event.
	ID int
	// A string description of the Event.
	Title string `json:"title"`
	// The time that the Event took place, represented by the number of seconds since January 1, 1970.
	Date int64 `json:"date"`
	// The type of Event.
	Type EventType `json:"type"`
	// The ID of the Task to which this Event refers.
	TaskID int `json:"taskid"`
}
