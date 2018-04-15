// This package contains all of the task-related declarations in the anwork project.
//
// A Task is something that someone is working on. It could be something like "mow the lawn" or "buy
// sister a holiday present."
//
// Every Task is in one of a number of different State's: Waiting, Blocked, Running, or Finished. A
// Task also has a priority which describes its relative importance to all other Task's.
//
// A Manager is an interface through which Task's can be created, read, updated, and deleted.
//
// A Manager also keeps track of the changes that are made to the Task's it oversees. It keeps a list
// of things (Event's) that happen to a Manager (i.e., a note is added, a Task is created, a Task is
// updated, etc.).
package task

// A State describes the status of some Task.
type State int

// These are the states that a Task could be in.
const (
	StateWaiting  State = 0
	StateBlocked  State = 1
	StateRunning  State = 2
	StateFinished State = 3
)

// These are the names of the State's that a Task can occupy, indexed by the State integer value.
var StateNames = [...]string{
	"Waiting",
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
	// number of State changes over the course of its life. All Tasks start out in the StateWaiting
	// State.
	State State `json:"state"`
}
