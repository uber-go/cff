// +build cff

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
		cff.Tasks(
			func(int) string { return "foo" },
			func(int) t1 { return t1{} },
			func() t2 { return t2{} },
		),
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
