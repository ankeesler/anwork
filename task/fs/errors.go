package fs

import "fmt"

type unknownTaskError struct {
	name string
	id   int
}

func (ute *unknownTaskError) Error() string {
	return fmt.Sprintf("unknown task with name '%s' and id %d", ute.name, ute.id)
}

type duplicateTaskError struct {
	name string
	id   int
}

func (dte *duplicateTaskError) Error() string {
	return fmt.Sprintf("duplicate task with name '%s' and id %d", dte.name, dte.id)
}

type unknownEventError struct {
	title string
	date  int64
}

func (uee *unknownEventError) Error() string {
	return fmt.Sprintf("unknown event with title '%s' and date %d", uee.title, uee.date)
}

type duplicateEventError struct {
	title string
	date  int64
}

func (dee *duplicateEventError) Error() string {
	return fmt.Sprintf("duplicate event with title '%s' and date %d", dee.title, dee.date)
}
