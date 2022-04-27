package insidegeneric

import (
	"context"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJoinTwo(t *testing.T) {
	var calledA, calledB, calledC bool
	c, err := JoinTwo(
		func(context.Context) (int, error) {
			calledA = true
			return 42, nil
		},
		func(context.Context) (string, error) {
			calledB = true
			return "42", nil
		},
		func(i int, s string) bool {
			calledC = true
			return strconv.Itoa(i) == s
		},
	)
	require.NoError(t, err)
	assert.True(t, c, `str(42) should be "42"`)
	assert.True(t, calledA, "A wasn't called")
	assert.True(t, calledB, "B wasn't called")
	assert.True(t, calledC, "C wasn't called")
}

func TestJoinMany(t *testing.T) {
	const N = 5

	called := make([]bool, N)
	defer func() {
		for i, called := range called {
			assert.True(t, called, "producer %d was not called", i)
		}
	}()

	producers := make([]Producer[string], N)
	for i := 0; i < N; i++ {
		i := i
		producers[i] = func(context.Context) (string, error) {
			called[i] = true
			return strconv.Itoa(i), nil
		}
	}

	got, err := JoinMany(producers...)
	require.NoError(t, err)
	assert.Equal(t, []string{"0", "1", "2", "3", "4"}, got)
}
