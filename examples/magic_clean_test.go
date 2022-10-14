package example_test

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGoldenMagic provides a test that asserts that the checked-in generated code in magic_gen.go is cleanly generated
// so that updates to the template without updating magic_gen.go trigger a failing test.
func TestGoldenMagic(t *testing.T) {
	expectedPath := path.Join(os.Getenv("TEST_SRCDIR"), "__main__/src/go.uber.org/cff/examples/cff_/magic_gen.go")
	actualPath := path.Join(os.Getenv("TEST_SRCDIR"), "__main__/src/go.uber.org/cff/examples/magic_gen.go")

	expected, err := os.ReadFile(expectedPath)
	require.NoError(t, err)
	actual, err := os.ReadFile(actualPath)
	require.NoError(t, err)

	if !assert.Equal(t, string(expected), string(actual), "magic_gen.go is out of date") {
		from := "$GOPATH/bazel-bin/src/go.uber.org/cff/examples/cff_/magic_gen.go"
		to := "$GOPATH/src/go.uber.org/cff/examples/magic_gen.go"
		t.Log(
			"Try running these commands:\n" +
				"bazel build //src/go.uber.org/cff/examples:cff &&\n" +
				fmt.Sprintf("  cp %q %q", from, to),
		)
	}
}

// TestGoldenMagic2 provides a test that asserts that the checked-in generated code in magic_gen_v2.go is cleanly generated
// so that updates to the template without updating magic_gen_v2.go trigger a failing test.
func TestGoldenMagic2(t *testing.T) {
	expectedPath := path.Join(os.Getenv("TEST_SRCDIR"), "__main__/src/go.uber.org/cff/examples/cff_v2_/magic_v2_gen.go")
	actualPath := path.Join(os.Getenv("TEST_SRCDIR"), "__main__/src/go.uber.org/cff/examples/magic_v2_gen.go")

	expected, err := os.ReadFile(expectedPath)
	require.NoError(t, err)
	actual, err := os.ReadFile(actualPath)
	require.NoError(t, err)

	if !assert.Equal(t, string(expected), string(actual), "magic_v2_gen.go is out of date") {
		from := "$GOPATH/bazel-bin/src/go.uber.org/cff/examples/cff_v2_/magic_v2_gen.go"
		to := "$GOPATH/src/go.uber.org/cff/examples/magic_v2_gen.go"
		t.Log(
			"Try running these commands:\n" +
				"bazel build //src/go.uber.org/cff/examples:cff_v2 &&\n" +
				fmt.Sprintf("  cp %q %q", from, to),
		)
	}
}
