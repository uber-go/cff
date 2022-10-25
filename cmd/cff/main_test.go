package main

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/cff/internal/flag"
)

func TestParseArgs(t *testing.T) {
	tests := []struct {
		desc    string
		give    []string
		want    params
		wantErr string
	}{
		{
			desc:    "empty",
			wantErr: "please provide an import path",
		},
		{
			desc: "import path",
			give: []string{"example.com/foo"},
			want: params{
				ImportPath: "example.com/foo",
				GenMode:    flag.BaseMode,
			},
		},
		{
			desc:    "too many import paths",
			give:    []string{"example.com/foo", "example.com/bar"},
			wantErr: "too many import paths",
		},
		{
			desc: "files",
			give: []string{"-file", "foo.go", "-file", "bar.go=baz.go", "example.com/foo"},
			want: params{
				GenMode: flag.BaseMode,
				Files: []flag.InOutPair{
					{Input: "foo.go"},
					{Input: "bar.go", Output: "baz.go"},
				},
				ImportPath: "example.com/foo",
			},
		},
		{
			desc: "instrument all",
			give: []string{"-auto-instrument", "example.com/foo"},
			want: params{
				GenMode:        flag.BaseMode,
				AutoInstrument: true,
				ImportPath:     "example.com/foo",
			},
		},
		{
			desc: "gen mode",
			give: []string{"-genmode", "source-map", "example.com/foo"},
			want: params{
				GenMode:    flag.SourceMapMode,
				ImportPath: "example.com/foo",
			},
		},
		{
			desc: "quiet",
			give: []string{"-quiet", "example.com/foo"},
			want: params{
				GenMode:    flag.BaseMode,
				Quiet:      true,
				ImportPath: "example.com/foo",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			_, got, err := parseArgs(io.Discard, tt.give)
			if tt.wantErr != "" {
				assert.ErrorContains(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, *got)
		})
	}
}

func TestMain_ErrorNoPanic(t *testing.T) {
	assert.NotPanics(t, func() {
		err := run([]string{"-unknown-flag"})
		assert.Error(t, err)
	})
}

func TestGenFilename(t *testing.T) {
	tests := []struct {
		desc    string
		give    string
		want    string
	}{
		{
			desc: "source file",
			give: "x/y/foo.go",
			want: "x/y/foo_gen.go",
		},
		{
			desc: "test file",
			give: "x/y/foo_test.go",
			want: "x/y/foo_gen_test.go",
		},
		{
			desc: "source file abs path",
			give: "/x/y/foo.go",
			want: "/x/y/foo_gen.go",
		},
		{
			desc: "test file abs path",
			give: "/x/y/foo_test.go",
			want: "/x/y/foo_gen_test.go",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			assert.Equal(t, tt.want, genFilename(tt.give))
		})
	}
}
