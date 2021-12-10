package parallel

import (
	"context"
	"io/fs"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/multierr"
)

func TestTasksAndTask(t *testing.T) {
	var m sync.Map
	require.NoError(t, TasksAndTask(&m))
	var count int
	m.Range(func(_, v_ interface{}) bool {
		count++
		return true
	})
	assert.Equal(t, 3, count)
}

func TestTasksWithError(t *testing.T) {
	require.Error(t, TasksWithError())
}

func TestTasksWithPanic(t *testing.T) {
	err := TasksWithPanic()
	require.Error(t, err)
	assert.Equal(t, "panic: sad times", err.Error())
}

func TestMultipleTasks(t *testing.T) {
	ch := make(chan<- string, 2)
	require.NoError(t, MultipleTasks(ch))
	assert.Len(t, ch, 2)
}

func TestContextErrorBefore(t *testing.T) {
	src := []int{1}
	target := make([]int, 1, 1)
	require.NotEqual(t, src, target)

	ctx, cFn := context.WithCancel(context.Background())
	cFn()

	require.Error(t, ContextErrorBefore(ctx, src, target))
	assert.NotEqual(t, src, target)
}

func TestContextErrorInFlight(t *testing.T) {
	src := []int{1}
	target := make([]int, 1, 1)
	require.NotEqual(t, src, target)

	ctx, cFn := context.WithCancel(context.Background())

	require.Error(t, ContextErrorInFlight(ctx, cFn, src, target))
	assert.NotEqual(t, src, target)
}

func TestTaskWithError(t *testing.T) {
	require.Error(t, TaskWithError())
}

func TestTaskWithPanic(t *testing.T) {
	err := TaskWithPanic()
	require.Error(t, err)
	assert.Equal(t, "panic: sad times", err.Error())
}

func TestMultipleTask(t *testing.T) {
	// Verifies that tasks are executed.
	src := []int{1, 2}
	target := make([]int, 2, 2)
	require.NoError(t, MultipleTask(src, target))
	assert.Equal(t, src, target)
}

func TestContinueOnError(t *testing.T) {
	// Verifies that non-errored tasks are executed.
	src := []int{1, 2}
	target := make([]int, 2, 2)
	err := ContinueOnError(src, target)
	require.Error(t, err)

	// Contains is used instead to verify non-deterministic ordering.
	assert.Contains(t, err.Error(), "sad times")
	assert.Contains(t, err.Error(), "panic: sadder times")

	assert.Equal(t, src, target)
}

func TestContinueOnError_BoolExprAndErrType(t *testing.T) {
	// Verifies that non-errored tasks are executed.
	src := []int{1, 2}
	target := make([]int, 2, 2)
	err := ContinueOnErrorBoolExpr(src, target, func() bool { return true })
	require.Error(t, err)

	// Verifies fidelity to task returned error type.
	errs := multierr.Errors(err)
	require.Len(t, errs, 1)

	assert.ErrorIs(t, errs[0], fs.ErrNotExist)

	assert.Equal(t, src, target)
}

func TestContinueOnError_Cancelled(t *testing.T) {
	src := []int{1}
	target := make([]int, 1, 1)
	require.NotEqual(t, src, target)

	ctx, cFn := context.WithCancel(context.Background())
	cFn()

	require.Error(t, ContinueOnErrorCancelled(ctx, src, target))
	assert.NotEqual(t, src, target)
}

func TestContinueOnError_CancelledDuring(t *testing.T) {
	src := []int{1}
	target := make([]int, 1, 1)
	require.NotEqual(t, src, target)

	ctx, cFn := context.WithCancel(context.Background())

	require.Error(t, ContinueOnErrorCancelledDuring(ctx, cFn, src, target))
	// Even with ContinueOnError tasks with cancelled contexts should not be
	// run by the scheduler.
	assert.NotEqual(t, src, target)
}

func TestSlice(t *testing.T) {
	src := []int{1, 2}
	target := make([]int, 2, 2)
	assert.NotEqual(t, src, target)

	require.NoError(t, Slice(src, target))

	assert.Equal(t, src, target)
}

func TestMultiple(t *testing.T) {
	src := []int{1, 2}
	targetA := make([]int, 2, 2)
	targetB := make([]int, 2, 2)
	assert.NotEqual(t, src, targetA)
	assert.NotEqual(t, src, targetB)

	require.NoError(t, SliceMultiple(src, src, targetA, targetB))

	assert.Equal(t, src, targetA)
	assert.Equal(t, src, targetB)
}

func TestSliceError(t *testing.T) {
	src := []int{1, 2}
	target := make([]int, 2)
	assert.NotEqual(t, src, target)

	err := SliceError(src, target)
	require.Error(t, err)

	assert.Equal(t, "sad times", err.Error())
	assert.NotEqual(t, src, target)
}

func TestSlicePanic(t *testing.T) {
	src := []int{1, 2}
	target := make([]int, 2, 2)
	assert.NotEqual(t, src, target)

	err := SlicePanic(src, target)
	require.Error(t, err)

	assert.Equal(t, "panic: sad times", err.Error())
	assert.NotEqual(t, src, target)
}

func TestSliceContinueOnError(t *testing.T) {
	src := []string{"copy", "error", "panic", "me"}
	target := make([]string, 4, 4)
	assert.NotEqual(t, src, target)

	err := SliceContinueOnError(src, target)
	require.Error(t, err)

	assert.Contains(t, err.Error(), "sad times")
	assert.Contains(t, err.Error(), "panic: sadder times")

	assert.Equal(t, []string{"copy", "", "", "me"}, target)
}

func TestMap(t *testing.T) {
	src := map[string]int{
		"test": 0,
		"one":  1,
	}
	keys, values := make([]string, 2), make([]int, 2)
	require.NoError(t, AssignMapItems(src, keys, values, false))

	assert.Equal(t, []string{"test", "one"}, keys)
	assert.Equal(t, []int{0, 1}, values)
}

func TestMapError(t *testing.T) {
	src := map[string]int{
		"error": 0,
	}
	var keys []string
	var values []int

	err := AssignMapItems(src, keys, values, false)
	require.Error(t, err)
	assert.Equal(t, "sad times", err.Error())
}

func TestMapPanic(t *testing.T) {
	src := map[string]int{
		"test":  0,
		"panic": 1,
	}
	keys, values := make([]string, 1), make([]int, 1)

	err := AssignMapItems(src, keys, values, false)
	require.Error(t, err)
	assert.EqualError(t, err, "panic: sadder times")

	assert.Equal(t, []string{"test"}, keys)
	assert.Equal(t, []int{0}, values)
}

func TestMapContinueOnError(t *testing.T) {
	src := map[string]int{
		"copy":  0,
		"error": 2,
		"panic": 3,
		"me":    1,
	}

	keys, values := make([]string, 2), make([]int, 2)

	err := AssignMapItems(src, keys, values, true)
	require.Error(t, err)

	// Using assert.Contains here because the returned error is non-deterministic.
	assert.Contains(t, err.Error(), "sad times")
	assert.Contains(t, err.Error(), "panic: sadder times")

	assert.Equal(t, []string{"copy", "me"}, keys)
	assert.Equal(t, []int{0, 1}, values)
}
