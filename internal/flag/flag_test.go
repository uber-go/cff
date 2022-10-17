package flag

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewSet(t *testing.T) {
	fset := NewSet("foo")
	fset.Int("flag", 42, "test flag")
	err := fset.Parse([]string{"-flag", "not an int"})
	require.Error(t, err)
}
