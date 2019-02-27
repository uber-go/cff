package main

import (
	"errors"
	"fmt"
	"go/token"
	"log"
	"os"

	"go.uber.org/cff/internal"
	"go.uber.org/multierr"
	"golang.org/x/tools/go/packages"
)

func main() {
	log.SetFlags(0)
	if err := run(os.Args[1:]); err != nil {
		log.Fatalf("%+v", err)
	}
}

func run(args []string) error {
	fset := token.NewFileSet()
	pkgs, err := packages.Load(&packages.Config{
		Mode:       packages.LoadSyntax,
		Fset:       fset,
		BuildFlags: []string{"-tags=cff"},
	}, args...)

	if err != nil {
		return fmt.Errorf("could not load packages: %v", err)
	}

	if len(pkgs) == 0 {
		return errors.New("no packages found")
	}

	for _, pkg := range pkgs {
		err = multierr.Append(err, internal.Process(fset, pkg, ""))
	}
	return err
}
