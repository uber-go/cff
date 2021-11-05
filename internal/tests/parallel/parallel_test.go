package parallel

import (
	"context"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	assert.Equal(t, "parallel function panic: sad times", err.Error())
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
	assert.Equal(t, "parallel function panic: sad times", err.Error())
}

func TestMultipleTask(t *testing.T) {
	// Verifies that tasks are executed.
	src := []int{1, 2}
	target := make([]int, 2, 2)
	require.NoError(t, MultipleTask(src, target))
	assert.Equal(t, src, target)
}
