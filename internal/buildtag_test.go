package internal

import (
	"bytes"
	"errors"
	"go/build/constraint"
	"go/parser"
	"go/token"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHasCFFTag(t *testing.T) {
	tests := []struct {
		desc string
		give string
		want bool
	}{
		{
			desc: "non cff tag",
			give: "// +build foo",
		},
		{
			desc: "non cff constraint",
			give: "//go:build foo",
		},
		{
			desc: "cff tag",
			give: "// +build cff",
			want: true,
		},
		{
			desc: "cff constraint",
			give: "//go:build cff",
			want: true,
		},
		{
			desc: "conditional cff constraint",
			give: "//go:build foo && cff",
			want: true,
		},
		{
			desc: "inverted cff",
			give: "//go:build (foo && !cff)",
			want: true,
		},
		{
			desc: "nested cff",
			give: "//go:build foo && (!bar || !(baz && cff))",
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			exp, err := constraint.Parse(tt.give)
			require.NoError(t, err)

			assert.Equal(t, tt.want, hasCFFTag(exp))
		})
	}
}

func TestFileHasCFFTag(t *testing.T) {
	tests := []struct {
		desc string
		give []string
		want bool
	}{
		{
			desc: "nothing",
			give: []string{
				"package foo",
			},
		},
		{
			desc: "non cff tags",
			give: []string{
				"// +build foo",
				"",
				"package foo",
			},
		},
		{
			desc: "non cff constraint",
			give: []string{
				"//go:build foo && bar",
				"",
				"package foo",
			},
		},
		{
			desc: "cff tag",
			give: []string{
				"// +build cff",
				"",
				"package foo",
			},
			want: true,
		},
		{
			desc: "cff constraint",
			give: []string{
				"//go:build cff",
				"",
				"package foo",
			},
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			contents := strings.Join(tt.give, "\n")
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, "foo.go", contents, parser.ParseComments)
			require.NoError(t, err)
			assert.Equal(t, tt.want, fileHasCFFTag(f))
		})
	}
}

func TestWriteInvertedCFFTag(t *testing.T) {
	tests := []struct {
		desc string
		give string
		want string
	}{
		{
			desc: "empty",
			give: "\n",
			want: "\n",
		},
		{
			desc: "tag/single",
			give: "// +build cff\n",
			want: "// +build !cff\n",
		},
		{
			desc: "tag/multiple",
			give: "// +build cff go1.19\n",
			want: "// +build !cff go1.19\n",
		},
		{
			desc: "tag/cff multiple times",
			give: "// +build cff cff cff\n",
			want: "// +build !cff !cff !cff\n",
		},
		{
			desc: "tag/multiple lines",
			give: "// +build cff\n// +build go1.19 cff\n",
			want: "// +build !cff\n// +build go1.19 !cff\n",
		},
		{
			desc: "tag/package clause",
			give: "// +build cff\n\npackage foo\n",
			want: "// +build !cff\n\npackage foo\n",
		},
		{
			desc: "tag/not",
			give: "// +build !go1.19 cff\n",
			want: "// +build !go1.19 !cff\n",
		},
		{
			desc: "constraint/single",
			give: "//go:build cff\n",
			want: "//go:build !cff\n",
		},
		{
			desc: "constraint/multiple and",
			give: "//go:build cff && go1.19\n",
			want: "//go:build !cff && go1.19\n",
		},
		{
			desc: "constraint/multiple or",
			give: "//go:build cff || go1.19\n",
			want: "//go:build !cff || go1.19\n",
		},
		{
			desc: "constraint/complex expression",
			give: "//go:build (go1.19 && cff) || (go1.17 && !cff)\n",
			want: "//go:build (go1.19 && !cff) || (go1.17 && cff)\n",
		},
		{
			desc: "constraint/package clause",
			give: "//go:build cff\n\npackage foo\n",
			want: "//go:build !cff\n\npackage foo\n",
		},
		{
			desc: "constraint/complex not",
			give: "//go:build !(go1.19 && cff)\n",
			want: "//go:build !(go1.19 && !cff)\n",
		},
		{
			desc: "both",
			give: "//go:build cff\n// +build cff\n\npackage foo\n",
			want: "//go:build !cff\n// +build !cff\n\npackage foo\n",
		},
		{
			desc: "invalid constraint",
			give: "//go:build not a cff constraint\n",
			want: "//go:build not a cff constraint\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var got bytes.Buffer
			err := writeInvertedCFFTag(&got, []byte(tt.give))
			require.NoError(t, err)
			assert.Equal(t, tt.want, got.String())
		})
	}
}

func TestStickyErrWriter(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		var buff bytes.Buffer
		w := &stickyErrWriter{W: &buff}
		io.WriteString(w, "hello")
		assert.NoError(t, w.Err)
	})

	t.Run("failure", func(t *testing.T) {
		giveErr := errors.New("great sadness")
		w := &stickyErrWriter{
			W: &errWriter{Err: giveErr},
		}
		io.WriteString(w, "hello")
		io.WriteString(w, "world")
		assert.ErrorIs(t, w.Err, giveErr)
	})
}

// Writer that always fails with the given error.
type errWriter struct{ Err error }

func (w *errWriter) Write([]byte) (int, error) {
	return 0, w.Err
}
