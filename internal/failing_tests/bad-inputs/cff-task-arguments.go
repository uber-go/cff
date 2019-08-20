// +build cff

package badinputs

import (
	"context"

	"go.uber.org/cff"
)

// ExpectsFunctionCallExpression is a flow that doesn't provide a function to cff.Task.
func ExpectsFunctionCallExpression() {
	cff.Flow(context.Background(), cff.Task(
		true,
	))
}

// ExpectedFlowArgumentsSelectorExpression is a function that calls cff.Task with the wrong arguments but trickily
// passes the type checks.
// TODO: note doesn't trigger due to string being present in ExpectedFlowArgumentsNotCFF. Leaving for
// illustration purposes.
func ExpectedFlowArgumentsSelectorExpression() {
	cff.Flow(
		context.Background(),
		cff.Task(
			func() {},
			(cff.TaskOption)(nil),
		),
	)
}

// ExpectedFlowArgumentsCallExpressions is a function that calls cff.Task with the wrong arguments but trickily
// passes the type checks.
func ExpectedFlowArgumentsCallExpressions() {
	bad := cff.Instrument("")
	cff.Flow(
		context.Background(),
		cff.Task(
			func() {},
			bad,
		),
	)
}

// ExpectedFlowArgumentsNotCFF is a function that calls cff.Task with the wrong arguments but trickily
// passes the type checks.
func ExpectedFlowArgumentsNotCFF() {
	badFn := struct{ ProvideBad func() cff.TaskOption }{func() cff.TaskOption { return cff.Instrument("") }}
	cff.Flow(
		context.Background(),
		cff.Task(
			func() {},
			badFn.ProvideBad(),
		),
	)
}

// ExpectedTasksBad is a function that calls cff.Task with the wrong arguments.
func ExpectedTasksBad() {
	cff.Flow(
		context.Background(),
		cff.Task(
			nil,
		),
	)
}

// ExpectedTasksBadCallExpr is a function that calls cff.Task with the wrong arguments.
func ExpectedTasksBadCallExpr() {
	cff.Flow(
		context.Background(),
		cff.Task(
			cff.Params(),
		),
	)
}

// ExpectedTasksBadCallExprNotCFF is a function that calls cff.Task with the wrong arguments.
func ExpectedTasksBadCallExprNotCFF() {
	badFn := struct{ Task func() int }{func() int { return 0 }}
	cff.Flow(
		context.Background(),
		cff.Task(
			badFn.Task(),
		),
	)
}
