// Package lag provides bare-bones logging functionality
// for the ANWORK project.
//
// The simple usage is as follows.
//   log := createLog() // stdlib log
//   l := lag.New(log, lag.I)
//   l.P(lag.I, "here is a formatted string %s", "foo")
package lag

import (
	"fmt"
	"log"
)

// Level describes the heat of the information being logged.
// Is it just an fyi (Debug), is it info (Info), it is vital
// (Error), etc.
type Level int

// These are the supported log Level's.
const (
	N Level = iota // Don't print anything
	E              // Only print errors
	I              // Print helpful information and errors
	D              // Print everything
)

// L is a log. You can print formatted stuff to it with P().
type L struct {
	Level // The logging Level at which this L will operate.

	log *log.Logger
}

// New creates an L that writes to the provided *log.Logger.
// It is initialized with the Default log Level.
//
// It is totally not thread-safe.
func New(log *log.Logger, level Level) *L {
	return &L{log: log, Level: level}
}

// P a printf-formatted message to the log.Logger that this L was
// initialized with. The message will only be printed if this L's
// current Level is greater than or equal to the provided Level.
func (l *L) P(level Level, message string, stuff ...interface{}) {
	if l.Level >= level {
		msg := fmt.Sprintf(message, stuff...)
		l.log.Printf("%s: %s", formatLevel(level), msg)
	}
}

func formatLevel(level Level) string {
	switch level {
	case N:
		return "NONE"
	case E:
		return "ERROR"
	case I:
		return "INFO"
	case D:
		return "DEBUG"
	default:
		return "???"
	}
}
