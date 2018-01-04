package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ankeesler/anwork/storage"
	"github.com/ankeesler/anwork/task"
)

// This is the version of this anwork app.
const version = 1

// These variables are used to store command line flag values.
var (
	debug         bool
	context, root string
)

func dbgfln(output io.Writer, format string, stuff ...interface{}) {
	if debug {
		fmt.Fprint(output, "anwork: dbg: ")
		fmt.Fprintf(output, format, stuff...)
		fmt.Fprintln(output)
	}
}

func readManager(persister *storage.Persister, context string, manager *task.Manager) {
	err := persister.Unpersist(context, manager)
	if err != nil {
		panic(err.Error())
	}
}

func writeManager(persister *storage.Persister, context string, manager *task.Manager) {
	err := persister.Persist(context, manager)
	if err != nil {
		panic(err.Error())
	}
}

func run(args []string, output io.Writer) int {
	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.SetOutput(output)

	flags.BoolVar(&debug, "debug", false, "Enable debug printing")

	flags.StringVar(&context, "context", "default-context", "Set the persistence context")
	flags.StringVar(&root, "root", ".", "Set the persistence root directory")

	if err := flags.Parse(args[1:]); err == flag.ErrHelp {
		// Looks like help is printed by the flag package...
		return 0
	} else if err != nil {
		// I think the flag package prints out the error and the usage...
		return 1
	}

	if flags.NArg() == 0 {
		fmt.Fprintln(output, "Error! Expected command arguments")
		flags.Usage()
		return 1
	}
	firstArg := flags.Arg(0)

	persister := storage.Persister{root}
	manager := task.NewManager()
	if persister.Exists(context) {
		dbgfln(output, "Context %s exists at root %s", context, root)
		readManager(&persister, context, manager)
	} else {
		dbgfln(output, "Context %s does not exist at root %s; creating it", context, root)
		writeManager(&persister, context, manager)
	}
	dbgfln(output, "Manager is %s", manager)

	command := findCommand(firstArg)
	if command == nil {
		fmt.Fprintln(output, "Error! Unknown command line argument:", firstArg)
		flags.Usage()
		return 1
	} else {
		if command.action(firstArg, output, manager) {
			dbgfln(output, "Persisting manager back to disk")
			writeManager(&persister, context, manager)
		} else {
			dbgfln(output, "NOT persisting manager back to disk")
		}
	}

	return 0
}

func main() {
	os.Exit(run(os.Args, os.Stdout))
}
