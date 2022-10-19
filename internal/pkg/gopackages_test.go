package pkg

import (
	"errors"
	"go/token"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/cff/internal/flag"
	"golang.org/x/tools/go/packages"
)

func TestGoPackagesLoader_Integration(t *testing.T) {
	wd, err := os.Getwd()
	require.NoError(t, err)
	testdata := filepath.Join(wd, "testdata/gopackages")
	factory := GoPackagesLoaderFactory{
		dir: testdata,
	}

	parser := flag.NewSet("cff")
	loader := factory.RegisterFlags(parser)
	require.NoError(t, parser.Parse([]string{"-tags", "foo,cff"}),
		"parse arguments")

	fset := token.NewFileSet()
	pkgs, err := loader.Load(fset, "example.com/foo")
	require.NoError(t, err, "load package")

	require.Len(t, pkgs, 1, "expected one package")
	pkg := pkgs[0]

	assert.Len(t, pkg.CompiledGoFiles, 2, "wrong number of files")
	assert.Len(t, pkg.Syntax, 2, "wrong number ASTs")
	assert.NotNil(t, pkg.Types, "missing Types")
	assert.NotNil(t, pkg.TypesInfo, "missing type information")

	baseNames := make([]string, len(pkg.CompiledGoFiles))
	for i, file := range pkg.CompiledGoFiles {
		baseNames[i] = filepath.Base(file)
	}
	assert.ElementsMatch(t, []string{"a.go", "b.go"}, baseNames)
}

func TestGoPackagesLoader_Errors(t *testing.T) {
	t.Run("load error", func(t *testing.T) {
		giveErr := errors.New("great sadness")
		loader := goPackagesLoader{
			load: func(c *packages.Config, s ...string) ([]*packages.Package, error) {
				return nil, giveErr
			},
		}

		_, err := loader.Load(token.NewFileSet(), "example.com/foo")
		assert.ErrorIs(t, err, giveErr)
	})

	t.Run("no packages", func(t *testing.T) {
		loader := goPackagesLoader{
			load: func(c *packages.Config, s ...string) ([]*packages.Package, error) {
				return nil, nil
			},
		}

		_, err := loader.Load(token.NewFileSet(), "example.com/foo")
		assert.ErrorContains(t, err, "no packages found")
	})

	t.Run("package errors", func(t *testing.T) {
		giveErrors := []packages.Error{
			{
				Pos:  "foo.go:1:2",
				Msg:  "great sadness",
				Kind: packages.ParseError,
			},
			{
				Pos:  "bar.go:1:2",
				Msg:  "massive fail",
				Kind: packages.TypeError,
			},
		}

		loader := goPackagesLoader{
			load: func(c *packages.Config, s ...string) ([]*packages.Package, error) {
				return []*packages.Package{
					{
						Errors: giveErrors,
					},
				}, nil
			},
		}

		_, err := loader.Load(token.NewFileSet(), "example.com/foo")
		for _, wantErr := range giveErrors {
			assert.ErrorIs(t, err, wantErr)
		}
	})
}
