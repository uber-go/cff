package mode

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseGenMode(t *testing.T) {
	type testCase struct {
		desc string
		give string
		want GenerationMode
	}

	tests := []testCase{
		{
			desc: "base",
			give: "base",
			want: Base,
		},
		{
			desc: "source mapped",
			give: "source-map",
			want: SourceMap,
		},
		{
			desc: "unknown",
			give: "sad",
			want: Unknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var genMode GenerationMode
			genMode.UnmarshalText([]byte(tt.give))
			assert.Equal(t, genMode, tt.want)
		})
	}
}
