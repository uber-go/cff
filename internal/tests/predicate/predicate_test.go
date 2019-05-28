package predicate_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	predicate_gen "go.uber.org/cff/internal/tests/predicate_gen"
)

func TestSimplePredicate(t *testing.T) {
	t.Run("true", func(t *testing.T) {
		called := false
		err := predicate_gen.Simple(func() { called = true }, true)
		assert.NoError(t, err)
		assert.True(t, called, "function was never called")
	})

	t.Run("false", func(t *testing.T) {
		err := predicate_gen.Simple(func() {
			t.Fatal("function must never be called")
		}, false)
		assert.NoError(t, err)
	})
}

func TestSimpleWithContextTask(t *testing.T) {
	require.NoError(t, predicate_gen.SimpleWithContextTask())
}

func TestSimpleWithContextPredicate(t *testing.T) {
	require.NoError(t, predicate_gen.SimpleWithContextPredicate())
}

func TestSimpleWithContextAndPredicate(t *testing.T) {
	require.NoError(t, predicate_gen.SimpleWithContextTaskAndPredicate())
}

func TestExtraDependencies(t *testing.T) {
	require.NoError(t, predicate_gen.ExtraDependencies())
}
