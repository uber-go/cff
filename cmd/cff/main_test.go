package main

import (
	"testing"

	"code.uber.internal/go/importer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
