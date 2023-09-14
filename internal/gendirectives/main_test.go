package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRunErrors(t *testing.T) {
	dir := t.TempDir()

	t.Run("too few arguments", func(t *testing.T) {
		require.Error(t, run([]string{"foo"}))
	})

	t.Run("unable to parse", func(t *testing.T) {
		input := filepath.Join(dir, "bad_syntax.go")
		require.NoError(t, os.WriteFile(input, []byte("foo"), 0o644))

		err := run([]string{input, filepath.Join(dir, "out.go")})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "expected 'package'")
	})

	t.Run("unable to write", func(t *testing.T) {
		input := filepath.Join(dir, "cff.go")
		require.NoError(t, os.WriteFile(input, []byte(_sampleFile), 0o644))

		err := run([]string{input, filepath.Join(dir, "does_not_exist", "out.go")})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no such file or directory")
	})
}

const _sampleFile = `
package whatever

func Foo() {}

type Bar struct{}

func (*Bar) Baz()
`

func TestRun(t *testing.T) {
	dir := t.TempDir()

	input := filepath.Join(dir, "cff.go")
	require.NoError(t, os.WriteFile(input, []byte(_sampleFile), 0o644))

	output := filepath.Join(dir, "out.go")
	require.NoError(t, run([]string{input, output}))

	out, err := os.ReadFile(output)
	require.NoError(t, err)

	got := string(out)
	assert.Contains(t, got, `"Foo"`)
	assert.NotContains(t, got, `"Bar"`)
	assert.NotContains(t, got, `"Baz"`)
}
