//go:build cff && failing
// +build cff,failing

package badinputs

import (
	"context"

	"go.uber.org/cff"
)

// ContextPredicate is a flow that has a predicate whose context argument is not the first positional argument.
func ContextPredicate() {
	cff.Flow(context.Background(),
		cff.Task(
			func(string) bool {
				return true
			},
			cff.Predicate(func(string, context.Context) bool {
				return true
			}),
		),
	)
}
