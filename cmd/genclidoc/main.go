// This is a very simple application that regenerates the command line interface documentation
// that is used in the anwork (see github.com/ankeesler/anwork/cmd/anwork) application. To be
// specific, this executable generates documentation in Github markdown format.
//
// To run, simply run the binary and pass the output file as an argument.
//
// Usage: genclidoc <output-file>
package main

import (
	"fmt"
	"io"
	"os"

	"github.com/ankeesler/anwork/runner"
)

func gendoc(output io.Writer) {
	fmt.Fprintln(output, "Generated by genclidoc. DO NOT EDIT.")
	fmt.Fprintln(output)

	fmt.Fprintln(output, "# _anwork_ CLI commands, version", runner.Version)
	fmt.Fprintln(output)

	runner.MarkdownUsage(output)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: genclidoc <output-file>")
		os.Exit(1)
	}

	output := os.Args[1]
	fmt.Printf("Generating documentation into %s\n", output)

	outputFile, err := os.Create(output)
	defer outputFile.Close()
	if err != nil {
		fmt.Printf("Error! %s\n", err.Error())
		os.Exit(1)
	}

	gendoc(outputFile)

	fmt.Println("Done.")
}
