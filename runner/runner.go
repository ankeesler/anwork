// Package runner contains the code for driving the anwork command line interface.
//
// A Runner is an object that can run the various pieces of anwork functionality, e.g.,
// create tasks, show tasks, print out version information, etc.
package runner

import (
	"fmt"
	"io"
	"strings"

	"github.com/ankeesler/anwork/manager"
)

// Print the usage of every anwork runner command to the provided output writer.
func Usage(output io.Writer) {
	for _, c := range commands {
		fmt.Fprintf(output, "  %s %s\n", c.Name, strings.Join(c.Args, " "))
		fmt.Fprintf(output, "        %s", c.Description)
		if c.Alias != "" {
			fmt.Fprintf(output, " (alias: %s)", c.Alias)
		}
		fmt.Fprintln(output)
	}
}

// Print the usage of every anwork runner command in Github markdown format
// to the provided output writer.
func MarkdownUsage(output io.Writer) {
	for _, c := range commands {
		fmt.Fprintf(output, "### `anwork %s", c.Name)
		for _, a := range c.Args {
			fmt.Fprintf(output, " %s", a)
		}
		fmt.Fprintln(output, "`")

		fmt.Fprintf(output, "* %s\n", c.Description)

		if c.Alias != "" {
			fmt.Fprintf(output, "* Alias: `%s`\n", c.Alias)
		}
	}
}

// This is a struct used to communicate the git hash and date at which this package was built.
type BuildInfo struct {
	Hash, Date string
}

// A Runner is an object that can run the various pieces of anwork functionality, e.g.,
// create tasks, show tasks, print out version information, etc.
type Runner struct {
	buildInfo                 *BuildInfo
	manager                   manager.Manager
	stdoutWriter, debugWriter io.Writer
}

// New creates a new Runner. The manager.Manager will be used to perform the task
// operations. The Runner will write its regular output to the stdoutWriter and its
// debug output to the debugWriter.
func New(buildInfo *BuildInfo, manager manager.Manager, stdoutWriter, debugWriter io.Writer) *Runner {
	return &Runner{
		buildInfo:    buildInfo,
		manager:      manager,
		stdoutWriter: stdoutWriter,
		debugWriter:  debugWriter,
	}
}

// Run the functionality specified via the arguments. The Runner will parse the args
// and run the appropriate piece of functionality. See runner.Runner.Usage() for
// a print out of the usage of this Runner.
func (a *Runner) Run(args []string) error {
	cmd := findCommand(args[0])
	if cmd == nil {
		return fmt.Errorf("Unknown command: '%s'", args[0])
	}

	if !validateArgs(cmd, args) {
		return fmt.Errorf("Invalid argument passed to command '%s':\n\tGot: %s\n\tExpected: %s",
			cmd.Name, args[1:], cmd.Args)
	}

	a.debug("Manager is %s\n", a.manager)

	if err := cmd.Action(cmd, args, a.stdoutWriter, a.manager, a.buildInfo); err != nil {
		return fmt.Errorf("Command '%s' failed: %s", args[0], err.Error())
	}

	return nil
}

func (a *Runner) debug(format string, args ...interface{}) {
	fmt.Fprintf(a.debugWriter, format, args...)
}

func validateArgs(cmd *command, args []string) bool {
	if len(cmd.Args) == (len(args) - 1) {
		return true
	}

	// optional argument (e.g., [task-name])
	if (len(cmd.Args) - 1) == (len(args) - 1) {
		lastArg := cmd.Args[len(cmd.Args)-1]
		if strings.HasPrefix(lastArg, "[") && strings.HasSuffix(lastArg, "]") {
			return true
		}
	}

	return false
}
