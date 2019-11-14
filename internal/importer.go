package internal

import (
	"fmt"
	"go/token"
	"go/types"
	"os"
	"strings"

	"go.uber.org/multierr"
	"golang.org/x/tools/go/gcexportdata"
)

type cffImporter struct {
	fset       *token.FileSet
	stdlibRoot string
	imports    map[string]*types.Package
}

func newImporter(fset *token.FileSet, archives []Archive, root string) (types.Importer, error) {
	importer := &cffImporter{
		fset:       fset,
		stdlibRoot: root,
		imports:    make(map[string]*types.Package), // import path -> package
	}
	// archives here should be the archive files for the dependencies (imports)
	// of srcs. We emulate the GOPATH by reading the archives for the dependencies
	// and creating an importer from it.
	for _, archive := range archives {
		if _, err := importer.readArchive(archive.File, archive.ImportMap); err != nil {
			return nil, fmt.Errorf("error reading archive %s: %v", archive, err)
		}
	}
	return importer, nil
}

func (i *cffImporter) Import(path string) (*types.Package, error) {
	// if the import is part of GoStdLib, it should read the archive file for it
	// to load the package.
	if isStdlibImport(path) {
		archiveFile := fmt.Sprintf("%v/%v.a", i.stdlibRoot, path)
		return i.readArchive(archiveFile, path)
	}

	// if the import is not part of GoStdLib, `imports` should have the package
	// loaded already.
	if pkg, ok := i.imports[path]; ok {
		return pkg, nil
	}
	return nil, fmt.Errorf("error finding package %q in importer: please double check dependencies for the cff Bazel rule", path)
}

func (i *cffImporter) readArchive(archiveFile, importPath string) (_ *types.Package, err error) {
	f, err := os.Open(archiveFile)
	if err != nil {
		return nil, fmt.Errorf("error opening archive file: %v", err)
	}
	defer func() {
		err = multierr.Append(err, f.Close())
	}()

	r, err := gcexportdata.NewReader(f)
	if err != nil {
		return nil, fmt.Errorf("error preparing reader for archive file %s: %v", archiveFile, err)
	}
	pkg, err := gcexportdata.Read(r, i.fset, i.imports, importPath)
	if err != nil {
		return nil, fmt.Errorf("error reading archive file: %v", err)
	}
	return pkg, nil
}

// By Go standards from v1.13, the import prefix before the first "/" includes a
// "." if it is an external import. However, there are some import paths in the
// monorepo that breaks this. We check against these to determine that they are
// not part of the Go standard library.
var knownImportPrefixes = map[string]struct{}{
	"glue":      {},
	"gogoproto": {},
	"mock":      {},
	"thriftrw":  {},
}

func isStdlibImport(path string) bool {
	if i := strings.IndexByte(path, '/'); i >= 0 {
		path = path[:i]
	}

	// If the prefix of the import path contains a ".", it should be considered
	// to be a external package (not part of Go standard lib).
	if strings.Contains(path, ".") {
		return false
	}

	// TODO: before moving to go1.13, we have import paths for some generated
	// code in the monorepo that don't have a "." in them, i.e. thriftrw/code.uber.internal/foo/bar,
	// because Bazel is able to resolve arbitrary import paths. Starting from
	// go1.13, there is a hard requirement for the prefix of the import path to
	// include a "." (Issue tracked: T4289809)
	// For now, filter the known import prefixes in the monorepo to be recognized
	// as an imported dependency / not part of the Go standard lib.
	if _, ok := knownImportPrefixes[path]; ok {
		return false
	}

	return true
}
