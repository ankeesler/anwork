package task2

// A State describes the status of some Task.
type State int

// These are the states that a Task could be in.
const (
	StateReady    State = 0
	StateBlocked  State = 1
	StateRunning  State = 2
	StateFinished State = 3
)

// These are the names of the State's that a Task can occupy, indexed by the State integer value.
var StateNames = [...]string{
	"Ready",
	"Blocked",
	"Running",
	"Finished",
}

// This is the default priority that a Task gets when created.
const DefaultPriority = 10

// A Task is something that someone is working on. It could be something like "mow the lawn" or "buy
// sister a holiday present." A Task also has a priority which describes its relative importance to
// all other Task's.
type Task struct {
	// The name of the Task, i.e., "mow the lawn" or "PROJECT-123-fix-infinite-recursion."
	Name string `json:"name"`

	// This is a unique ID. Every Task has a different ID.
	ID int `json:"id"`

	// This is a description of the Task.
	Description string `json:"description"`

	// This is when the Task was created, represented by the number of seconds since January 1, 1970.
	StartDate int64 `json:"startDate"`

	// This is the priority of the Task. The lower the number, the higher the importance.
	Priority int `json:"priority"`

	// This is the State of the Task. See State* for possible values. A Task can go through any
	// number of State changes over the course of its life. All Tasks start out in the StateReady
	// State.
	State State `json:"state"`
}

// An EventType describes the type of Event that took place in the Manager.
type EventType int

// These are the types of Event's that can occur.
const (
	EventTypeCreate      EventType = 0
	EventTypeDelete      EventType = 1
	EventTypeSetState    EventType = 2
	EventTypeNote        EventType = 3
	EventTypeSetPriority EventType = 4
)

// An Event is something that took place. Each Event is associated with only one Task.
type Event struct {
	// A string description of the Event.
	Title string `json:"title"`
	// The time that the Event took place, represented by the number of seconds since January 1, 1970. This is unique for each Event.
	Date int64 `json:"date"`
	// The type of Event.
	Type EventType `json:"type"`
	// The ID of the Task to which this Event refers.
	TaskID int `json:"taskid"`
}

// Repo is an object that allows CRUD operations on Task's and Event's.
type Repo interface {
	CreateTask(*Task) error
	Tasks() ([]*Task, error)
	FindTaskByID(int) (*Task, error)
	FindTaskByName(string) (*Task, error)
	UpdateTask(*Task) error
	DeleteTask(*Task) error

	CreateEvent(*Event) error
	Events() ([]*Event, error)
	DeleteEvent(*Event) error
}
