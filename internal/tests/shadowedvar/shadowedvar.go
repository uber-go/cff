//go:build cff
// +build cff

package shadowedvar

import (
	"context"

	cff2 "go.uber.org/cff"
)

// CtxConflict introduces a variable conflict with ctx to demonstrate that
// CFF does not shadow variables.
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

// CtxConflictParallel introduces a variable conflict with ctx within cff.Parallel Task
// to demonstrate that CFF does not shadow variables.
func CtxConflictParallel(ctx string) (string, string, error) {
	var result1 string
	var result2 string
	err := cff2.Parallel(
		context.Background(),
		cff2.Task(func() {
			result1 = ctx
		}),
		cff2.Task(func() {
			result2 = ctx
		}),
	)
	return result1, result2, err
}

// CtxConflictSlice introduces a variable conflict with ctx within cff.Slice function
// to demonstrate that CFF does not shadow variables.
func CtxConflictSlice(ctx string, target []string) error {
	return cff2.Parallel(
		context.Background(),
		cff2.Concurrency(2),
		cff2.Slice(
			func(idx int, val string) error {
				target[idx] = ctx + val
				return nil
			},
			target,
		),
	)
}

// CtxConflictMap introduces a variable conflict with ctx within cff.Map function
// to demonstrate that CFF does not shadow variables.
func CtxConflictMap(ctx int, input map[int]int) ([]int, error) {
	slice := make([]int, len(input))
	err := cff2.Parallel(
		context.Background(),
		cff2.Concurrency(2),
		cff2.Map(
			func(key int, val int) {
				slice[key] = ctx + val
			},
			input,
		),
	)
	return slice, err
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
