package parallel

import (
	"context"
	"errors"
	"fmt"
	"io/fs"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/multierr"
	"golang.org/x/exp/maps"
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
	src := []string{"1", "2"}
	target := make([]string, len(src))
	assert.NotEqual(t, src, target)

	require.NoError(t, AssignSliceItems(src, target, false))

	assert.Equal(t, src, target)
}

func TestMultiple(t *testing.T) {
	src := []int{1, 2}
	targetA := make([]int, len(src))
	targetB := make([]int, len(src))
	assert.NotEqual(t, src, targetA)
	assert.NotEqual(t, src, targetB)

	require.NoError(t, SliceMultiple(src, src, targetA, targetB))

	assert.Equal(t, src, targetA)
	assert.Equal(t, src, targetB)
}

func TestSliceNoIndex(t *testing.T) {
	src := []int{0, 1}

	targetA := make([]int, len(src))
	targetB := make([]int, len(src))
	assert.NotEqual(t, src, targetA)
	assert.NotEqual(t, src, targetB)

	require.NoError(t, SliceNoIndex(src, src, targetA, targetB))

	assert.Equal(t, src, targetA)
	assert.Equal(t, src, targetB)
}

func TestSliceWrapped(t *testing.T) {
	src := []int{1, 2}
	target := make([]int, len(src))
	require.NoError(t, SliceWrapped(src, target))
	assert.Equal(t, src, target)
}

func TestSliceError(t *testing.T) {
	src := []string{"1", "error"}
	target := make([]string, len(src))
	assert.NotEqual(t, src, target)

	err := AssignSliceItems(src, target, false)
	require.Error(t, err)

	assert.Equal(t, "sad times", err.Error())
	assert.Equal(t, "error", target[1])
}

func TestSlicePanic(t *testing.T) {
	src := []string{"1", "panic"}
	target := make([]string, len(src))
	assert.NotEqual(t, src, target)

	err := AssignSliceItems(src, target, false)
	require.Error(t, err)

	assert.Equal(t, "panic: sadder times", err.Error())
	assert.Equal(t, "panic", target[1])
}

func TestSliceContinueOnError(t *testing.T) {
	src := []string{"copy", "error", "panic", "me"}
	target := make([]string, len(src))
	assert.NotEqual(t, src, target)

	err := AssignSliceItems(src, target, true)
	require.Error(t, err)

	assert.Contains(t, err.Error(), "sad times")
	assert.Contains(t, err.Error(), "panic: sadder times")
	assert.Equal(t, []string{"copy", "error", "panic", "me"}, target)
}

func TestSliceEnd(t *testing.T) {
	var src, target []int
	errSadTimes := errors.New("sad times")
	errSadderTimes := errors.New("panic: sadder times")
	tests := []struct {
		desc       string
		sliceEndFn func()
		src        []int
		wantTarget []int
		wantErr    error
	}{
		{
			desc: "success",
			sliceEndFn: func() {
				assert.Equal(t, src, target)
			},
			src:        []int{1, 2, 3},
			wantTarget: []int{1, 2, 3},
		},
		{
			desc: "panic",
			sliceEndFn: func() {
				assert.Equal(t, src, target)
				panic("sadder times")
			},
			src:        []int{1, 2, 3},
			wantTarget: []int{1, 2, 3},
			wantErr:    errSadderTimes,
		},
		{
			desc: "not called on error",
			sliceEndFn: func() {
				t.Fatal("SliceEnd shouldn't run after a error")
			},
			src:        []int{-1},
			wantTarget: []int{},
			wantErr:    errSadTimes,
		},
		{
			desc: "not called on panic",
			sliceEndFn: func() {
				t.Fatal("SliceEnd shouldn't run after a panic")
			},
			src:        []int{-2},
			wantTarget: []int{},
			wantErr:    errSadderTimes,
		},
	}

	assignItemsFn := func(idx, val int) error {
		switch val {
		case -1:
			return errSadTimes
		case -2:
			panic("sadder times")
		}
		target[idx] = val
		return nil
	}

	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			src = tt.src
			target = make([]int, len(tt.wantTarget))
			err := SliceEnd(src, assignItemsFn, tt.sliceEndFn)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantTarget, target)
		})
	}
}

