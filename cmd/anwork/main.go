// This is the main anwork command line executable. This command line executable provides the
// ability to create, read, update, and delete anwork Task objects.
//
// Versioning is done with a single 32-bit integer. Version names start with a lowercase 'v' and are
// then followed by the number of the release. For example, the first version of the release was
// named _v1_. The second version of the release will be _v2_. There are no minor version
// numbers. This version number is controlled via the "version" global in the runner package. See
// the CLI command "anwork version" for more information.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ankeesler/anwork/runner"
	"github.com/ankeesler/anwork/task/local"
)

type debugWriter struct {
	debug bool
}

func (dw *debugWriter) Write(data []byte) (int, error) {
	if dw.debug {
		return fmt.Print(string(data))
	}
	return 0, nil
}

func usage(r *runner.Runner, f *flag.FlagSet) func() {
	return func() {
		fmt.Println("Usage of anwork")
		fmt.Println("Flags")
		f.PrintDefaults()
		fmt.Println("Commands")
		r.Usage(os.Stdout)
	}
}

func main() {
	var (
		context, root string
		dw            debugWriter
	)

	flags := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	flags.BoolVar(&dw.debug, "d", false, "Enable debug printing")

	flags.StringVar(&context, "c", "default-context", "Set the persistence context")
	flags.StringVar(&root, "o", ".", "Set the persistence root directory")

	factory := local.NewManagerFactory(root, context)
	r := runner.New(factory, os.Stdout, &dw)
	flags.Usage = usage(r, flags)

	if err := flags.Parse(os.Args[1:]); err == flag.ErrHelp {
		// Looks like help is printed by the flag package...
		os.Exit(0)
	} else if err != nil {
		// I think the flag package prints out the error and the usage...
		os.Exit(1)
	}

	if flags.NArg() == 0 {
		// If there are no arguments, return success. People might use this to simply check if the anwork
		// executable is on their machine.
		flags.Usage()
		os.Exit(0)
	}

	if err := r.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}
}
