package internal

import (
	"go/token"

	"go.uber.org/multierr"
	"golang.org/x/tools/go/packages"
)

// Process processes the provided Go package with cff.
func Process(fset *token.FileSet, pkg *packages.Package, outputPath string) error {
	var errors error
	for _, e := range pkg.Errors {
		errors = multierr.Append(errors, e)
	}
	if errors != nil {
		return errors
	}

	c := newCompiler(fset, pkg.TypesInfo, pkg.Types)

	var files []*file
	for _, file := range pkg.Syntax {
		f, err := c.CompileFile(file)
		if err != nil {
			continue
		}
		files = append(files, f)
	}

	for _, err := range c.errors {
		errors = multierr.Append(errors, err)
	}

	g := newGenerator(fset, outputPath)
	for _, f := range files {
		if err := g.GenerateFile(f); err != nil {
			errors = multierr.Append(errors, err)
		}
	}

	return errors
}
