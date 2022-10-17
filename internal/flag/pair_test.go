package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInOutPair(t *testing.T) {
	tests := []struct {
		desc string
		give string
		want InOutPair
		str  string
	}{
		{
			desc: "input",
			give: "foo",
			want: InOutPair{Input: "foo"},
			str:  "foo",
		},
		{
			desc: "input output",
			give: "foo=bar",
			want: InOutPair{Input: "foo", Output: "bar"},
			str:  "foo=bar",
		},
		{
			desc: "empty output",
			give: "foo=",
			want: InOutPair{Input: "foo"},
			str:  "foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var got InOutPair
			fset := NewSet("foo")
			fset.Var(&got, "x", "")
			require.NoError(t, fset.Parse([]string{"-x", tt.str}))
			assert.Equal(t, tt.want, got)

			t.Run("String", func(t *testing.T) {
				assert.Equal(t, tt.str, tt.want.String())
			})

			t.Run("Get", func(t *testing.T) {
				assert.Equal(t, &tt.want, got.Get())
			})
		})
	}
}

func TestInOutPair_Error(t *testing.T) {
	var p InOutPair
	assert.ErrorContains(t, p.Set("=foo"), "cannot be empty")
}
