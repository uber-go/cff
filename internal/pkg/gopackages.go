// Package pkg defines the interface for loading Go packages.
//
// It provides a pluggable means for systems like Bazel and Buck
// to inject an alternative package loading mechanism instead of go/packages.
package pkg

import (
	"errors"
	"fmt"
	"go/token"
	"sort"
	"strings"

	"go.uber.org/cff/internal/flag"
	"go.uber.org/multierr"
	"golang.org/x/tools/go/packages"
)

// GoPackagesLoaderFactory builds a Loader
// that uses the go/packages to load package information.
type GoPackagesLoaderFactory struct {
	dir string // used for testing
}

var _ LoaderFactory = (*GoPackagesLoaderFactory)(nil)

// RegisterFlags registers no new flags for GoPackagesLoaderFactory.
func (f *GoPackagesLoaderFactory) RegisterFlags(fset *flag.Set) Loader {
	loader := goPackagesLoader{
		dir:  f.dir,
		load: packages.Load,
	}

	fset.Var(flag.AsList(&loader.tags), "tags",
		"Build tags to load packages with in addition to the 'cff' tag.\n"+
			"This flag may be provided multiple times.")

	return &loader
}

type goPackagesLoader struct {
	tags []flag.String
	dir  string

	// load is a reference to the packages.Load function.
	// By putting it in a function reference,
	// we can easily swap it out for unit tests.
	load func(*packages.Config, ...string) ([]*packages.Package, error)
}

var _ Loader = (*goPackagesLoader)(nil)

const _goPackagesLoadMode = packages.NeedName |
	packages.NeedFiles |
	packages.NeedCompiledGoFiles |
	packages.NeedImports |
	packages.NeedDeps |
	packages.NeedTypes |
	packages.NeedSyntax |
	packages.NeedTypesInfo |
	packages.NeedTypesSizes

func (l *goPackagesLoader) Load(fset *token.FileSet, importPath string) ([]*Package, error) {
	tags := make(map[string]struct{}, len(l.tags)+1)
	tags["cff"] = struct{}{}
	for _, tag := range l.tags {
		// For convenience, split on "," since that's what "go build -tags"
		// expects.
		for _, tag := range strings.Split(string(tag), ",") {
			tags[tag] = struct{}{}
		}
	}

	uniqueTags := make([]string, 0, len(tags))
	for tag := range tags {
		uniqueTags = append(uniqueTags, tag)
	}
	sort.Strings(uniqueTags)

	pkgs, err := l.load(&packages.Config{
		Mode:       _goPackagesLoadMode,
		Fset:       fset,
		BuildFlags: []string{"-tags", strings.Join(uniqueTags, ",")},
		Dir:        l.dir,
		Tests:      true,
	}, importPath)
	if err != nil {
		return nil, fmt.Errorf("packages.Load: %w", err)
	}
	if len(pkgs) == 0 {
		return nil, errors.New("no packages found")
	}

	ipkgs := make([]*Package, len(pkgs))
	for i, pkg := range pkgs {
		// pkg.Errors is a []packages.Error so we can't
		// use mutlierr.Combine.
		for _, e := range pkg.Errors {
			err = multierr.Append(err, e)
		}
		ipkgs[i] = &Package{
			CompiledGoFiles: pkg.CompiledGoFiles,
			Syntax:          pkg.Syntax,
			Types:           pkg.Types,
			TypesInfo:       pkg.TypesInfo,
		}
	}

	return ipkgs, err
}
