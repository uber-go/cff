package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestString(t *testing.T) {
	tests := []struct {
		desc  string
		str   string
		value String
	}{
		{desc: "empty"},
		{
			desc:  "simple",
			str:   "foo",
			value: String("foo"),
		},
		{
			desc:  "space",
			str:   "foo bar",
			value: String("foo bar"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var got String
			fset := NewSet("foo")
			fset.Var(&got, "x", "")
			require.NoError(t, fset.Parse([]string{"-x", tt.str}))
			assert.Equal(t, tt.value, got)

			t.Run("String", func(t *testing.T) {
				assert.Equal(t, tt.str, got.String())
			})

			t.Run("Get", func(t *testing.T) {
				assert.Equal(t, tt.str, got.Get())
			})
		})
	}
}
