// Package lag provides bare-bones logging functionality
// for the ANWORK project.
//
// The simple usage is as follows.
//   l := lag.New(os.Stdout)
//   l.P(lag.I, "here is a formatted string %s", "foo")
package lag

import (
	"fmt"
	"io"
	"time"
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

	Default = E // The default log level is E
)

// L is a log. You can print formatted stuff to it with P().
type L struct {
	Level // The logging Level at which this L will operate.

	writer io.Writer
}

// New creates an L that writes to the provided writer.
// It is initialized with the Default log Level.
//
// It is totally not thread-safe.
func New(writer io.Writer) *L {
	return &L{writer: writer, Level: Default}
}

// P a printf-formatted message to the io.Writer that this L was
// initialized with. The message will only be printed if this L's
// current Level is greater than or equal to the provided Level.
func (l *L) P(level Level, message string, stuff ...interface{}) {
	if l.Level >= level {
		msg := fmt.Sprintf(message, stuff...)
		fmt.Fprintf(l.writer, "[%s] (%s) %s\n", formatDate(), formatLevel(level), msg)
	}
}

func formatDate() string {
	return time.Now().Format(time.Stamp)
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
