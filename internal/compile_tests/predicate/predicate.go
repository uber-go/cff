// +build cff

package predicate

import (
	"context"

	"go.uber.org/cff"
)

// ExampleFlowWithPredicates provides a flow that deploys multiple
// cff.Predicates. This flow is compiled to test the CFF compiler's
// internal state.
func ExampleFlowWithPredicates(f func(), pred bool) error {
	var s string
	var b bool
	return cff.Flow(
		context.Background(),
		cff.Results(&s, &b),
		cff.Task(
			func() string {
				f()
				return "foo"
			},
			cff.Predicate(func() bool { return pred }),
		),
		cff.Task(
			func() bool {
				f()
				return true
			},
			cff.Predicate(func() bool { return !pred }),
		),
	)
}
