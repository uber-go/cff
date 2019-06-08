package predicate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSimplePredicate(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		called := false
		err := Simple(func() { called = true }, true)
		assert.NoError(t, err)
		assert.True(t, called, "function was never called")
	})

	t.Run("false", func(t *testing.T) {
		err := Simple(func() {
			t.Fatal("function must never be called")
		}, false)
		assert.NoError(t, err)
	})
}

func TestSimpleWithContextTask(t *testing.T) {
	require.NoError(t, SimpleWithContextTask())
}

func TestSimpleWithContextPredicate(t *testing.T) {
	require.NoError(t, SimpleWithContextPredicate())
}

func TestSimpleWithContextAndPredicate(t *testing.T) {
	require.NoError(t, SimpleWithContextTaskAndPredicate())
}

func TestExtraDependencies(t *testing.T) {
	require.NoError(t, ExtraDependencies())
}
