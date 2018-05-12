// A Runner is an object that can run the various pieces of anwork functionality, e.g.,
// create tasks, show tasks, print out version information, etc.
package runner

import (
	"fmt"
	"io"
	"strings"

	"github.com/ankeesler/anwork/task"
)

// Print the usage of every anwork runner command to the provided output writer.
func Usage(output io.Writer) {
	for _, c := range commands {
		fmt.Fprintf(output, "  %s %s\n", c.Name, strings.Join(c.Args, " "))
		fmt.Fprintf(output, "        %s\n", c.Description)
	}
}

// A Runner is an object that can run the various pieces of anwork functionality, e.g.,
// create tasks, show tasks, print out version information, etc.
type Runner struct {
	factory                   task.ManagerFactory
	stdoutWriter, debugWriter io.Writer
}

// Create a new Runner. The task.ManagerFactory argument will be used to create a
// manager for use by the Runner. The Runner will write its regular output to
// the stdoutWriter and its debug output to the debugWriter.
func New(factory task.ManagerFactory, stdoutWriter, debugWriter io.Writer) *Runner {
	return &Runner{
		factory:      factory,
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
			args[0], args[1:], cmd.Args)
	}

	manager, err := a.factory.Create()
	if err != nil {
		return fmt.Errorf("Could not create manager: %s", err.Error())
	}
	a.debug("Manager is %s", manager)

	if err := cmd.Action(cmd, args, a.stdoutWriter, manager); err != nil {
		if _, ok := err.(*resetError); ok {
			if err = a.factory.Reset(); err != nil {
				return fmt.Errorf("Could not reset factory: %s", err.Error())
			}
		} else {
			return fmt.Errorf("Command '%s' failed: %s", args[0], err.Error())
		}
	}

	if err := a.factory.Save(manager); err != nil {
		return fmt.Errorf("Could not save manager: %s", err.Error())
	}

	return nil
}

func (a *Runner) debug(format string, args ...interface{}) {
	fmt.Fprintln(a.debugWriter, format, args)
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
