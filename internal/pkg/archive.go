package pkg

import (
	"fmt"
	"go/token"
	"strings"

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
func (f *ArchiveLoaderFactory) RegisterFlags(parser Command) (Loader, error) {
	var opts struct {
		// Archives is a list of archive paths for
		// dependencies of the package being loaded.
		Archives []string `long:"archive" value-name:"IMPORTPATHS=IMPORTMAP=FILE=EXPORT"`
		// Sources is a list of source files for the package.
		Sources []string `long:"source"`
		// StdlibRoot specifies where archive files for
		// the standard library are stored.
		StdlibRoot string `long:"stdlibroot"`
	}
	if _, err := parser.AddGroup("Archive Data", "", &opts); err != nil {
		return nil, err
	}

	parser.FindOptionByLongName("archive").Description =
		"Use the given archive FILE for import path IMPORTMAP when parsing the " +
			"source files. IMPORTPATHS is a colon-separated list of import paths; " +
			"IMPORTMAP is the actual import path of the library this archive " +
			"holds; FILE is the path to the archive file; EXPORT is the path to " +
			"the corresponding export file. Currently, IMPORTPATHS and EXPORT " +
			"arguments are ignored."
	parser.FindOptionByLongName("source").Description =
		"When using archives to parse the source code, specifies the filepaths to " +
			"all Go code in the package, so that CFF can parse the entire " +
			"package."
	parser.FindOptionByLongName("stdlibroot").Description =
		"When using archives to parse the source code, specifies the path containing " +
			"archive files for the Go standard library."

	return &archiveLoader{
		archives:    &opts.Archives,
		srcs:        &opts.Sources,
		stdlibRoot:  &opts.StdlibRoot,
		loadArchive: importer.LoadArchive,
	}, nil
}

type archiveLoader struct {
	// These pointers will be filled when the flags are parsed.
	archives   *[]string
	srcs       *[]string
	stdlibRoot *string

	// loadArchive is a reference to importer.LoadArchive
	// that we can swap out for tests.
	loadArchive func(importer.LoadParams) (*importer.Package, error)
}

var _ Loader = (*archiveLoader)(nil)

func (l *archiveLoader) Load(fset *token.FileSet, importPath string) ([]*Package, error) {
	archives := make([]importer.Archive, len(*l.archives))
	for i, archive := range *l.archives {
		// TODO(abg): This should happen as part of flag parsing.
		a, err := parseArchive(archive)
		if err != nil {
			return nil, fmt.Errorf("invalid argument --archive=%q: %w", archive, err)
		}
		archives[i] = a
	}

	pkg, err := l.loadArchive(importer.LoadParams{
		Fset:       fset,
		ImportPath: importPath,
		Srcs:       *l.srcs,
		StdlibRoot: *l.stdlibRoot,
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

// parseArchive parses the archive string to the internal.Archive type.
//
// The following is the flag format:
//
//	--archive=IMPORTPATHS=IMPORTMAP=FILE=EXPORT
//
// For example,
//
//	--archive=github.com/foo/bar:github.com/foo/baz=github.com/foo/bar=bar.go=bar_export.go
//
// The flag is structured in this format to closely follow https://github.com/bazelbuild/rules_go/blob/8ea79bbd5e6ea09dc611c245d1dc09ef7ab7118a/go/private/actions/compile.bzl#L20;
// however, the IMPORTPATHS and EXPORT elements are ignored. There may be future
// work involved in resolving import aliases, using IMPORTPATHS.
func parseArchive(archive string) (importer.Archive, error) {
	args := strings.Split(archive, "=")
	if len(args) != 4 {
		return importer.Archive{}, fmt.Errorf("expected 4 elements, got %d", len(args))
	}

	// Currently, we ignore the IMPORTPATHS and EXPORT elements.
	return importer.Archive{
		ImportMap: args[1],
		File:      args[2],
	}, nil
}
