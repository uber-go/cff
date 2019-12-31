package noinfra

import (
	"go/ast"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"testing"

	"code.uber.internal/base/testing/envtest"
	"go.uber.org/cff"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/analysis/analysistest"
)

func TestCFF2Import(t *testing.T) {
	tests := []struct {
		desc    string
		pkgpath string
		// Note, expectation is set in the // want ... comment on specific line. Not my choice.
		files map[string]string
	}{
		{
			desc:    "cff import present",
			pkgpath: "main",
			files: map[string]string{"main/main.go": `package main

import (
	"context"

	"go.uber.org/cff" // want "found usage of go.uber.org/cff"
)


func SimpleFlow() error {
	return cff.Flow(context.Background())
}`,
				"go.uber.org/cff/cff.go": `package cff

import "context"

type FlowOption interface {}

func Flow(ctx context.Context, opts ...FlowOption) error {
	panic("code not generated; run cff")
}
`,
			},
		},
		{
			desc:    "cff import present aliased",
			pkgpath: "main",
			files: map[string]string{"main/main.go": `package main

import (
	"context"

	cff2 "go.uber.org/cff" // want "found usage of go.uber.org/cff"
)


func SimpleFlow() error {
	return cff2.Flow(context.Background())
}`,
				"go.uber.org/cff/cff.go": `package cff

import "context"

type FlowOption interface {}

func Flow(ctx context.Context, opts ...FlowOption) error {
	panic("code not generated; run cff")
}
`,
			},
		},
		{
			desc:    "cff import not present",
			pkgpath: "main",
			files: map[string]string{"main/main.go": `package main

import (
	"fmt"
)
func main() {
    fmt.Println("hello")
}
`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			defer envtest.OverrideVar(t, "PATH", os.ExpandEnv("$TEST_SRCDIR/__main__/external/go_sdk/bin:$PATH")).Reset()
			// Go executable requires a GOCACHE to be set after go1.12.
			tmpDir := os.Getenv("TEST_TMPDIR")
			if file, err := os.Stat(tmpDir); err == nil {
				if file.IsDir() {
					// Go executable requires a GOCACHE to be set after go1.12.
					defer envtest.OverrideVar(t, "GOCACHE", filepath.Join(tmpDir, "/cache")).Reset()
				}
			}
			dir, cleanup, err := analysistest.WriteFiles(tt.files)
			require.NoError(t, err)
			defer cleanup()
			analysistest.Run(t, dir, Analyzer, tt.pkgpath)
		})
	}
}

func TestImportPath(t *testing.T) {
	tests := []struct {
		desc string
		give string
		want string
	}{
		{
			desc: "unquote",
			give: `"go.uber.org/cff"`,
			want: "go.uber.org/cff",
		},
		{
			desc: "no match",
			give: `"context"`,
			want: "context",
		},
		{
			desc: "unquoted import",
			give: "context",
			want: "", // ErrSyntax in strconv.Unquote
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			n := &ast.ImportSpec{
				Path: &ast.BasicLit{
					Value: tt.give,
				},
			}
			assert.Equal(t, tt.want, importPath(n))
		})
	}
}

func TestIsCFF2ImportPath(t *testing.T) {
	tests := []struct {
		desc string
		give string
		want bool
	}{
		{
			desc: "exact match",
			give: "go.uber.org/cff",
			want: true,
		},
		{
			desc: "no match",
			give: "context",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.want, isCFF2Import(tt.give))
		})
	}
}

// TestMatchCFF2Location checks that linter doesn't become stale if CFF2's location is changed.
func TestMatchCFF2Location(t *testing.T) {
	out := runtime.FuncForPC(reflect.ValueOf(cff.Flow).Pointer()).Name()
	currentImport := strings.TrimSuffix(out, ".Flow") // this should match _cff2Pkg.
	assert.True(t, isCFF2Import(currentImport))
}
