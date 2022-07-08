package shadowedvar

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCtxConflict(t *testing.T) {
	msg, err := CtxConflict("hello")
	assert.NoError(t, err)
	assert.Equal(t, "hello", msg)
}

func TestCtxConflictParallel(t *testing.T) {
	msg1, msg2, err := CtxConflictParallel("hello")
	assert.NoError(t, err)
	assert.Equal(t, "hello", msg1)
	assert.Equal(t, "hello", msg2)
}

func TestCtxConflictSlice(t *testing.T) {
	ctx := "hello"
	target := []string{"A", "B", "C"}
	require.NoError(t, CtxConflictSlice(ctx, target))
	assert.Equal(t, []string{"helloA", "helloB", "helloC"}, target)
}

func TestCtxConflictMap(t *testing.T) {
	ctx := 5
	input := map[int]int{
		0: 10,
		1: 15,
	}
	out, err := CtxConflictMap(ctx, input)
	require.NoError(t, err)
	assert.Equal(t, 15, out[0])
	assert.Equal(t, 20, out[1])
}

func TestPredicateCtxConflict(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		called := false
		err := PredicateCtxConflict(func() { called = true }, true)
		assert.NoError(t, err)
		assert.True(t, called, "function was never called")
	})

	t.Run("false", func(t *testing.T) {
		err := PredicateCtxConflict(func() {
			t.Fatal("function must never be called")
		}, false)
		assert.NoError(t, err)
	})
}
