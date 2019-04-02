// +build cff

package badinputs

import (
	"go.uber.org/cff"
	"context"
)

// ExpectsFunctionCallExpression is a flow that doesn't provide a function to cff.Task.
func ExpectsFunctionCallExpression() {
	cff.Flow(context.Background(), cff.Task(
		true,
	))
}

// ExpectedFlowArgumentsSelectorExpression is a function that calls cff.Task with the wrong arguments but trickily
// passes the type checks.
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
		cff.Tasks(
			nil,
		),
	)
}

// ExpectedTasksBadCallExpr is a function that calls cff.Task with the wrong arguments.
func ExpectedTasksBadCallExpr() {
	cff.Flow(
		context.Background(),
		cff.Tasks(
			cff.Params(),
		),
	)
}

// ExpectedTasksBadCallExprNonCFF is a function that calls cff.Task with the wrong arguments.
func ExpectedTasksBadCallExprNonCFF() {
	badFn := struct{ Task func() bool }{func() bool { return true }}
	cff.Flow(
		context.Background(),
		cff.Tasks(
			badFn.Task(),
		),
	)
}
