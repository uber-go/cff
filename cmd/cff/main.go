package main

import (
	"errors"
	"fmt"
	"go/token"
	"log"
	"os"

	"go.uber.org/cff/internal"
	flags "github.com/jessevdk/go-flags"
	"go.uber.org/multierr"
	"golang.org/x/tools/go/packages"
)

type options struct {
	Input  string `long:"input" required:"yes"`
	Output string `long:"output" required:"yes"`
}

func newCLIParser() (*flags.Parser, *options) {
	var opts options
	parser := flags.NewParser(&opts, flags.HelpFlag)
	parser.Name = "cff"

	// This is more readable than embedding the descriptions in the options
	// above.
	parser.FindOptionByLongName("input").Description =
		"Pattern to search for the CFF source code."
	parser.FindOptionByLongName("output").Description =
		"Path to which the output file should be generated."

	return parser, &opts
}

func main() {
	if err := run(os.Args); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run(args []string) error {
	parser, f := newCLIParser()
	if _, err := parser.ParseArgs(args); err != nil {
		return err
	}

	fset := token.NewFileSet()
	pkgs, err := packages.Load(&packages.Config{
		Mode:       packages.LoadSyntax,
		Fset:       fset,
		BuildFlags: []string{"-tags=cff"},
	}, f.Input)

	if err != nil {
		return fmt.Errorf("could not load packages: %v", err)
	}

	if len(pkgs) == 0 {
		return errors.New("no packages found")
	}

	for _, pkg := range pkgs {
		err = multierr.Append(err, internal.Process(fset, pkg, f.Output))
	}
	return err
}
