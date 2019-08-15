package internal

import (
	"go/ast"
	"go/token"

	"go.uber.org/multierr"
	"golang.org/x/tools/go/packages"
)

// Process processes the provided Go package with cff.
func Process(fset *token.FileSet, pkg *packages.Package, file *ast.File, outputPath string, compilerOpts CompilerOpts) error {
	var err error
	for _, e := range pkg.Errors {
		err = multierr.Append(err, e)
	}
	if err != nil {
		return err
	}

	c := newCompiler(fset, pkg.TypesInfo, pkg.Types, compilerOpts)

	f, err := c.CompileFile(file)
	if err != nil {
		return err
	}

	g := newGenerator(fset, outputPath)
	if err := g.GenerateFile(f); err != nil {
		return err
	}

	return nil
}
