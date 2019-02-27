package internal

import (
	"go/token"

	"go.uber.org/multierr"
	"golang.org/x/tools/go/packages"
)

// Process processes the provided Go package with cff.
func Process(fset *token.FileSet, pkg *packages.Package, outputDir string) error {
	var err error
	for _, e := range pkg.Errors {
		err = multierr.Append(err, e)
	}
	if err != nil {
		return err
	}

	c := newCompiler(fset, pkg.TypesInfo, pkg.Types)

	var files []*file
	for _, file := range pkg.Syntax {
		f, err := c.CompileFile(file)
		if err != nil {
			return err
		}
		files = append(files, f)
	}

	g := newGenerator(fset, outputDir)
	for _, f := range files {
		if err := g.GenerateFile(f); err != nil {
			return err
		}

	}

	return nil
}