func TestSliceEndWithErr(t *testing.T) {
	var src, target []int
	errSadTimes := errors.New("sad times")
	errSadderTimes := errors.New("panic: sadder times")
	tests := []struct {
		desc       string
		sliceEndFn func() error
		src        []int
		wantTarget []int
		wantErr    error
	}{
		{
			desc: "success",
			sliceEndFn: func() error {
				assert.Equal(t, src, target)
				return nil
			},
			src:        []int{1, 2, 3},
			wantTarget: []int{1, 2, 3},
		},
		{
			desc: "error",
			sliceEndFn: func() error {
				assert.Equal(t, src, target)
				return errSadTimes
			},
			src:        []int{1, 2, 3},
			wantTarget: []int{1, 2, 3},
			wantErr:    errSadTimes,
		},
		{
			desc: "panic",
			sliceEndFn: func() error {
				assert.Equal(t, src, target)
				panic("sadder times")
				return nil
			},
			src:        []int{1, 2, 3},
			wantTarget: []int{1, 2, 3},
			wantErr:    errSadderTimes,
		},
		{
			desc: "not called on error",
			sliceEndFn: func() error {
				t.Fatal("SliceEnd shouldn't run after a error")
				return nil
			},
			src:        []int{-1},
			wantTarget: []int{},
			wantErr:    errSadTimes,
		},
		{
			desc: "not called on panic",
			sliceEndFn: func() error {
				t.Fatal("SliceEnd shouldn't run after a panic")
				return nil
			},
			src:        []int{-2},
			wantTarget: []int{},
			wantErr:    errSadderTimes,
		},
	}

	assignItemsFn := func(idx, val int) error {
		switch val {
		case -1:
			return errSadTimes
		case -2:
			panic("sadder times")
		}
		target[idx] = val
		return nil
	}

	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			src = tt.src
			target = make([]int, len(tt.wantTarget))
			err := SliceEndWithErr(src, assignItemsFn, tt.sliceEndFn)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantTarget, target)
		})
	}
}

func TestSliceEndWithCtx(t *testing.T) {
	var src, target []int
	errSadTimes := errors.New("sad times")
	errSadderTimes := errors.New("panic: sadder times")
	tests := []struct {
		desc       string
		sliceEndFn func(ctx context.Context)
		src        []int
		wantTarget []int
		wantErr    error
	}{
		{
			desc: "success",
			sliceEndFn: func(ctx context.Context) {
				assert.NotNil(t, ctx)
				assert.Equal(t, src, target)
			},
			src:        []int{1, 2, 3},
			wantTarget: []int{1, 2, 3},
		},
		{
			desc: "panic",
			sliceEndFn: func(ctx context.Context) {
				assert.NotNil(t, ctx)
				assert.Equal(t, src, target)
				panic("sadder times")
			},
			src:        []int{1, 2, 3},
			wantTarget: []int{1, 2, 3},
			wantErr:    errSadderTimes,
		}, {
			desc: "not called on error",
			sliceEndFn: func(context.Context) {
				t.Fatal("SliceEnd shouldn't run after a error")
			},
			src:        []int{-1},
			wantTarget: []int{},
			wantErr:    errSadTimes,
		},
		{
			desc: "not called on panic",
			sliceEndFn: func(context.Context) {
				t.Fatal("SliceEnd shouldn't run after a panic")
			},
			src:        []int{-2},
			wantTarget: []int{},
			wantErr:    errSadderTimes,
		},
	}

	assignItemsFn := func(idx, val int) error {
		switch val {
		case -1:
			return errSadTimes
		case -2:
			panic("sadder times")
		}
		target[idx] = val
		return nil
	}

	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			src = tt.src
			target = make([]int, len(tt.wantTarget))
			err := SliceEndWithCtx(src, assignItemsFn, tt.sliceEndFn)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantTarget, target)
		})
	}
}

func TestSliceEndWithCtxAndErr(t *testing.T) {
	var src, target []int
	errSadTimes := errors.New("sad times")
	errSadderTimes := errors.New("panic: sadder times")
	tests := []struct {
		desc       string
		sliceEndFn func(ctx context.Context) error
		src        []int
		wantTarget []int
		wantErr    error
	}{
		{
			desc: "success",
			sliceEndFn: func(ctx context.Context) error {
				assert.NotNil(t, ctx)
				assert.Equal(t, src, target)
				return nil
			},
			src:        []int{1, 2, 3},
			wantTarget: []int{1, 2, 3},
		},
		{
			desc: "panic",
			sliceEndFn: func(ctx context.Context) error {
				assert.NotNil(t, ctx)
				assert.Equal(t, src, target)
				panic("sadder times")
				return nil
			},
			src:        []int{1, 2, 3},
			wantTarget: []int{1, 2, 3},
			wantErr:    errSadderTimes,
		}, {
			desc: "not called on error",
			sliceEndFn: func(context.Context) error {
				t.Fatal("SliceEnd shouldn't run after a error")
				return nil
			},
			src:        []int{-1},
			wantTarget: []int{},
			wantErr:    errSadTimes,
		},
		{
			desc: "not called on panic",
			sliceEndFn: func(context.Context) error {
				t.Fatal("SliceEnd shouldn't run after a panic")
				return nil
			},
			src:        []int{-2},
			wantTarget: []int{},
			wantErr:    errSadderTimes,
		},
	}

	assignItemsFn := func(idx, val int) error {
		switch val {
		case -1:
			return errSadTimes
		case -2:
			panic("sadder times")
		}
		target[idx] = val
		return nil
	}

	for _, tt := range tests {
		t.Run(t.Name(), func(t *testing.T) {
			src = tt.src
			target = make([]int, len(tt.wantTarget))
			err := SliceEndWithCtxAndErr(src, assignItemsFn, tt.sliceEndFn)
			assert.Equal(t, tt.wantErr, err)
			assert.Equal(t, tt.wantTarget, target)
		})
	}
}

