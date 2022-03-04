package parallel

import (
	"context"
	"errors"
	"os"
	"sync"

	"go.uber.org/cff"
)

// TasksAndTask runs cff.Tasks and cff.Task in parallel to populate the
// provided map.
func TasksAndTask(m *sync.Map) error {
	return cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Tasks(
			func() {
				m.Store("foo", "bar")
			},
			func(_ context.Context) {
				m.Store("fiz", "buzz")
			},
		),
		cff.Task(
			func(_ context.Context) {
				m.Store("go", "lang")
			},
		),
	)
}

// TasksWithError runs a parallel cff.Tasks that errors.
func TasksWithError() error {
	return cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Tasks(
			func() error {
				return errors.New("sad times")
			},
		),
	)
}

// TasksWithPanic runs a parallel cff.Tasks that panics.
func TasksWithPanic() error {
	return cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Tasks(
			func() {
				panic("sad times")
			},
		),
	)
}

// MultipleTasks runs multiple cff.Tasks in parallel to populate the provided
// channel.
func MultipleTasks(c chan<- string) error {
	send := func(_ context.Context) error {
		c <- "send"
		return nil
	}
	return cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Tasks(
			func() {
				c <- "multiple"
			},
		),
		cff.Tasks(
			send,
		),
	)
}

// ContextErrorBefore runs a cff.Tasks function to test that the function is
// not run if the context errors before scheduler execution.
func ContextErrorBefore(ctx context.Context, src, target []int) error {
	return cff.Parallel(
		ctx,
		cff.Concurrency(2),
		cff.Tasks(
			func() {
				target[0] = src[0]
			},
		),
	)
}

// ContextErrorInFlight runs a cff.Tasks function to test that the function is
// not run if the context errors during scheduler execution.
func ContextErrorInFlight(ctx context.Context, cancel func(), src, target []int) error {
	blocker := make(chan struct{})
	return cff.Parallel(
		ctx,
		cff.Concurrency(2),
		cff.Tasks(
			// Busy both workers during context cancellation so that the
			// third function is processed after cancellation.
			func() {
				cancel()
				close(blocker)
			},
			func() {
				<-blocker
			},
			func() {
				target[0] = src[0]
			},
		),
	)
}

// TaskWithError runs a parallel cff.Task that errors.
func TaskWithError() error {
	return cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Task(
			func() error {
				return errors.New("sad times")
			},
		),
	)
}

// TaskWithPanic runs a parallel cff.task that panics.
func TaskWithPanic() error {
	return cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Task(
			func() {
				panic("sad times")
			},
		),
	)
}

// MultipleTask runs multiple cff.Task in parallel to populate the provided
// slice.
func MultipleTask(src, target []int) error {
	send := func(ctx context.Context) error {
		_, _ = ctx.Deadline()
		target[1] = src[1]
		return nil
	}
	return cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Task(
			func() {
				target[0] = src[0]
			},
		),
		cff.Task(
			send,
		),
	)
}

// ContinueOnError runs multiple parallel tasks through errors in the
// scheduler. Despite tasks that error, the contents of the src slice should
// be copied into the target slice.
func ContinueOnError(src []int, target []int) error {
	blockerA := make(chan struct{})
	blockerB := make(chan struct{})
	return cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.ContinueOnError(true),
		cff.Tasks(
			func(_ context.Context) error {
				close(blockerA)
				return errors.New("sad times")
			},
			func() {
				// Ensure erroring task has run before unblocking this copying
				// task.
				<-blockerA
				target[0] = src[0]
			},
		),
		cff.Task(
			func() {
				// Don't panic until after the error task is run.
				<-blockerA
				close(blockerB)
				panic("sadder times")
			},
		),
		cff.Task(
			func() {
				// Ensure panicing task has run before unblocking this copying
				// task.
				<-blockerB
				target[1] = src[1]
			},
		),
	)
}

// ContinueOnErrorBoolExpr parameterizes the ContinueOnError Option
// with a boolean expression. Despite tasks that error, the contents of
// the src slice should be copied into the target slice. The error returned
// by this function is checked to verify that it's underlying type is
// unchanged.
func ContinueOnErrorBoolExpr(src, target []int, fn func() bool) error {
	blockerA := make(chan struct{})
	blockerB := make(chan struct{})
	return cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.ContinueOnError(fn()),
		cff.Tasks(
			func(_ context.Context) error {
				close(blockerA)
				// Use the error from os.Open so we can also assert with
				// errors.Is for fs.ErrNotExist.
				_, err := os.Open("non-existing")
				return err
			},
			func() {
				// Ensure erroring task has run before unblocking this copying
				// task.
				<-blockerA
				target[0] = src[0]
				close(blockerB)
			},
		),
		cff.Task(
			func() {
				// Given 2 workers and blocking on two earlier declared tasks,
				// this copying task runs after the error task has been
				// processed by the scheduler.
				<-blockerB
				target[1] = src[1]
			},
		),
	)
}

