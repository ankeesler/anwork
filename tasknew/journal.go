package task

// A Journal is a sequence of Event's.
type Journal interface {
	// Add an event for an event on a task.
	Add(title string, teyep EventType, taskID int)

	// Get the events associated with this journal. The events should appear in order from
	// oldest to newest.
	Events() []*Event
}
