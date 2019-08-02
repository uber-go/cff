// +build cff

package badinputs

import (
	"context"

	"go.uber.org/cff"
)

type foo struct{}
type bar struct{}
type baz struct{}
type qux struct{}
type quux struct{}
type quuz struct{}
type corge struct{}
type grault struct{}

// EarlyResult fails because task with string for return value is unused.
func EarlyResult(ctx context.Context) error {
	var out *foo
	return cff.Flow(
		ctx,
		cff.Results(&out),
		cff.Task(
			func(*foo) *bar {
				return &bar{}
			}),
		cff.Task(func() *foo {
			return &foo{}
		}),
		cff.Task(
			func(*bar) *baz {
				return &baz{}
			}),

		cff.Task(
			func(*baz) string {
				return ""
			}),
	)
}

// EarlyResultDiamond fails because task with int32 return type is unused. Also testing
// diamond formation.
func EarlyResultDiamond(ctx context.Context) error {
	var out *foo
	return cff.Flow(
		ctx,
		cff.Results(&out),
		cff.Task(
			func() *bar {
				return &bar{}
			}),

		cff.Task(func() *foo {
			return &foo{}
		}),
		cff.Task(
			func(*foo, *bar) *baz {
				return &baz{}
			}),

		cff.Task(
			func(*baz) (int64, int32) {
				return int64(1), int32(5)
			}),
		cff.Task(
			func(int64) error {
				return nil
			}),
	)
}

// EarlyResultMultipleFlows fails because task with string for return value is unused.
func EarlyResultMultipleFlows(ctx context.Context) error {
	var out *foo
	return cff.Flow(
		ctx,
		cff.Results(&out),
		cff.Task(
			func(*foo) *bar {
				return &bar{}
			}),
		cff.Task(
			func() *foo {
				return &foo{}
			}),
		cff.Task(
			func(*baz) *quuz {
				return &quuz{}
			}),
		cff.Task(
			func(*bar) *baz {
				return &baz{}
			}),
		cff.Task(
			func(*baz) *corge {
				return &corge{}
			}),
		cff.Task(
			func(*baz) *grault {
				return &grault{}
			}),
		cff.Task(
			func(*baz) *qux {
				return &qux{}
			}),
		cff.Task(
			func(*qux) *quux {
				return &quux{}
			}),
		cff.Task(
			func(*quux) error {
				return nil
			}),
	)
}
