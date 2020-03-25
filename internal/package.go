package internal

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/packages"
)

// Package is the scope of *packages.Package that is needed for cff.
type Package struct {
	CompiledGoFiles []string
	Syntax          []*ast.File
	Types           *types.Package
	TypesInfo       *types.Info
}

// NewPackage creates a internal.Package from a packages.Package.
func NewPackage(pkg *packages.Package) *Package {
	return &Package{
		CompiledGoFiles: pkg.CompiledGoFiles,
		Syntax:          pkg.Syntax,
		Types:           pkg.Types,
		TypesInfo:       pkg.TypesInfo,
	}
}

// Archive holds information about a library and its corresponding archive.
//
// It is mainly used as a flag through the cff CLI. Its format closely follows
// the format used in Bazel rules_go: https://github.com/bazelbuild/rules_go/blob/8ea79bbd5e6ea09dc611c245d1dc09ef7ab7118a/go/private/actions/compile.bzl#L20
//
// The following is the flag format:
//
//  --archive=IMPORTPATHS=IMPORTMAP=FILE=EXPORT
//
// For example,
//
//  --archive=github.com/foo/bar:github.com/foo/baz=github.com/foo/bar=bar.go=bar_export.go
//
// However, we only use the ImportMap and File attribute from this flag. In the
// future, we may use ImportPaths to resolve import aliases.
type Archive struct {
	// ImportMap refers to the actual import path to the library this archive
	// represents. While the naming may be confusing, this closely follows Bazel
	// rules_go conventions.
	//
	// See https://github.com/bazelbuild/rules_go/blob/f7a8cb6b9158006e5dfc91074f9636820a446921/go/core.rst#go_library.
	ImportMap string
	File      string
}

// LoadParams is the arguments needed for loading the package for the CFF
// compiler.
type LoadParams struct {
	Fset       *token.FileSet
	ImportPath string
	Srcs       []string
	StdlibRoot string
	Archives   []Archive
}

// PackagesArchive reads Go archive information for any dependencies to build
// the package to be processed by the CFF compiler. This is a solution for
// Bazel-driven CFF, in which virtualizing GOPATH is expensive but required for
// using the go/packages API.
func PackagesArchive(p LoadParams) ([]*Package, error) {
	// Parse the CFF files and its package contents by hand, and build the
	// necessary components for a package to be processed by the CFF compiler.
	files := make([]*ast.File, 0, len(p.Srcs))
	for _, src := range p.Srcs {
		f, err := parser.ParseFile(p.Fset, src, nil, parser.ParseComments)
		if err != nil {
			return nil, fmt.Errorf("error parsing %q: %v", src, err)
		}
		files = append(files, f)
	}

	// Build an importer using the imports map built by reading dependency
	// archives, and use it to build the *types.Package and *types.PosInfo for the
	// source files.
	imp, err := newImporter(p.Fset, p.Archives, p.StdlibRoot)
	if err != nil {
		return nil, err
	}
	conf := types.Config{Importer: imp}
	typesInfo := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Implicits:  make(map[ast.Node]types.Object),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
		Scopes:     make(map[ast.Node]*types.Scope),
	}
	pkg, err := conf.Check(p.ImportPath, p.Fset, files, typesInfo)
	if err != nil {
		return nil, fmt.Errorf("error building pkg %q: %v", p.ImportPath, err)
	}
	return []*Package{
		{
			CompiledGoFiles: p.Srcs,
			Syntax:          files,
			Types:           pkg,
			TypesInfo:       typesInfo,
		},
	}, nil
}
