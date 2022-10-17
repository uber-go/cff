package pkg

import (
	"fmt"
	"go/token"
	"strings"

	"go.uber.org/cff/internal/flag"
	"code.uber.internal/go/importer"
)

// ArchiveLoaderFactory builds a Loader
// that loads package information for a package
// from source files for that package,
// and pre-compiled archive data for its dependencies.
// This is the optimal way of loading package data in a Bazel environment.
type ArchiveLoaderFactory struct{}

var _ LoaderFactory = (*ArchiveLoaderFactory)(nil)

// RegisterFlags registers new flags needed to get archive data.
func (f *ArchiveLoaderFactory) RegisterFlags(fset *flag.Set) Loader {
	loader := archiveLoader{
		loadArchive: importer.LoadArchive,
	}

	fset.Var(flag.AsList(&loader.archives), "archive", "A value in the form 'IMPORTPATHS=IMPORTMAP=FILE=EXPORT' where,\n"+
		"  - IMPORTPATHS is a list of colon-separated import paths\n"+
		"  - IMPORTMAP is a package import path\n"+
		"  - FILE is the path to the archive file for the package at IMPORTMAP\n"+
		"  - EXPORT is the path to the export file for the package at IMPORTMAP\n"+
		"Pass this zero or more times to specify archives for dependencies of a package.")

	fset.Var(flag.AsList(&loader.srcs), "source",
		"Path to a Go source file for the package cff is generating code for.\n"+
			"This flag may be passed multiple times.")

	fset.Var(&loader.stdlibRoot, "stdlibroot",
		"Specifies the path containing archive files for the Go standard library.")

	return &loader
}

type archiveLoader struct {
	// These pointers will be filled when the flags are parsed.
	archives   []archiveValue
	srcs       []flag.String
	stdlibRoot flag.String

	// loadArchive is a reference to importer.LoadArchive
	// that we can swap out for tests.
	loadArchive func(importer.LoadParams) (*importer.Package, error)
}

var _ Loader = (*archiveLoader)(nil)

func (l *archiveLoader) Load(fset *token.FileSet, importPath string) ([]*Package, error) {
	archives := make([]importer.Archive, len(l.archives))
	for i, a := range l.archives {
		// We're recording all fields, but go/importer ignores
		// everything except ImportMap and File right now.
		archives[i] = importer.Archive{
			ImportPaths: a.ImportPaths,
			ImportMap:   a.ImportMap,
			File:        a.File,
			ExportFile:  a.Export,
		}
	}

	srcs := make([]string, len(l.srcs))
	for i, s := range l.srcs {
		srcs[i] = string(s)
	}

	pkg, err := l.loadArchive(importer.LoadParams{
		Fset:       fset,
		ImportPath: importPath,
		Srcs:       srcs,
		StdlibRoot: string(l.stdlibRoot),
		Archives:   archives,
	})
	if err != nil {
		return nil, err
	}
	return []*Package{
		{
			CompiledGoFiles: pkg.CompiledGoFiles,
			Syntax:          pkg.Syntax,
			Types:           pkg.Types,
			TypesInfo:       pkg.TypesInfo,
		},
	}, nil
}

// archiveValue is a Go package archive passed as a command line argument.
//
// The following is the flag format:
//
//	--archive=IMPORTPATHS=IMPORTMAP=FILE=EXPORT
//
// For example,
//
//	--archive=github.com/foo/bar:github.com/foo/baz=github.com/foo/bar=bar.go=bar_export.go
//
// The flag is structured in this format to closely follow
// https://github.com/bazelbuild/rules_go/blob/8ea79bbd5e6ea09dc611c245d1dc09ef7ab7118a/go/private/actions/compile.bzl#L20.
type archiveValue struct {
	ImportPaths []string
	ImportMap   string
	File        string
	Export      string
}

var _ flag.Getter = (*archiveValue)(nil)

func (a *archiveValue) String() string {
	importPaths := strings.Join(a.ImportPaths, ":")
	return fmt.Sprintf("%v=%v=%v=%v", importPaths, a.ImportMap, a.File, a.Export)
}

// Set receives a flag value.
func (a *archiveValue) Set(name string) error {
	args := strings.Split(name, "=")
	if len(args) != 4 {
		return fmt.Errorf("expected 4 elements, got %d", len(args))
	}

	*a = archiveValue{
		ImportPaths: strings.Split(args[0], ":"),
		ImportMap:   args[1],
		File:        args[2],
		Export:      args[3],
	}
	return nil
}

// Get returns the current value of the archiveValue.
func (a *archiveValue) Get() any { return a }
