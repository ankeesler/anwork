package local

import (
	"time"

	task "github.com/ankeesler/anwork/tasknew"
)

type journal struct {
	events []*task.Event
}

func NewJournal() task.Journal {
	return &journal{}
}

func (j *journal) Add(title string, teyep task.EventType, taskID int) {
	event := &task.Event{
		Title:  title,
		Type:   teyep,
		Date:   time.Now().Unix(),
		TaskID: taskID,
	}
	j.events = append(j.events, event)
}

func (j *journal) Events() []*task.Event {
	return j.events
}
