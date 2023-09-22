package panic

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/cff"
)

func TestCatchesPanicParallel(t *testing.T) {
	var p Panicker
	err := p.FlowPanicsParallel()
	assert.ErrorContains(t, err, "panic: panic\ngoroutine")
	var panicError *cff.PanicError
	require.ErrorAs(t, err, &panicError, "error returned should be a cff.PanicError")
	assert.Equal(t, "panic", panicError.Value, "PanicError.Value should be recovered value")
	stacktrace := string(panicError.Stacktrace)
	assert.Contains(t, stacktrace, "panic({", "panic should be included in the stack trace")
	assert.Contains(t, stacktrace, ".FlowPanicsParallel.func", "function that panicked should be in the stack")
}

func TestCatchesPanicSerial(t *testing.T) {
	var p Panicker
	err := p.FlowPanicsSerial()
	assert.ErrorContains(t, err, "panic: panic\ngoroutine")
	var panicError *cff.PanicError
	require.ErrorAs(t, err, &panicError, "error returned should be a cff.PanicError")
	assert.Equal(t, "panic", panicError.Value, "PanicError.Value should be recovered value")
	stacktrace := string(panicError.Stacktrace)
	assert.Contains(t, stacktrace, "panic({", "panic should be included in the stack trace")
	assert.Contains(t, stacktrace, ".FlowPanicsSerial.func", "function that panicked should be in the stack")
}
