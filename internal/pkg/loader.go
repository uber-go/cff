package pkg

import (
	"go/ast"
	"go/token"
	"go/types"

	flags "github.com/jessevdk/go-flags"
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

// Command is a CLI parser command that we can register new options with.
// This conforms to the go-flags.Command struct's API.
type Command interface {
	AddGroup(short, long string, data any) (*flags.Group, error)
	FindOptionByLongName(longName string) (option *flags.Option)
}

var _ Command = (*flags.Command)(nil)

// LoaderFactory builds Loaders from command line flags.
type LoaderFactory interface {
	// RegisterFlags builds a Loader and registers the flags
	// necessary for it with the given Command.
	//
	// The returned Loader will not be used until the parser
	// has finished parsing its arguments.
	RegisterFlags(Command) (Loader, error)
}

// Loader loads information about a Go package from its import path.
type Loader interface {
	Load(fset *token.FileSet, importPath string) ([]*Package, error)
}
