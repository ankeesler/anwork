package task

//go:generate counterfeiter . ManagerFactory
//go:generate counterfeiter . Manager

// A ManagerFactory is an object that can create and save Manager instances.
type ManagerFactory interface {
	// Create a Manager, or return an error.
	Create() (Manager, error)
	// Save a Manager back into the factory.
	Save(Manager) error
	// Completely forget all of the state associated with this factory.
	Reset() error
}

// A Manager is an interface through which Task's can be created, read, updated, and deleted.
type Manager interface {
	// Create a task with a name. Return an error if the task name is not unique.
	Create(name string) error

	// Delete a task with a name. Returns an error if the task was not able to be deleted.
	Delete(name string) error

	// Find a task with an ID.
	FindByID(id int) *Task
	// Find a task with a name.
	FindByName(name string) *Task

	// Get all of the Tasks contained in this manager, ordered from highest priority (lowest integer
	// value) to lowest priority (highest integer value).
	//
	// When multiple tasks have the same priority, the Task's will be ordered by their (unique) ID in
	// ascending order. This means that the older Task's will come first. This is a conscious decision.
	// The Task's that have been around the longest are assumed to need to be completed first.
	//
	// This function will never return nil!
	Tasks() []*Task

	// Add a note for a task.
	Note(name, note string) error
	// Set the priority of a task.
	SetPriority(name string, priority int) error
	// Set the state of a task.
	SetState(name string, state State) error

	// Get the events associated with this manager.
	Events() []*Event
}
