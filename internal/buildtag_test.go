package internal

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
