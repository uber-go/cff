package internal

import (
	"fmt"
	"go/token"
	"go/types"
	"os"

	"github.com/bazelbuild/bazel-gazelle/language/go"
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
	if golang.IsStandard(path) {
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
