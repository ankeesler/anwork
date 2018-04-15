// This is the main anwork command line executable. This command line executable provides the
// ability to create, read, update, and delete anwork Task objects.
//
// Versioning is done with a single 32-bit integer. Version names start with a lowercase 'v' and are
// then followed by the number of the release. For example, the first version of the release was
// named _v1_. The second version of the release will be _v2_. There are no minor version
// numbers. This version number is controlled via the "version" property in the command package. See
// the CLI command "anwork version" for more information.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/ankeesler/anwork/cmd/anwork/command"
	"github.com/ankeesler/anwork/task/local"
)

// These variables are used to store command line flag values.
var (
	debug bool
)

func dbgfln(output io.Writer, format string, stuff ...interface{}) {
	if debug {
		fmt.Fprint(output, "anwork: dbg: ")
		fmt.Fprintf(output, format, stuff...)
		fmt.Fprintln(output)
	}
}

func usage(f *flag.FlagSet, output io.Writer) func() {
	return func() {
		fmt.Fprintln(output, "Usage of anwork")
		fmt.Fprintln(output, "Flags")
		f.PrintDefaults()
		fmt.Fprintln(output, "Commands")
		for _, c := range command.Commands {
			c.Usage(output)
		}
	}
}

func run(args []string, output io.Writer) int {
	var (
		context, root string
	)

	flags := flag.NewFlagSet(args[0], flag.ContinueOnError)
	flags.SetOutput(output)

	flags.BoolVar(&debug, "d", false, "Enable debug printing")

	flags.StringVar(&context, "c", "default-context", "Set the persistence context")
	flags.StringVar(&root, "o", ".", "Set the persistence root directory")

	flags.Usage = usage(flags, output)

	if err := flags.Parse(args[1:]); err == flag.ErrHelp {
		// Looks like help is printed by the flag package...
		return 0
	} else if err != nil {
		// I think the flag package prints out the error and the usage...
		return 1
	}

	if flags.NArg() == 0 {
		// If there are no arguments, return success. People might use this to simply check if the anwork
		// executable is on their machine.
		flags.Usage()
		return 0
	}
	firstArg := flags.Arg(0)

	factory := local.NewManagerFactory(root, context)
	manager, err := factory.Create()
	if err != nil {
		fmt.Fprintf(output, "Error! Could not create manager: %s\n", err.Error())
		return 1
	}
	dbgfln(output, "Manager is %s", manager)

	cmd := command.FindCommand(firstArg)
	if cmd == nil {
		fmt.Fprintln(output, "Error! Unknown command:", firstArg)
		flags.Usage()
		return 1
	} else {
		switch cmd.Run(flags, output, manager) {
		case command.ResponsePersist:
			dbgfln(output, "Persisting manager back to disk")
			if err := factory.Save(manager); err != nil {
				fmt.Fprintf(output, "Error! Could not save manager: %s\n", err.Error())
				return 1
			}
		case command.ResponseNoPersist:
			dbgfln(output, "NOT persisting manager back to disk")
		case command.ResponseReset:
			dbgfln(output, "Completely deleting everything in context %s", context)
			if err := factory.Reset(); err != nil {
				fmt.Fprintf(output, "Error! Could not reset context %s: %s\n", context, err.Error())
				return 1
			}
		case command.ResponseArgumentError:
			fmt.Fprintln(output, "Error! Wrong arguments passed to command")
			cmd.Usage(output)
			return 1
		case command.ResponseFailure:
			return 1
		}
	}

	return 0
}

func main() {
	os.Exit(run(os.Args, os.Stdout))
}
