package pkg

import (
	"errors"
	"go/token"
	"testing"

	"go.uber.org/cff/internal/flag"
	"code.uber.internal/go/importer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArchiveLoaderFactory_Flags(t *testing.T) {
	tests := []struct {
		desc           string
		give           []string
		wantArchives   []archiveValue
		wantSources    []flag.String
		wantStdlibRoot flag.String
	}{
		{desc: "empty"},
		{
			desc: "archive",
			give: []string{
				"--archive", "github.com/foo/bar:github.com/foo/baz=github.com/foo/bar=bar.go=bar_export.go",
			},
			wantArchives: []archiveValue{
				{
					ImportPaths: []string{"github.com/foo/bar", "github.com/foo/baz"},
					ImportMap:   "github.com/foo/bar",
					File:        "bar.go",
					Export:      "bar_export.go",
				},
			},
		},
		{
			desc: "sources",
			give: []string{
				"--source", "foo.go",
				"--source", "bar.go",
			},
			wantSources: []flag.String{"foo.go", "bar.go"},
		},
		{
			desc:           "stdlibroot",
			give:           []string{"--stdlibroot", "foo"},
			wantStdlibRoot: "foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var factory ArchiveLoaderFactory

			parser := flag.NewSet("cff")
			loader := factory.RegisterFlags(parser)
			al, ok := loader.(*archiveLoader)
			require.True(t, ok, "expected an *archiveLoader, got %T", loader)

			require.NoError(t, parser.Parse(tt.give))
			assert.Empty(t, parser.Args(), "unexpected unparsed arguments")

			assert.Equal(t, tt.wantArchives, al.archives, "archives")
			assert.Equal(t, tt.wantSources, al.srcs, "sources")
			assert.Equal(t, tt.wantStdlibRoot, al.stdlibRoot, "stdlibroot")
		})
	}
}

func TestArchiveLoader(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var called bool
		defer func() {
			assert.True(t, called, "loadArchive was never called")
		}()

		importPath := "example.com/foo"
		loader := archiveLoader{
			archives: []archiveValue{
				{
					ImportPaths: []string{"github.com/foo/bar", "github.com/foo/baz"},
					ImportMap:   "github.com/foo/bar",
					File:        "bar.go",
					Export:      "bar_export.go",
				},
			},
			srcs:       []flag.String{"foo.go", "bar.go"},
			stdlibRoot: "std",
			loadArchive: func(lp importer.LoadParams) (*importer.Package, error) {
				called = true
				assert.Equal(t, importPath, lp.ImportPath)
				assert.Equal(t, []string{"foo.go", "bar.go"}, lp.Srcs)
				assert.Equal(t, "std", lp.StdlibRoot)
				assert.Equal(t, []importer.Archive{
					{
						ImportPaths: []string{
							"github.com/foo/bar",
							"github.com/foo/baz",
						},
						ImportMap:  "github.com/foo/bar",
						File:       "bar.go",
						ExportFile: "bar_export.go",
					},
				}, lp.Archives)
				return &importer.Package{}, nil
			},
		}

		pkgs, err := loader.Load(token.NewFileSet(), importPath)
		require.NoError(t, err)
		assert.Len(t, pkgs, 1)
	})

	t.Run("load error", func(t *testing.T) {
		giveErr := errors.New("great sadness")
		loader := archiveLoader{
			loadArchive: func(lp importer.LoadParams) (*importer.Package, error) {
				return nil, giveErr
			},
		}

		_, err := loader.Load(token.NewFileSet(), "example.com/foo")
		assert.ErrorIs(t, err, giveErr)
	})
}
