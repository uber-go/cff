package pkg

import (
	"errors"
	"fmt"
	"go/token"

	"go.uber.org/multierr"
	"golang.org/x/tools/go/packages"
)

// GoPackagesLoaderFactory builds a Loader
// that uses the go/packages to load package information.
type GoPackagesLoaderFactory struct {
	// BuildFlags is a list of flags to load packages with.
	//
	// For example, the following will enable the "cff" build tag
	// when loading packages.
	//
	//	BuildFlags: []string{"-tags", "cff"},
	BuildFlags []string

	dir string // used for testing
}

var _ LoaderFactory = (*GoPackagesLoaderFactory)(nil)

// RegisterFlags registers no new flags for GoPackagesLoaderFactory.
func (f *GoPackagesLoaderFactory) RegisterFlags(cmd Command) (Loader, error) {
	return &goPackagesLoader{
		buildFlags: f.BuildFlags,
		dir:        f.dir,
		load:       packages.Load,
	}, nil
}

type goPackagesLoader struct {
	buildFlags []string
	dir        string

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
	pkgs, err := l.load(&packages.Config{
		Mode:       _goPackagesLoadMode,
		Fset:       fset,
		BuildFlags: l.buildFlags,
		Dir:        l.dir,
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
