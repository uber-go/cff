package example_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGoldenMagic provides a test that asserts that the checked-in generated code in magic_gen.go is cleanly generated
// so that updates to the template without updating magic_gen.go trigger a failing test.
func TestGoldenMagic(t *testing.T) {
	expectedPath := filepath.Join(t.TempDir(), "magic_gen.go")
	actualPath := "magic_gen.go"

	cmd := exec.Command("cff", "-genmode=source-map", "-file", "magic.go="+expectedPath, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	require.NoError(t, cmd.Run())

	expected, err := os.ReadFile(expectedPath)
	require.NoError(t, err)
	actual, err := os.ReadFile(actualPath)
	require.NoError(t, err)

	if !assert.Equal(t, string(expected), string(actual), "magic_gen.go is out of date") {
		t.Log("Try running 'make generate'")
	}
}

// TestGoldenMagic2 provides a test that asserts that the checked-in generated code in magic_gen_v2.go is cleanly generated
// so that updates to the template without updating magic_gen_v2.go trigger a failing test.
func TestGoldenMagic2(t *testing.T) {
	expectedPath := filepath.Join(t.TempDir(), "magic_v2.go")
	actualPath := "magic_v2_gen.go"

	cmd := exec.Command("cff", "-genmode=modifier", "-tags=v2", "-file", "magic_v2.go="+expectedPath, ".")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	require.NoError(t, cmd.Run())

	expected, err := os.ReadFile(expectedPath)
	require.NoError(t, err)
	actual, err := os.ReadFile(actualPath)
	require.NoError(t, err)

	if !assert.Equal(t, string(expected), string(actual), "magic_v2_gen.go is out of date") {
		t.Log("Try running 'make generate'")
	}
}
