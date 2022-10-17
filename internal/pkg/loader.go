package pkg

import (
	"go/ast"
	"go/token"
	"go/types"

	"go.uber.org/cff/internal/flag"
)

// Package is a Go package that cff is going to generate code for.
type Package struct {
	// CompiledGoFiles is a list of absolute file paths of Go files.
	CompiledGoFiles []string
	// Syntax is a ASTs of files in CompiledGoFiles in the same order.
	Syntax []*ast.File
	// Types is type information about CompiledGoFiles.
	Types *types.Package
	// TypesInfo provides type information about ASTs.
	TypesInfo *types.Info
}

// LoaderFactory builds Loaders from command line flags.
type LoaderFactory interface {
	// RegisterFlags builds a Loader and registers the flags
	// necessary for it with the given flag set.
	//
	// The returned Loader will not be used until the parser
	// has finished parsing its arguments.
	RegisterFlags(*flag.Set) Loader
}

// Loader loads information about a Go package from its import path.
type Loader interface {
	Load(fset *token.FileSet, importPath string) ([]*Package, error)
}
