package flag

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMode(t *testing.T) {
	tests := []struct {
		name string
		mode Mode
	}{
		{
			name: "base",
			mode: BaseMode,
		},
		{
			name: "source-map",
			mode: SourceMapMode,
		},
		{
			name: "modifier",
			mode: ModifierMode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got Mode
			fset := NewSet("foo")
			fset.Var(&got, "x", "")
			require.NoError(t, fset.Parse([]string{"-x", tt.name}))
			assert.Equal(t, tt.mode, got)

			t.Run("String", func(t *testing.T) {
				assert.Equal(t, tt.name, tt.mode.String())
			})

			t.Run("Get", func(t *testing.T) {
				assert.Equal(t, tt.mode, got.Get())
			})

			t.Run("UnmarshalText", func(t *testing.T) {
				var got Mode
				require.NoError(t, got.UnmarshalText([]byte(tt.name)))
				assert.Equal(t, tt.mode, got)
			})
		})
	}
}

func TestModeUnknown_String(t *testing.T) {
	tests := []Mode{0, 10, 20}
	for _, tt := range tests {
		t.Run(fmt.Sprint(int(tt)), func(t *testing.T) {
			assert.Equal(t, "unknown", tt.String())
		})
	}
}

func TestModeUnknown_Unmarshal(t *testing.T) {
	tests := []string{"foo", "bar", "unknown"}
	for _, tt := range tests {
		t.Run(tt, func(t *testing.T) {
			var m Mode
			assert.ErrorContains(t, m.Set(tt), "unknown mode")
		})
	}
}
