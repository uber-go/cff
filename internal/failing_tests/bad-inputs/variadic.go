//go:build cff && failing
// +build cff,failing

package badinputs

import (
	"context"

	"go.uber.org/cff"
)

// Variadic is a flow that has a task that is variadic.
func Variadic() {
	cff.Flow(context.Background(),
		cff.Task(func(s ...string) bool {
			return true
		}),
	)
}

// VariadicPredicate is a flow that has a task whose predicate is variadic.
func VariadicPredicate() {
	cff.Flow(context.Background(),
		cff.Task(
			func() bool {
				return true
			},
			cff.Predicate(func(s ...string) bool {
				return true
			}),
		),
	)
}
