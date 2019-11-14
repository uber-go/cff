package internal

import (
	"go/token"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArchiveLoad(t *testing.T) {
	type testCase struct {
		desc            string
		importPath      string
		srcs            []string
		stdlibroot      string
		depArchives     []Archive
		expectedImports []string
	}

	tests := []testCase{
		{
			desc:       "success",
			importPath: "internal/archive_tests",
			srcs:       []string{"archive_tests/testfile.go"},
			stdlibroot: "archive_tests/stdlibroot",
			depArchives: []Archive{
				Archive{
					ImportMap: "example.import/archivedata",
					File:      "archive_tests/archiveroot/errors-ae16.a",
				},
			},
			expectedImports: []string{"context", "example.import/archivedata"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			pkgs, err := PackagesArchive(LoadParams{
				Fset:       token.NewFileSet(),
				ImportPath: tt.importPath,
				Srcs:       tt.srcs,
				StdlibRoot: tt.stdlibroot,
				Archives:   tt.depArchives,
			})
			require.NoError(t, err)
			require.Len(t, pkgs, 1)
			pkg := pkgs[0]
			// packages should have the same number of CompiledGoFiles and Syntax files
			assert.Equal(t, tt.srcs, pkg.CompiledGoFiles)
			assert.Len(t, pkg.Syntax, len(pkg.CompiledGoFiles))

			// the types.Package should have the provided importPath
			typesPkg := pkg.Types
			assert.Equal(t, tt.importPath, typesPkg.Path())
			imports := typesPkg.Imports()
			actualImports := make([]string, len(imports))
			for i, imp := range imports {
				actualImports[i] = imp.Path()
			}
			assert.ElementsMatch(t, tt.expectedImports, actualImports)
		})
	}
}
