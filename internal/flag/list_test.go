package flag

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	tests := []struct {
		desc       string
		give       []string
		want       []String
		wantString string
		wantGet    any
	}{
		{
			desc:    "empty",
			wantGet: []any{},
		},
		{
			desc:       "one",
			give:       []string{"-x", "foo"},
			want:       []String{"foo"},
			wantString: "foo",
			wantGet:    []any{"foo"},
		},
		{
			desc:       "multiple",
			give:       []string{"-x", "foo", "-x", "bar"},
			want:       []String{"foo", "bar"},
			wantString: "foo; bar",
			wantGet:    []any{"foo", "bar"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			var got []String
			fset := NewSet("foo")
			fset.Var(AsList(&got), "x", "")
			require.NoError(t, fset.Parse(tt.give))
			assert.Equal(t, tt.want, got)

			t.Run("String", func(t *testing.T) {
				assert.Equal(t, tt.wantString, AsList(&got).String())
			})

			t.Run("Get", func(t *testing.T) {
				assert.Equal(t, tt.wantGet, AsList(&got).Get())
			})
		})
	}
}

func TestList_Error(t *testing.T) {
	var list List[InOutPair, *InOutPair]
	assert.Error(t, list.Set(""))
}
