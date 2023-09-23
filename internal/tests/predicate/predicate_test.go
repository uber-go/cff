package predicate

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/cff"
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

func TestMultiplePredicates(t *testing.T) {
	require.NoError(t, MultiplePredicates())
}

func TestPanicRecovered(t *testing.T) {
	var err error
	require.NotPanics(
		t,
		func() {
			err = Panicked()
		},
	)
	var panicError *cff.PanicError
	require.ErrorAs(t, err, &panicError, "error returned should be a cff.PanicError")
	assert.Equal(t, "sad times", panicError.Value, "PanicError.Value should be recovered value")
	stacktrace := string(panicError.Stacktrace)
	assert.Contains(t, stacktrace, "panic({", "panic should be included in the stack trace")
	assert.Contains(t, stacktrace, ".Panicked.func", "function that panicked should be in the stack")
}

func TestPanicFallback(t *testing.T) {
	var (
		s   string
		err error
	)
	require.NotPanics(
		t,
		func() {
			s, err = PanickedWithFallback()
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, s, "predicate-fallback")
}
