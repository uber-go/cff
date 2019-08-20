// +build cff

package badinputs

import (
	"context"

	"go.uber.org/cff"
)

// PredicateReturnsNonbool is a function with a predicate that doesn't return a boolean, instead returns a string.
func PredicateReturnsNonbool() {
	cff.Flow(context.Background(),
		cff.Task(
			func() bool {
				return true
			},
			cff.Predicate(func() string {
				return ""
			}),
		),
	)
}

// PredicateReturnsMultipleValues is a function with a predicate that returns too many values.
func PredicateReturnsMultipleValues() {
	cff.Flow(context.Background(),
		cff.Task(
			func() bool {
				return true
			},
			cff.Predicate(func() (string, bool) {
				return "", true
			}),
		),
	)
}
