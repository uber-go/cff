package internal

import (
	"go/ast"
	"go/token"
)

// Process processes the provided Go package with cff.
func Process(fset *token.FileSet, pkg *Package, file *ast.File, outputPath string, compilerOpts CompilerOpts) error {
	c := newCompiler(fset, pkg.TypesInfo, pkg.Types, compilerOpts)

	f, err := c.CompileFile(file, pkg)
	if err != nil {
		return err
	}

	g := newGenerator(fset, outputPath)
	if err := g.GenerateFile(f); err != nil {
		return err
	}

	return nil
}
