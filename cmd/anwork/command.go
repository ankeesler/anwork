package main

import (
	"github.com/ankeesler/anwork/task"
)

// A command is a keyword (see name field) passed to the anwork executable that provokes some
// functionality (see action field).
type command struct {
	name, usage string
	action      func(name string, manager *task.Manager) bool
}