// ContinueOnErrorCancelled runs a Parallel with a cancelled context.
func ContinueOnErrorCancelled(ctx context.Context, src []int, target []int) error {
	return cff.Parallel(
		ctx,
		cff.Concurrency(2),
		cff.ContinueOnError(true),
		cff.Task(
			func() {
				target[0] = src[0]
			},
		),
	)
}

// ContinueOnErrorCancelledDuring runs a parallel with a context that is
// cancelled during scheduler operation.
func ContinueOnErrorCancelledDuring(ctx context.Context, cancel func(), src []int, target []int) error {
	blocker := make(chan struct{})
	return cff.Parallel(
		ctx,
		cff.Concurrency(2),
		cff.ContinueOnError(true),
		// Busy both workers while the context is cancelled so that the third
		// function is processed after cancellation.
		cff.Tasks(
			func() {
				cancel()
				close(blocker)
			},
			func() {
				<-blocker
			},
		),
		cff.Task(
			func() {
				target[0] = src[0]
			},
		),
	)
}

// SliceMultiple runs multiple cff.Slices to populate the provided slice.
// Slice will panic if len(target) < len(src).
func SliceMultiple(srcA, srcB, targetA, targetB []int) error {
	return cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Slice(
			func(idx int, val int) error {
				targetA[idx] = val
				return nil
			},
			srcA,
		),
		cff.Slice(
			func(_ context.Context, idx int, val int) {
				targetB[idx] = val
			},
			srcB,
		),
	)
}

// AssignSliceItems runs cff.Slice in parallel to populate the provided slices.
// AssignSliceItems expects len(target) >= len(src).
func AssignSliceItems(src, target []string, keepgoing bool) error {
	return cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.ContinueOnError(keepgoing),
		cff.Slice(
			func(idx int, val string) error {
				switch val {
				case "error":
					return errors.New("sad times")
				case "panic":
					panic("sadder times")
				default:
					target[idx] = val
					return nil
				}
			},
			src,
		),
	)
}

// SliceEnd runs cff.Slice in parallel and calls sliceEndFn after all items
// in the slice have finished.
func SliceEnd(src []int, sliceFn func(idx, val int) error, sliceEndFn func()) (err error) {
	err = cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Slice(
			sliceFn,
			src,
			cff.SliceEnd(sliceEndFn),
		),
	)
	return err
}

// SliceEndWithErr runs cff.Slice in parallel and calls an erroring sliceEndFn
// after all items in the slice have finished.
func SliceEndWithErr(src []int, sliceFn func(idx, val int) error, sliceEndFn func() error) (err error) {
	err = cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Slice(
			sliceFn,
			src,
			cff.SliceEnd(sliceEndFn),
		),
	)
	return err
}

// SliceEndWithCtx runs cff.Slice in parallel and calls sliceEndFn after all items
// in the slice have finished.
func SliceEndWithCtx(src []int, sliceFn func(idx, val int) error, sliceEndFn func(context.Context)) (err error) {
	err = cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Slice(
			sliceFn,
			src,
			cff.SliceEnd(sliceEndFn),
		),
	)
	return err
}

// SliceEndWithCtxAndErr runs cff.Slice in parallel and calls an erroring sliceEndFn
// after all items in the slice have finished.
func SliceEndWithCtxAndErr(src []int, sliceFn func(idx, val int) error, sliceEndFn func(context.Context) error) (err error) {
	err = cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.Slice(
			sliceFn,
			src,
			cff.SliceEnd(sliceEndFn),
		),
	)
	return err
}

// AssignMapItems runs cff.Map in parallel to populate the provided slices.
func AssignMapItems(src map[string]int, keys []string, values []int, keepgoing bool) error {
	return cff.Parallel(
		context.Background(),
		cff.Concurrency(2),
		cff.ContinueOnError(keepgoing),
		cff.Map(
			func(key string, val int) error {
				switch key {
				case "error":
					return errors.New("sad times")
				case "panic":
					panic("sadder times")
				default:
					keys[val] = key
					values[val] = val
					return nil
				}
			},
			src,
		),
	)
}
