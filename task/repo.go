package task

//go:generate counterfeiter . Repo

// Repo is an object that allows CRUD operations on Task's and Event's.
type Repo interface {
	// CreateTask creates a Task. The Task.ID field is set by the Repo.
	CreateTask(*Task) error
	// Tasks returns all of the Task's in this Repo.
	Tasks() ([]*Task, error)
	// FindTaskByID tries to find a Task with the provided ID. If the Task does not
	// exist, it will return nil, nil.
	FindTaskByID(int) (*Task, error)
	// FindTaskByName tries to find a Task with the provided name. If the Task does not
	// exist, it will return nil, nil.
	FindTaskByName(string) (*Task, error)
	// UpdateTask finds the Task in the Repo with the provided ID and updates its
	// values to those provided.
	UpdateTask(*Task) error
	// DeleteTask deletes a Task with the provided ID.
	// If the Task does not exist, this function will return nil.
	DeleteTask(*Task) error

	// CreateEvent creates a Event. The Event.ID field is set by the Repo.
	CreateEvent(*Event) error
	// Events returns all of the Event's in this Repo.
	Events() ([]*Event, error)
	// FindEventByID will try to find an Event with the provided ID. If the Event does
	// not exist, it will return nil, nil.
	FindEventByID(int) (*Event, error)
	// DeleteEvent deletes an Event with the provided ID.
	// If the Event does not exist, this function will return nil.
	DeleteEvent(*Event) error
}
