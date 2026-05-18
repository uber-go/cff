//go:build cff
// +build cff

package predicate

import (
	"context"
	"time"

	"go.uber.org/cff"
)

// Simple runs the provided function in a task flow if the provided boolean
// is true.
func Simple(f func(), pred bool) error {
	var s string
	return cff.Flow(
		context.Background(),
		cff.Results(&s),
		cff.Task(
			func() string {
				f()
				return "foo"
			},
			cff.Predicate(func() bool { return pred }),
		),
	)
}

// SimpleWithContextTask is a task flow which checks that context can be passed into Task w/out
// errors.
func SimpleWithContextTask() error {
	var s string
	return cff.Flow(
		context.Background(),
		cff.Results(&s),
		cff.Params(int64(2)),
		cff.Task(
			func(ctx context.Context) string {
				return "foo"
			},
			cff.Predicate(
				func(int64) bool {
					return false
				}),
		),
	)
}

// SimpleWithContextPredicate is a task flow which checks that context can be passed into
// Predicate but not Task.
func SimpleWithContextPredicate() error {
	var s string
	return cff.Flow(
		context.Background(),
		cff.Results(&s),
		cff.Params(int64(2)),
		cff.Task(
			func() string {
				return "foo"
			},
			cff.Predicate(
				func(context.Context, int64) bool {
					return false
				}),
		),
	)
}

// SimpleWithContextTaskAndPredicate is a task flow which checks that context can be passed into
// Predicate and Task.
func SimpleWithContextTaskAndPredicate() error {
	var s string
	return cff.Flow(
		context.Background(),
		cff.Results(&s),
		cff.Params(int64(2)),
		cff.Task(
			func(ctx context.Context) string {
				return "foo"
			},
			cff.Predicate(
				func(context.Context, int64) bool {
					return false
				}),
		),
	)
}

// ExtraDependencies is a task flow where the predicate has more dependencies
// than the task.
func ExtraDependencies() error {
	type t1 struct{}
	type t2 struct{}
	type t3 struct{}

	var out t3
	return cff.Flow(
		context.Background(),
		cff.Params(int(42)),
		cff.Results(&out),
		cff.Task(
			func(int) t1 { return t1{} }),
		cff.Task(
			func() t2 { return t2{} }),
		cff.Task(
			func(t2) t3 { return t3{} },
			cff.Predicate(
				func(int, t1) bool {
					return true
				},
			),
		),
	)
}

// MultiplePredicates is a task flow which checks that the outputs of multiple
// predicates can be distinguished.
func MultiplePredicates() error {
	var s string
	var b bool
	return cff.Flow(
		context.Background(),
		cff.Results(&s, &b),
		cff.Task(
			func() string {
				return "foo"
			},
			cff.Predicate(func() bool { return true }),
		),
		cff.Task(
			func() bool {
				return true
			},
			cff.Predicate(func() bool { return false }),
		),
	)
}

// Panicked is a task flow that contains a task predicate that panics.
func Panicked() error {
	var s string
	return cff.Flow(
		context.Background(),
		cff.Results(&s),
		cff.Task(
			func(ctx context.Context) string {
				return "foo"
			},
			cff.Predicate(
				func() bool {
					panic("sad times")
					return true
				},
			),
		),
	)
}

// PanickedWithFallback is a flow that runs a panicing task predicate with a
// fallback.
func PanickedWithFallback() (string, error) {
	var s string
	err := cff.Flow(
		context.Background(),
		cff.Results(&s),
		cff.Task(
			func(ctx context.Context) (string, error) {
				return "foo", nil
			},
			cff.Predicate(
				func() bool {
					panic("sad times")
					return true
				},
			),
			cff.FallbackWith("predicate-fallback"),
		),
	)
	return s, err
}

// BlockingInputs builds a flow that probes whether cff.Predicate
// short-circuits the scheduler's wait on the predicated task's input
// dependencies.
//
// Shape: a slow source feeds the predicated task; a fast source feeds
// the predicate. The predicate returns false (so the task body is
// skipped via FallbackWith). A downstream consumer measures the
// elapsed wall-clock time from flow start.
//
// If predicates short-circuit input deps:
//   elapsed ≈ 0 (slow path is not on the critical path)
// If predicates only short-circuit the task body:
//   elapsed ≈ slowDelay (scheduler still waits on slowOut)
func BlockingInputs(slowDelay time.Duration) (time.Duration, error) {
	type slowOut struct{}
	type fastOut struct{}
	type skippedOut struct{}

	var elapsed time.Duration
	start := time.Now()
	err := cff.Flow(
		context.Background(),
		cff.Results(&elapsed),
		cff.Task(func() slowOut {
			time.Sleep(slowDelay)
			return slowOut{}
		}),
		cff.Task(func() fastOut { return fastOut{} }),
		cff.Task(
			func(slowOut) (skippedOut, error) { return skippedOut{}, nil },
			cff.Predicate(func(fastOut) bool { return false }),
			cff.FallbackWith(skippedOut{}),
		),
		cff.Task(func(skippedOut) time.Duration {
			return time.Since(start)
		}),
	)
	return elapsed, err
}
