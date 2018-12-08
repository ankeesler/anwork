package manager

import "fmt"

type unknownTaskError struct {
	name string
}

func (ute unknownTaskError) Error() string {
	return fmt.Sprintf("unknown task with name '%s'", ute.name)
}
