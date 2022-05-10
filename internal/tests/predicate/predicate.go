package predicate

import (
	"context"

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
