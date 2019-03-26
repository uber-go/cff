// +build cff

package badinputs

import (
	"go.uber.org/cff"
	"context"
)

// Variadic is a flow that has a task that is variadic
func Variadic() {
	cff.Flow(context.Background(),
		cff.Task(func(s ...string) bool {
			return true
		}),
	)
}

// VariadicPredicate is a flow that has a task whose predicate is variadic
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
