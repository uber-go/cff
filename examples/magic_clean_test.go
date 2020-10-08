package example_test

import (
	"io/ioutil"
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

	expected, err := ioutil.ReadFile(expectedPath)
	require.NoError(t, err)
	actual, err := ioutil.ReadFile(actualPath)
	require.NoError(t, err)

	assert.Equal(t, string(expected), string(actual), "magic_gen.go is out of date, try running these commands:\n"+
		os.ExpandEnv("cd $GOPATH/src/go.uber.org/cff/examples\n")+
		"rm -f magic_gen.go &&  $GOPATH/bin/cff --file=magic.go go.uber.org/cff/examples")
}
