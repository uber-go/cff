package pkg

import (
	"errors"
	"go/token"
	"testing"

	"code.uber.internal/go/importer"
	flags "github.com/jessevdk/go-flags"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArchiveLoaderFactory_Flags(t *testing.T) {
	tests := []struct {
		desc           string
		give           []string
		wantArchives   []string
		wantSources    []string
		wantStdlibRoot string
	}{
		{desc: "empty"},
		{
			desc: "archive",
			give: []string{
				"--archive", "github.com/foo/bar:github.com/foo/baz=github.com/foo/bar=bar.go=bar_export.go",
			},
			wantArchives: []string{
				"github.com/foo/bar:github.com/foo/baz=github.com/foo/bar=bar.go=bar_export.go",
			},
		},
		{
			desc: "sources",
			give: []string{
				"--source", "foo.go",
				"--source", "bar.go",
			},
			wantSources: []string{"foo.go", "bar.go"},
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

			parser := flags.NewParser(&struct{}{}, 0)
			loader, err := factory.RegisterFlags(parser)
			require.NoError(t, err)
			al, ok := loader.(*archiveLoader)
			require.True(t, ok, "expected an *archiveLoader, got %T", loader)

			extra, err := parser.ParseArgs(tt.give)
			require.NoError(t, err)
			assert.Empty(t, extra, "unexpected unparsed arguments")

			assert.Equal(t, tt.wantArchives, *al.archives, "archives")
			assert.Equal(t, tt.wantSources, *al.srcs, "sources")
			assert.Equal(t, tt.wantStdlibRoot, *al.stdlibRoot, "stdlibroot")
		})
	}
}

func TestArchiveLoader(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var called bool
		defer func() {
			assert.True(t, called, "loadArchive was never called")
		}()

		archives := []string{
			"github.com/foo/bar:github.com/foo/baz=github.com/foo/bar=bar.go=bar_export.go",
		}
		srcs := []string{"foo.go", "bar.go"}
		stdlibRoot := "std"
		importPath := "example.com/foo"

		loader := archiveLoader{
			archives:   &archives,
			srcs:       &srcs,
			stdlibRoot: &stdlibRoot,
			loadArchive: func(lp importer.LoadParams) (*importer.Package, error) {
				called = true
				assert.Equal(t, importPath, lp.ImportPath)
				assert.Equal(t, srcs, lp.Srcs)
				assert.Equal(t, stdlibRoot, lp.StdlibRoot)
				assert.Equal(t, []importer.Archive{
					{
						ImportMap: "github.com/foo/bar",
						File:      "bar.go",
					},
				}, lp.Archives)
				return &importer.Package{}, nil
			},
		}

		pkgs, err := loader.Load(token.NewFileSet(), importPath)
		require.NoError(t, err)
		assert.Len(t, pkgs, 1)
	})

	t.Run("bad archive", func(t *testing.T) {
		archives := []string{"not valid"}
		loader := archiveLoader{
			archives:   &archives,
			srcs:       new([]string),
			stdlibRoot: new(string),
			loadArchive: func(lp importer.LoadParams) (*importer.Package, error) {
				t.Fatal("loadArchive should not be called")
				return &importer.Package{}, nil
			},
		}

		_, err := loader.Load(token.NewFileSet(), "example.com/foo")
		assert.ErrorContains(t, err, "invalid argument --archive")
	})

	t.Run("load error", func(t *testing.T) {
		giveErr := errors.New("great sadness")
		loader := archiveLoader{
			archives:   new([]string),
			srcs:       new([]string),
			stdlibRoot: new(string),
			loadArchive: func(lp importer.LoadParams) (*importer.Package, error) {
				return nil, giveErr
			},
		}

		_, err := loader.Load(token.NewFileSet(), "example.com/foo")
		assert.ErrorIs(t, err, giveErr)
	})
}

func TestParseArchive(t *testing.T) {
	type testCase struct {
		desc    string
		give    string
		want    importer.Archive
		wantErr string
	}

	tests := []testCase{
		{
			desc: "success",
			give: "foo:bar=hello=world=baz",
			want: importer.Archive{
				ImportMap: "hello",
				File:      "world",
			},
		},
		{
			desc:    "failure",
			give:    "foo=bar=hello=world=baz",
			wantErr: "expected 4 elements, got 5",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			if tt.wantErr == "" {
				arc, err := parseArchive(tt.give)
				require.NoError(t, err, "expected no error parsing %q", tt.give)
				assert.Equal(t, tt.want, arc)
			} else {
				_, err := parseArchive(tt.give)
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErr)
			}
		})
	}
}
