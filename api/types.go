package api

import "github.com/ankeesler/anwork/task"

// This is a generic error message used in a response.
type ErrorResponse struct {
	Message string `json:"message"`
}

// Return the message that describes this error.
func (e ErrorResponse) Error() string {
	return e.Message
}

// This is the payload in a POST to /tasks. It is used to create a new task
// with a name.
type CreateRequest struct {
	Name string `json:"name"`
}

// This is the payload in a PUT to /tasks/:name. It is used to update the
// priority, state, or name of a task.
type UpdateTaskRequest struct {
	Priority int        `json:"priority"`
	State    task.State `json:"state"`
	Name     string     `json:"name"`
}

// This is the payload in a PUT to /events. It is used to add a new event.
type AddEventRequest struct {
	Title  string
	Date   int64
	Type   task.EventType
	TaskID int
}
