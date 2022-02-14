package shadowedvar

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCtxConflict(t *testing.T) {
	msg, err := CtxConflict("hello")
	assert.NoError(t, err)
	assert.Equal(t, "hello", msg)
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
