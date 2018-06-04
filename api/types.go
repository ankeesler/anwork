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
// priority or the state of a task.
type SetRequest struct {
	Priority int        `json:"priority"`
	State    task.State `json:"state"`
}
