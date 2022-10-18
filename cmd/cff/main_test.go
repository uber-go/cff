package main

import (
	"io"
	"testing"

	"go.uber.org/cff/internal/flag"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
