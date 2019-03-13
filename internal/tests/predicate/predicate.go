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
