// +build cff

package badinputs

import (
	"go.uber.org/cff"
	"context"
)

// PredicateParams is a flow that has a predicate that has bad arguments.
func PredicateParams() {
	cff.Flow(context.Background(),
		cff.Task(
			func(string) bool {
				return true
			},
			cff.Predicate(nil),
		),
	)
}
