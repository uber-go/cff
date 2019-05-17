package main

import (
	"errors"
	"flag"
	"fmt"
	"go/token"
	"log"
	"os"

	"go.uber.org/cff/internal"
	"go.uber.org/multierr"
	"golang.org/x/tools/go/packages"
)

type flags struct {
	// Input is the path to CFF2 source code.
	Input string
	// Output represents the path at which the generated code with be deposited.
	Output string
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run() error {
	fs := flag.NewFlagSet("cff", flag.PanicOnError)
	var f flags
	fs.StringVar(&f.Input, "input", "", "Path for CFF2 source.")
	fs.StringVar(&f.Output, "output", "", "Output file path for generated code.")
	if err := fs.Parse(os.Args[1:]); err != nil {
		return err
	}
	if f.Input == "" {
		return fmt.Errorf("must specify an input")
	}
	if f.Output == "" {
		return fmt.Errorf("must specify output path")
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
