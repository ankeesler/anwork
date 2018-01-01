package task

import (
	"time"
)

// An Event is something that took place. It is stored in a Journal.
type Event struct {
	Title string
	T     time.Time
}

// A Journal is a sequence of Event's.
type Journal struct {
	Events []*Event
}
