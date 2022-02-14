package shadowedvar

import (
	"context"

	cff2 "go.uber.org/cff"
)

// CtxConflict introduces a variable conflict with ctx to demonstrate that
// CFF2 does not shadow variables.
func CtxConflict(ctx string) (string, error) {
	var result string
	err := cff2.Flow(
		context.Background(),
		cff2.Results(&result),
		cff2.Task(func() (string, error) {
			var hello string
			hello = ctx
			return hello, nil
		}),
	)

	return result, err
}

// PredicateCtxConflict runs the provided function in a task flow if the
// provided boolean is true. This tests if the cff flow works even when the ctx
// variable is shadowed.
func PredicateCtxConflict(f func(), ctx bool) error {
	var s string
	return cff2.Flow(
		context.Background(),
		cff2.Results(&s),
		cff2.Task(
			func() string {
				f()
				return "foo"
			},
			cff2.Predicate(func() bool { return ctx }),
		),
	)
}