func TestMap(t *testing.T) {
	src := map[string]int{
		"test": 0,
		"one":  1,
	}
	keys, values := make([]string, len(src)), make([]int, len(src))
	require.NoError(t, AssignMapItems(src, keys, values, false))

	assert.Equal(t, []string{"test", "one"}, keys)
	assert.Equal(t, []int{0, 1}, values)
}

func TestMapError(t *testing.T) {
	src := map[string]int{
		"error": 0,
	}

	err := AssignMapItems(src, nil, nil, false)
	require.Error(t, err)
	assert.Equal(t, "sad times", err.Error())
}

func TestMapPanic(t *testing.T) {
	src := map[string]int{
		"panic": 1,
	}
	err := AssignMapItems(src, nil, nil, false)
	require.Error(t, err)
	assert.EqualError(t, err, "panic: sadder times")
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

func TestMapEnd(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

	t.Run("success", func(t *testing.T) {
		// Push keys from the map into an unbuffered channel.
		//
		// Reading on the channel will deadlock if the end function is
		// never called.
		keyc := make(chan string)
		go ForEachMapItem(m, func(k string, _ int) {
			keyc <- k
		}, func() {
			close(keyc)
			// If this does not get called,
			// the range on keyc will deadlock
			// because the channel is never closed.
		})

		keys := make([]string, 0, len(m))
		for k := range keyc {
			keys = append(keys, k)
		}

		assert.ElementsMatch(t, maps.Keys(m), keys)
	})

	t.Run("panic", func(t *testing.T) {
		called := false
		defer func() {
			assert.True(t, called, "MapEnd never called")
		}()

		err := ForEachMapItem(m, func(string, int) {
		}, func() {
			called = true
			panic("great sadness")
		})

		require.Error(t, err)
		assert.Contains(t, err.Error(), "panic")
	})

	t.Run("skipped on panic", func(t *testing.T) {
		err := ForEachMapItem(m, func(s string, _ int) {
			if s == "c" {
				panic("great sadness")
			}
		}, func() {
			t.Error("This function should not be called")
		})
		require.Error(t, err)
		assert.Contains(t, err.Error(), "panic")
	})
}

func TestMapEnd_ReturnsError(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

	t.Run("success", func(t *testing.T) {
		calledEnd := false
		defer func() {
			assert.True(t, calledEnd, "MapEnd not called")
		}()

		err := ForEachMapItemError(m, func(k string, _ int) error {
			switch k {
			case "a", "b", "c", "d":
				return nil
			default:
				return fmt.Errorf("unexpected key %q", k)
			}
		}, func() error {
			calledEnd = true
			return nil
		})
		require.NoError(t, err)
	})

	t.Run("errors", func(t *testing.T) {
		giveErr := errors.New("great sadness")

		err := ForEachMapItemError(m, func(string, int) error {
			// doesn't matter
			return nil
		}, func() error {
			return giveErr
		})
		assert.ErrorIs(t, err, giveErr)
	})

	t.Run("skipped on error", func(t *testing.T) {
		giveErr := errors.New("great sadness")

		err := ForEachMapItemError(m, func(s string, _ int) error {
			if s == "c" {
				return giveErr
			}
			return nil
		}, func() error {
			t.Errorf("this function should not be called")
			return nil
		})
		assert.ErrorIs(t, err, giveErr)
	})
}

func TestMapEnd_HasContext(t *testing.T) {
	m := map[string]int{"a": 1, "b": 2, "c": 3, "d": 4}

	t.Run("success", func(t *testing.T) {
		calledEnd := false
		defer func() {
			assert.True(t, calledEnd, "MapEnd not called")
		}()

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := ForEachMapItemContext(ctx, m,
			func(ctx context.Context, k string, _ int) {
				_, ok := ctx.Deadline()
				assert.True(t, ok, "task %q received no deadline", k)
			}, func(ctx context.Context) {
				_, ok := ctx.Deadline()
				assert.True(t, ok, "MapEnd received no deadline")
				calledEnd = true
			})
		require.NoError(t, err)
	})

	t.Run("expired context", func(t *testing.T) {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := ForEachMapItemContext(ctx, m,
			func(ctx context.Context, k string, _ int) {
				if k == "c" {
					cancel()
				}
			}, func(context.Context) {
				t.Errorf("this function should not be called")
			})
		assert.ErrorIs(t, err, context.Canceled)
	})
}
